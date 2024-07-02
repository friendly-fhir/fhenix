package cfg

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"text/template"

	"github.com/friendly-fhir/fhenix/internal/cfg"
	"github.com/friendly-fhir/fhenix/internal/templatefuncs"

	"gopkg.in/yaml.v3"
)

// Transform is a configuration for transforming input entities into templated
// output(s).
type Transform struct {
	// Include is a list of filters for conditions that an entity may satisfy to
	// be included in the transformation. At least one of these filters must be
	// satisfied for an entity to be included.
	Include []*TransformFilter `yaml:"include"`

	// Exclude is a list of filters for conditions that an entity may satisfy to
	// be excluded from the transformation. If any of these filters are satisfied,
	// the entity will be excluded, even if it satisfies an 'include' filter.
	Exclude []*TransformFilter `yaml:"exclude"`

	// OutputPath is a template for where the transformation will be output to.
	// The output path is always dynamic, and will be fed each individual entity
	// that has been matched by the filters to be transformed.
	OutputPath string `yaml:"output-path"`

	// Funcs is mapping of template function-names to files containing templates
	// that will perform textual transformations. This enables more complex logic
	// to be included in template pipelines.
	Funcs map[string]string `yaml:"funcs"`

	// Templates is a mapping of template names to files containing
	// templates that will be included in the main templates. This enables
	// re-use of common template snippets.
	//
	// Some special templates have implicit meaning in this configuration, but
	// this is otherwise freeform:
	//
	// - 'header': This template will be included at the top of the output file.
	// - 'footer': This template will be included at the bottom of the output file.
	// - 'type': This template will be called with each _individual entity_ 'type'
	//    that is matched by the filters.
	// - 'code-system': This template will be called with each _individual entity_
	//   'code-system' that is matched by the filters.
	// - 'main': This template will be called by the _list of all matched entities_.
	//   This template is provided by default by the implementation, which will call
	//   'header', followed by the appropriate intermediate template, followed by
	//   'footer'. Replacing this will replace all the above templates.
	Templates *TransformTemplates `yaml:"templates"`
}

func (t *Transform) UnmarshalYAML(node *yaml.Node) error {
	type transform Transform
	var out transform
	if err := node.Decode(&out); err != nil {
		return err
	}

	if err := verifyTemplate(out.OutputPath); err != nil {
		return &cfg.FieldError{Field: "transform.output-path", Err: err}
	}

	var errs []error
	for name := range out.Funcs {
		if !funcNameRegex.MatchString(name) {
			errs = append(errs, &cfg.FieldError{
				Field: "transform.funcs",
				Err:   fmt.Errorf("%w: '%v' is not a valid function name", cfg.ErrInvalidField, name),
			})
		}
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	*t = Transform(out)
	return nil
}

type TransformTemplates struct {
	// Header is a template that will be included at the top of the output file.
	Header string `yaml:"header"`

	// Footer is a template that will be included at the bottom of the output file.
	Footer string `yaml:"footer"`

	// Type is a template that will be called with each _individual entity_ 'type'
	// that is matched by the filters.
	Type string `yaml:"type"`

	// CodeSystem is a template that will be called with each _individual entity_
	// 'code-system' that is matched by the filters.
	CodeSystem string `yaml:"code-system"`

	// ValueSet is a template that will be called with each _individual entity_
	// 'value-set' that is matched by the filters.
	ValueSet string `yaml:"value-set"`

	// Main is a template that will be called by the _list of all matched entities_.
	// This template is provided by default by the implementation, which will call
	// 'header', followed by the appropriate intermediate template, followed by
	// 'footer'. Replacing this will replace all the above templates.
	Main string `yaml:"main"`

	// Partials is a custom set of additional templates to provide to the main
	// templates. This enables reuse of custom templates across transformations.
	Partials map[string]string `yaml:"-"`
}

func (tt *TransformTemplates) UnmarshalYAML(node *yaml.Node) error {
	out := map[string]string{}
	if err := node.Decode(&out); err != nil {
		return err
	}
	tt.Header = out["header"]
	tt.Footer = out["footer"]
	tt.Type = out["type"]
	tt.CodeSystem = out["code-system"]
	tt.ValueSet = out["value-set"]
	tt.Main = out["main"]

	delete(out, "header")
	delete(out, "footer")
	delete(out, "type")
	delete(out, "code-system")
	delete(out, "value-set")
	delete(out, "main")
	tt.Partials = out
	return nil
}

type TransformFilterType string

const (
	TransformFilterTypeStructureDefinition TransformFilterType = "StructureDefinition"
	TransformFilterTypeValueSet            TransformFilterType = "ValueSet"
	TransformFilterTypeCodeSystem          TransformFilterType = "CodeSystem"
)

type TransformFilter struct {
	// Name is a filter on the name of the input entity.
	// This may be a regular expression.
	Name string `yaml:"name"`

	// Type is the type of the input entity.
	Type string `yaml:"type"`

	// URL is an exact-match filter on the URL of the input entity.
	URL string `yaml:"url"`

	// Package is a filter on the package-name of the input entity.
	// This is only relevant if 'input.include-dependencies' is enabled.
	Package string `yaml:"package"`

	// Source is a filter on the source-file that the input entity is defined in.
	Source string `yaml:"source"`

	// Condition is a custom, template-filter that the input entity must satisfy.
	// This is meant as a back-door to allow for more complex conditions than
	// the other filters allow.
	Condition string `yaml:"condition"`
}

func verifyTemplate(v string) error {
	if strings.TrimSpace(v) == "" {
		return nil
	}
	_, err := template.New("").Funcs(templatefuncs.DefaultFuncs).Parse(v)
	return err
}

var (
	alnumRegex    = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	funcNameRegex = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)
)

func verifyRegex(v string) error {
	v = strings.TrimSpace(v)
	// Alphanumeric are exact matches, empty are ignored
	if v == "" || alnumRegex.MatchString(v) {
		return nil
	}

	_, err := regexp.Compile(v)
	return err
}

func hasSet(entries ...string) bool {
	for _, entry := range entries {
		if strings.TrimSpace(entry) != "" {
			return true
		}
	}
	return false
}

func (tf *TransformFilter) UnmarshalYAML(node *yaml.Node) error {
	type transformFilter TransformFilter
	var out transformFilter
	if err := node.Decode(&out); err != nil {
		return err
	}

	if !hasSet(out.Name, out.Type, out.URL, out.Package, out.Source, out.Condition) {
		return &cfg.FieldError{
			Field: "transform.filter",
			Err:   fmt.Errorf("%w: at least one filter option must be specified", cfg.ErrMissingField),
		}
	}

	if err := verifyRegex(out.Name); err != nil {
		return &cfg.FieldError{Field: "transform.filter.name", Err: err}
	}

	if err := verifyTemplate(out.Condition); err != nil {
		return &cfg.FieldError{Field: "transform.filter.condition", Err: err}
	}

	*tf = TransformFilter(out)

	return nil
}
