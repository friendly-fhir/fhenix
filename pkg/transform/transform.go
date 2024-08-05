package transform

import (
	"io"
	"path/filepath"
	"strings"

	"github.com/friendly-fhir/fhenix/internal/templatefuncs"
	"github.com/friendly-fhir/fhenix/pkg/config"
	"github.com/friendly-fhir/fhenix/pkg/filter"
	"github.com/friendly-fhir/fhenix/pkg/transform/internal/template"
	"github.com/friendly-fhir/fhenix/pkg/transform/internal/transformer"
)

// Transform represents a transformation to be applied to the input definitions.
type Transform struct {
	include  filter.Filters
	exclude  filter.Filters
	output   template.Template
	template template.Template
}

type transformConfig struct {
	funcs    Funcs
	reporter templatefuncs.Reporter
}

type Option interface {
	set(*transformConfig)
}

type option func(*transformConfig)

func (o option) set(c *transformConfig) {
	o(c)
}

func WithFuncs(fns Funcs) Option {
	return option(func(c *transformConfig) {
		c.funcs = fns
	})
}

func WithReporter(r templatefuncs.Reporter) Option {
	return option(func(c *transformConfig) {
		c.reporter = r
	})
}

type Funcs map[string]any

func New(mode config.Mode, transform *config.Transform, opts ...Option) (*Transform, error) {
	if transform == nil {
		panic("transform: New called with nil transform")
	}
	var cfg transformConfig
	for _, opt := range opts {
		opt.set(&cfg)
	}
	engine, err := template.FromString(string(mode))
	if err != nil {
		return nil, err
	}

	funcs, err := transformer.FuncsFromConfig(transform.Funcs, cfg.reporter)
	if err != nil {
		return nil, err
	}
	for name, fn := range cfg.funcs {
		funcs[name] = fn
	}

	tmpl, err := transformer.NewTemplate(engine, transform.Templates,
		transformer.WithFuncs(funcs),
		transformer.WithReporter(cfg.reporter),
	)
	if err != nil {
		return nil, err
	}

	output, err := engine.New("").Funcs(templatefuncs.NewFuncs(cfg.reporter)).Parse(transform.OutputPath)
	if err != nil {
		return nil, err
	}

	result := &Transform{
		template: tmpl,
		include:  filter.New(transform.Include...),
		exclude:  filter.New(transform.Exclude...),
		output:   output,
	}

	return result, nil
}

// CanTransform returns true if the given value should be included in the
// output transformation.
func (t *Transform) CanTransform(v any) bool {
	if t == nil {
		return false
	}

	if len(t.include) > 0 && !t.include.Matches(v) {
		return false
	}

	return !t.exclude.Matches(v)
}

// OutputPath returns the output path for the given value.
// The output path is always specified as an absolute path.
func (t *Transform) OutputPath(v any) (string, error) {
	if t == nil {
		return "", nil
	}

	var sb strings.Builder
	if err := t.output.Execute(&sb, v); err != nil {
		return "", err
	}
	return filepath.FromSlash(strings.TrimSpace(sb.String())), nil
}

// Execute the transformation with the given data.
func (t *Transform) Execute(w io.Writer, data any) error {
	if t == nil || t.template == nil {
		return nil
	}
	return t.template.Execute(w, data)
}
