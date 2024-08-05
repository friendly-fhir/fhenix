package transformer

import (
	"os"

	"github.com/friendly-fhir/fhenix/internal/templatefuncs"
	"github.com/friendly-fhir/fhenix/pkg/transform/internal/template"
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

type config struct {
	funcs    map[string]any
	reporter templatefuncs.Reporter
}

type Option interface {
	apply(*config)
}

type option func(*config)

func (o option) apply(c *config) {
	o(c)
}

func WithFuncs(fns map[string]any) Option {
	return option(func(c *config) {
		c.funcs = fns
	})
}

func WithReporter(r templatefuncs.Reporter) Option {
	return option(func(c *config) {
		c.reporter = r
	})
}

// NewTemplate creates a new template using the underlying template engine.
func NewTemplate(engine template.Engine, templates map[string]string, opts ...Option) (template.Template, error) {
	var cfg config
	for _, opt := range opts {
		opt.apply(&cfg)
	}
	tmpl := engine.New("").Funcs(templatefuncs.NewFuncs(cfg.reporter)).Funcs(cfg.funcs)

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
