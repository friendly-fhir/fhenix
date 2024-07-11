/*
Package driver is the primary driver behind the application.
It is responsible for the main loop of the application which does the actual
code generation.
*/
package driver

import (
	"context"
	"fmt"

	"github.com/friendly-fhir/fhenix/config"
	"github.com/friendly-fhir/fhenix/driver/job"
	"github.com/friendly-fhir/fhenix/fhirsource"
	"github.com/friendly-fhir/fhenix/model"
	"github.com/friendly-fhir/fhenix/model/conformance"
	"github.com/friendly-fhir/fhenix/transform"
)

type Listener interface {
	fhirsource.Listener
}

type BaseListener struct {
	fhirsource.BaseListener
}

type Driver struct {
	Source     fhirsource.Source
	model      *model.Model
	transforms []*transform.Transform
	listener   Listener
	outputPath string

	includeDependencies bool
}

func commonBase(t1, t2 *model.Type) *model.Type {
	if t1 == t2 {
		return t1
	}
	if t1.Base == nil && t2.Base == nil {
		return nil
	}
	if t1.Base != nil {
		t1 = t1.Base
	}
	if t2.Base != nil {
		t2 = t2.Base
	}
	return commonBase(t1, t2)
}

func CommonBase(types []*model.Type) *model.Type {
	if len(types) == 0 {
		return nil
	}
	base := types[0]
	for _, t := range types[1:] {
		base = commonBase(base, t)
	}
	return base
}

func New(config *config.Config, listener Listener) (*Driver, error) {
	source, err := fhirsource.New(config.Input, listener)
	if err != nil {
		return nil, err
	}

	transforms := make([]*transform.Transform, 0, len(config.Transforms))
	for _, t := range config.Transforms {
		fmt.Println("Creating transform " + t.OutputPath)
		t, err := transform.New(config.Mode, t, transform.Funcs{
			"commonbase": CommonBase,
		})
		if err != nil {
			return nil, err
		}
		transforms = append(transforms, t)
	}

	driver := &Driver{
		Source:     source,
		transforms: transforms,
		listener:   listener,
		outputPath: config.OutputDir,

		includeDependencies: config.Input.IncludeDependencies,
	}
	return driver, nil
}

func (d *Driver) Run(ctx context.Context) error {
	module := conformance.DefaultModule()

	bundles, err := d.Source.Bundles(ctx)
	if err != nil {
		return err
	}

	for _, bundle := range bundles {
		for _, file := range bundle.Files {
			if err := module.ParseFile(file, bundle.Package); err != nil {
				fmt.Println("Error parsing file '"+file+"':", err)
			}
		}
	}
	d.model = model.NewModel(module)
	if err := d.model.DefineAllTypes(); err != nil {
		return err
	}

	var jobs []*job.Job
	for _, t := range d.transforms {
		j, err := job.New(d.model, d.outputPath, t)
		if err != nil {
			return err
		}
		jobs = append(jobs, j...)
	}

	fmt.Println("Outputs:")
	for _, j := range jobs {
		fmt.Println("\t" + j.OutputPath())
		if err := j.Execute(ctx); err != nil {
			return err
		}
	}
	return nil
}
