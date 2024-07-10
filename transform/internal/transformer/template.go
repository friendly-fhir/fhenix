package transformer

import (
	"os"

	"github.com/friendly-fhir/fhenix/internal/templatefuncs"
	"github.com/friendly-fhir/fhenix/transform/internal/template"
)

const (
	// DefaultMainTemplate is the default template for the main template execution.
	DefaultMainTemplate string = `
{{- range .StructureDefinitions }}{{ template "structure-definition" . }}{{ end -}}
{{- range .ValueSets }}{{ template "value-set" . }}{{ end -}}
{{- range .CodeSystems }}{{ template "code-system" . }}{{ end -}}
`

	// DefaultEntryTemplate is the default template used by the template engine.
	DefaultEntryTemplate string = `
{{- template "header" . }}{{ template "main" . }}{{ template "footer" . -}}
`
)

// NewTemplate creates a new template using the underlying template engine.
func NewTemplate(engine template.Engine, templates map[string]string, fns map[string]any) (template.Template, error) {
	tmpl := engine.New("").Funcs(templatefuncs.DefaultFuncs).Funcs(fns)

	defaults := map[string]string{
		"main":                 DefaultMainTemplate,
		"header":               "",
		"footer":               "",
		"structure-definition": "",
		"code-system":          "",
		"value-set":            "",
	}

	for name, path := range templates {
		if err := parse(tmpl, name, path); err != nil {
			return nil, err
		}
	}

	for name, template := range defaults {
		if _, ok := templates[name]; ok {
			continue
		}

		_, _ = tmpl.New(name).Parse(template)
	}

	_, _ = tmpl.Parse(DefaultEntryTemplate)
	return tmpl, nil
}

func parse(tmpl template.Template, name string, path string) error {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	_, err = tmpl.New(name).Parse(string(bytes))
	return err
}
