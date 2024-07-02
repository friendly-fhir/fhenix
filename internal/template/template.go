/*
Package template provides the templated solutions to working with FHENIX
configurations.

This package wraps the [text/template] package to provide a more user-friendly
interface, and provides sensible defaults for the substitution.

[text/template]: https://golang.org/pkg/text/template/
*/
package template

import (
	"encoding/json"
	"strconv"
	"strings"
	"text/template"

	"github.com/friendly-fhir/fhenix/internal/templatefuncs"
	"gopkg.in/yaml.v3"
)

type FuncMap = template.FuncMap

type Template struct {
	*template.Template
}

// New allocates a new, undefined template with the given name.
func New(name string) *Template {
	return &Template{template.New(name).Funcs(templatefuncs.DefaultFuncs)}
}

// Parse parses a string into a new template with the given name.
func Parse(name, text string) (*Template, error) {
	return New(name).Parse(text)
}

// MustParse parses a string into a new template with the given name.
// It panics if the template cannot be parsed.
func MustParse(name, text string) *Template {
	return Must(Parse(name, text))
}

func (t *Template) UnmarshalYAML(node *yaml.Node) error {
	if t.Template == nil {
		t.Template = template.New("").Funcs(templatefuncs.DefaultFuncs)
	}

	var text string
	if err := node.Decode(&text); err != nil {
		return err
	}

	_, err := t.Parse(text)
	return err
}

// Templates returns a slice of defined templates associated with t.
func (t *Template) Templates() []*Template {
	templates := t.Template.Templates()
	result := make([]*Template, 0, len(templates))
	for _, tmpl := range templates {
		result = append(result, &Template{tmpl})
	}
	return result
}

// Funcs adds the elements of the argument map to the template's function map.
// It must be called before the template is parsed.
// It panics if a value in the map is not a function with appropriate return
// type or if the name cannot be used syntactically as a function in a template.
// It is legal to overwrite elements of the map. The return value is the template,
// so calls can be chained.
func (t *Template) Funcs(funcs FuncMap) *Template {
	t.Template.Funcs(funcs)
	return t
}

// Parse parses text as a template body for t.
// Named template definitions ({{define ...}} or {{block ...}} statements) in text
// define additional templates associated with t and are removed from the
// definition of t itself.
//
// Templates can be redefined in successive calls to Parse.
// A template definition with a body containing only white space and comments
// is considered empty and will not replace an existing template's body.
// This allows using Parse to add new named template definitions without
// overwriting the main template body.
func (t *Template) Parse(text string) (*Template, error) {
	tmpl, err := t.Template.Parse(text)
	if err != nil {
		return nil, err
	}
	t.Template = tmpl
	return t, nil
}

// Lookup returns the template with the given name that is associated with t.
// It returns nil if there is no such template or the template has no definition.
func (t *Template) Lookup(name string) *Template {
	return &Template{t.Template.Lookup(name)}
}

// ExecuteBool applies a parsed template to the specified data object,
// and returns the resulting value as a boolean.
//
// This returns true if under the following conditions
//
//   - if the evaluated result is numeric with a value greater than 0
//   - if the evaluated result is boolean with a value of true
//   - if the evaluated result is a JSON sequence or mapping with at least 1 entry
//   - if the evaluated result is a JSON string with a non-zero length
//   - if the result is not valid JSON and is a boolean-parseable value of true
//   - if the result is not valid JSON and the output is a non-zero-length
//     and non-empty string.
//
// This returns false if under the following conditions
//
//   - if the evaluated result is numeric with a value of 0
//   - if the evaluated result is boolean with a value of false
//   - if the evaluated result is a JSON sequence or mapping with no entries
//   - if the evaluated result is a JSON string with a zero length
//   - if the result is not valid JSON and is a boolean-parseable value of false
//   - if the result is not valid JSON and the output is a zero-length.
//
// In all other cases, this returns false.
func (t *Template) ExecuteBool(data any) (bool, error) {
	var sb strings.Builder
	err := t.Execute(&sb, data)
	if err != nil {
		return false, err
	}

	str := strings.TrimSpace(sb.String())
	var j any
	if err := json.Unmarshal([]byte(str), &j); err == nil {
		switch j := j.(type) {
		case []any:
			return len(j) > 0, nil
		case map[string]any:
			return len(j) > 0, nil
		case nil:
			return false, nil
		case string:
			str := strings.TrimSpace(j)
			return len(str) > 0, nil
		case bool:
			return j, nil
		case float64:
			return j > 0, nil
		}
	}

	// This enables TRUE to be used as input
	if truth, err := strconv.ParseBool(str); err == nil {
		return truth, nil
	}
	return len(str) > 0, nil
}

// ExecuteString applies a parsed template to the specified data object,
// and returns the resulting value as a string.
func (t *Template) ExecuteString(data any) (string, error) {
	var sb strings.Builder
	err := t.Execute(&sb, data)
	return strings.TrimSpace(sb.String()), err
}

// Must is a helper that wraps a call to a function returning (*Template, error)
// and panics if the error is non-nil. It is intended for use in variable
// initializations such as
//
//	var t = template.Must(template.New("name").Parse("text"))
func Must(t *Template, err error) *Template {
	return &Template{template.Must(t.Template, err)}
}
