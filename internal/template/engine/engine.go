package engine

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/friendly-fhir/fhenix/internal/config"
	"github.com/friendly-fhir/fhenix/internal/model"
	"github.com/friendly-fhir/fhenix/internal/template"
)

// Engine is a template engine that processes a model and applies transformations.
type Engine struct {
	cfg    *config.Config
	output string
}

type transform[T any] struct {
	Template *template.Template
	Inputs   map[string][]*T
}

func parse(tmpl **template.Template, name string, filename string) error {
	if filename == "" {
		return nil
	}
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	content, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	out, err := template.Parse(name, string(content))
	*tmpl = out
	return err
}

func (e *Engine) buildTypeTransforms(m *model.Model, transformation *config.Transformation) (*transform[model.Type], error) {
	result := &transform[model.Type]{
		Inputs: map[string][]*model.Type{},
	}
	if err := parse(&result.Template, "content", transformation.Template); err != nil {
		return nil, err
	}

	for _, t := range m.Types().All() {
		can := transformation.Input.If.Evaluate(t)
		if !can {
			continue
		}
		path, err := transformation.Output.Evaluate(t)
		if err != nil {
			return nil, err
		}
		result.Inputs[path] = append(result.Inputs[path], t)
	}
	return result, nil
}

func (e *Engine) buildCodeSystemTransforms(m *model.Model, transformation *config.Transformation) (*transform[model.CodeSystem], error) {
	result := &transform[model.CodeSystem]{
		Inputs: map[string][]*model.CodeSystem{},
	}
	if err := parse(&result.Template, "content", transformation.Template); err != nil {
		return nil, err
	}

	for _, cs := range m.CodeSystems() {
		can := transformation.Input.If.Evaluate(cs)
		if !can {
			continue
		}
		path, err := transformation.Output.Evaluate(cs)
		if err != nil {
			return nil, err
		}
		result.Inputs[path] = append(result.Inputs[path], cs)
	}
	return result, nil
}

// Run executes the template engine.
func (e *Engine) Run(m *model.Model) error {
	var errs []error
	for _, transform := range e.cfg.Transformations {
		if transform.Input.Type == "StructureDefinition" {
			transforms, err := e.buildTypeTransforms(m, &transform)
			if err != nil {
				errs = append(errs, err)
				continue
			}
			if err := run(m, e.output, transforms); err != nil {
				errs = append(errs, err)
				continue
			}
		}
		if transform.Input.Type == "CodeSystem" {
			transforms, err := e.buildCodeSystemTransforms(m, &transform)
			if err != nil {
				errs = append(errs, err)
				continue
			}
			if err := run(m, e.output, transforms); err != nil {
				errs = append(errs, err)
				continue
			}
		}
	}
	return errors.Join(errs...)
}

func run[T any](m *model.Model, output string, transforms *transform[T]) error {
	var errs []error
	for path, types := range transforms.Inputs {
		if len(types) != 1 {
			errs = append(errs, fmt.Errorf("expected 1 type for output %v, got %d", path, len(types)))
			continue
		}
		outpath := filepath.FromSlash(path)
		if output != "" {
			outpath = filepath.Join(output, outpath)
		}
		fmt.Println("output: ", outpath)

		if err := os.MkdirAll(filepath.Dir(outpath), 0755); err != nil {
			errs = append(errs, err)
			continue
		}
		file, err := os.Create(outpath)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		defer file.Close()
		if err := transforms.Template.Execute(file, types[0]); err != nil {
			errs = append(errs, err)
			continue
		}
	}
	return errors.Join(errs...)
}

// Option is a functional option for Engine.
type Option interface {
	set(*Engine)
}

type option func(*Engine)

func (o option) set(e *Engine) {
	o(e)
}

var _ Option = (*option)(nil)

// Output returns an option that sets the output directory on Engine creation.
func Output(output string) Option {
	return option(func(e *Engine) {
		e.output = output
	})
}

// New creates a new template Engine.
func New(cfg *config.Config, opts ...Option) *Engine {
	engine := &Engine{
		cfg: cfg,
	}
	for _, opt := range opts {
		opt.set(engine)
	}
	return engine
}
