/*
Package driver is the primary driver behind the application.
It is responsible for the main loop of the application which does the actual
code generation.
*/
package driver

import (
	"context"
	"errors"
	"runtime"

	"github.com/friendly-fhir/fhenix/internal/task"
	"github.com/friendly-fhir/fhenix/internal/templatefuncs"
	"github.com/friendly-fhir/fhenix/pkg/config"
	"github.com/friendly-fhir/fhenix/pkg/driver/job"
	"github.com/friendly-fhir/fhenix/pkg/model"
	"github.com/friendly-fhir/fhenix/pkg/model/conformance"
	"github.com/friendly-fhir/fhenix/pkg/model/loader"
	"github.com/friendly-fhir/fhenix/pkg/registry"
	"github.com/friendly-fhir/fhenix/pkg/transform"
)

type Option interface {
	set(*Driver)
}

type option func(*Driver)

func (o option) set(d *Driver) {
	o(d)
}

var _ Option = (*option)(nil)

type Reporter = templatefuncs.Reporter

type Stage int

const (
	StageDownload Stage = iota
	StageLoadTransform
	StageLoadConformance
	StageLoadModel
	StageTransform
)

type Listener interface {
	BeforeStage(stage Stage)
	AfterStage(stage Stage, err error)

	BeforeLoadTransform(n int)
	AfterLoadTransform(n int, err error)

	BeforeTransform(n int, jobs int)
	OnTransformOutput(n int, output string)
	AfterTransformOutput(n int, output string, err error)

	loader.Listener
	registry.CacheListener
}

type BaseListener struct {
	loader.BaseListener
	registry.BaseCacheListener
}

func (BaseListener) BeforeStage(Stage)       {}
func (BaseListener) AfterStage(Stage, error) {}

func (BaseListener) BeforeLoadTransform(n int)           {}
func (BaseListener) AfterLoadTransform(n int, err error) {}

func (BaseListener) BeforeTransform(i, jobs int)                          {}
func (BaseListener) OnTransformOutput(i int, output string)               {}
func (BaseListener) AfterTransformOutput(i int, output string, err error) {}

var _ Listener = (*BaseListener)(nil)

type Driver struct {
	outputPath string

	module *conformance.Module
	cache  *registry.Cache

	mode config.Mode

	transformConfigs []*config.Transform

	forceDownload    bool
	parallel         int
	explicitPackages []registry.PackageRef

	listeners []Listener
	reporter  templatefuncs.Reporter
}

// Cache returns an [Option] for the [Driver] that will set the cache to use
// for downloading packages.
func Cache(cache *registry.Cache) Option {
	return option(func(d *Driver) {
		d.cache = cache
	})
}

// Parallel returns an [Option] for the [Driver] that will set the number of
// parallel worker threads.
func Parallel(parallel int) Option {
	return option(func(d *Driver) {
		d.parallel = parallel
	})
}

// ConformanceModule returns an [Option] for the [Driver] that will set the
// conformance module to use for loading.
func ConformanceModule(module *conformance.Module) Option {
	return option(func(d *Driver) {
		d.module = module
	})
}

// ExplicitPackages returns an [Option] for the [Driver] that will set the
// explicit packages to download.
func ExplicitPackages(packages ...registry.PackageRef) Option {
	return option(func(d *Driver) {
		d.explicitPackages = append(d.explicitPackages, packages...)
	})
}

// ForceDownload returns an [Option] for the [Driver] that will set whether
// to force download of packages.
func ForceDownload(force bool) Option {
	return option(func(d *Driver) {
		d.forceDownload = force
	})
}

// Listeners returns an [Option] for the [Driver] that will set the
// listeners to notify when a package is downloaded or loaded.
func Listeners(listeners ...Listener) Option {
	return option(func(d *Driver) {
		d.listeners = append(d.listeners, listeners...)
	})
}

// TemplateReporter returns an [Option] for the [Driver] that will set the
// reporter to use for reporting template errors.
func TemplateReporter(reporter Reporter) Option {
	return option(func(d *Driver) {
		d.reporter = reporter
	})
}

// TemplateReportFunc returns an [Option] for the [Driver] that will set the
// reporter to use for reporting template errors.
func TemplateReportFunc(report func(error)) Option {
	return option(func(d *Driver) {
		d.reporter = templatefuncs.ReporterFunc(report)
	})
}

func New(config *config.Config, opts ...Option) (*Driver, error) {
	driver := &Driver{
		outputPath: config.OutputDir,
		parallel:   runtime.NumCPU(),

		mode: config.Mode,

		module: conformance.DefaultModule(),
		cache:  registry.DefaultCache(),

		forceDownload:    false,
		transformConfigs: config.Transforms,
	}
	for _, pkg := range config.Input {
		driver.explicitPackages = append(driver.explicitPackages, registry.NewPackageRef("default", pkg.Name, pkg.Version))
	}
	for _, opt := range opts {
		opt.set(driver)
	}

	return driver, nil
}

func (d *Driver) DownloadPackages(ctx context.Context) error {
	for _, listener := range d.listeners {
		listener.BeforeStage(StageDownload)
	}

	downloader := registry.NewDownloader(d.cache).Force(d.forceDownload).Workers(d.parallel)
	for _, listener := range toListeners[registry.CacheListener](d.listeners) {
		d.cache.AddListener(listener)
	}

	for _, ref := range d.explicitPackages {
		registry, name, version := ref.Parts()
		downloader.Add(registry, name, version, true)
	}

	err := downloader.Start(ctx)
	for _, listener := range d.listeners {
		listener.AfterStage(StageDownload, err)
	}
	return err
}

func (d *Driver) LoadTransforms() ([]*transform.Transform, error) {
	for _, listener := range d.listeners {
		listener.BeforeStage(StageLoadTransform)
	}

	var errs []error
	transforms := make([]*transform.Transform, 0, len(d.transformConfigs))
	for i, t := range d.transformConfigs {
		for _, listener := range d.listeners {
			listener.BeforeLoadTransform(i)
		}
		t, err := transform.New(d.mode, t, transform.WithFuncs(Funcs), transform.WithReporter(d.reporter))
		for _, listener := range d.listeners {
			listener.AfterLoadTransform(i, err)
		}
		if err != nil {
			errs = append(errs, err)
			continue
		}
		transforms = append(transforms, t)
	}

	err := errors.Join(errs...)
	for _, listener := range d.listeners {
		listener.AfterStage(StageLoadTransform, err)
	}
	return transforms, err
}

func (d *Driver) LoadConformanceModule() error {
	for _, listener := range d.listeners {
		listener.BeforeStage(StageLoadConformance)
	}
	loader := loader.New(d.cache,
		loader.WithModule(d.module),
		loader.WithWorkers(d.parallel),
		loader.WithListeners(toListeners[loader.Listener](d.listeners)...),
	)
	err := loader.Load(d.explicitPackages...)
	for _, listener := range d.listeners {
		listener.AfterStage(StageLoadConformance, err)
	}
	return err
}

func (d *Driver) LoadModel() (*model.Model, error) {
	for _, listener := range d.listeners {
		listener.BeforeStage(StageLoadModel)
	}
	model := model.NewModel(d.module)
	err := model.DefineAllTypes()
	for _, listener := range d.listeners {
		listener.AfterStage(StageLoadModel, err)
	}
	return model, err
}

func (d *Driver) Transform(ctx context.Context, model *model.Model, transforms []*transform.Transform) error {
	for _, listener := range d.listeners {
		listener.BeforeStage(StageTransform)
	}
	runner := task.NewRunner(d.parallel)

	for i, t := range transforms {
		jobs, err := job.New(model, d.outputPath, t)
		for _, listener := range d.listeners {
			listener.BeforeTransform(i, len(jobs))
		}
		if err != nil {
			return err
		}
		for _, job := range jobs {
			runner.Add(task.Func(func(ctx context.Context) error {
				for _, listener := range d.listeners {
					listener.OnTransformOutput(i, job.OutputPath())
				}

				err := job.Execute(ctx)

				for _, listener := range d.listeners {
					listener.AfterTransformOutput(i, job.OutputPath(), err)
				}
				return err
			}))
		}
	}

	_, err := runner.Run(ctx)
	for _, listener := range d.listeners {
		listener.AfterStage(StageTransform, err)
	}
	return err
}

func (d *Driver) Run(ctx context.Context) error {
	if err := d.DownloadPackages(ctx); err != nil {
		return err
	}
	if err := d.LoadConformanceModule(); err != nil {
		return err
	}
	model, err := d.LoadModel()
	if err != nil {
		return err
	}
	transforms, err := d.LoadTransforms()
	if err != nil {
		return err
	}
	return d.Transform(ctx, model, transforms)
}

func toListeners[T any](listeners []Listener) []T {
	result := make([]T, len(listeners))
	for i, l := range listeners {
		result[i] = l.(T)
	}
	return result
}
