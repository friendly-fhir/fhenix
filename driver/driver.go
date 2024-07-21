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

	"github.com/friendly-fhir/fhenix/config"
	"github.com/friendly-fhir/fhenix/driver/job"
	"github.com/friendly-fhir/fhenix/internal/task"
	"github.com/friendly-fhir/fhenix/internal/templatefuncs"
	"github.com/friendly-fhir/fhenix/model"
	"github.com/friendly-fhir/fhenix/model/conformance"
	"github.com/friendly-fhir/fhenix/model/loader"
	"github.com/friendly-fhir/fhenix/registry"
	"github.com/friendly-fhir/fhenix/transform"
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

type Listener interface {
	BeforeDownload()
	AfterDownload(err error)

	BeforeLoadTransform()
	AfterLoadTransform(err error)

	BeforeLoadConformance()
	AfterLoadConformance(err error)

	BeforeLoadModel()
	AfterLoadModel(err error)

	BeforeTransform()
	OnTransform(output string)
	AfterTransform(jobs int, err error)

	loader.Listener
	registry.CacheListener
}

type BaseListener struct {
	loader.BaseListener
	registry.BaseCacheListener
}

func (BaseListener) BeforeDownload()     {}
func (BaseListener) AfterDownload(error) {}

func (BaseListener) BeforeLoadTransform()     {}
func (BaseListener) AfterLoadTransform(error) {}

func (BaseListener) BeforeLoadConformance()     {}
func (BaseListener) AfterLoadConformance(error) {}

func (BaseListener) BeforeLoadModel()     {}
func (BaseListener) AfterLoadModel(error) {}

func (BaseListener) BeforeTransform()          {}
func (BaseListener) OnTransform(output string) {}
func (BaseListener) AfterTransform(int, error) {}

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
		explicitPackages: []registry.PackageRef{
			registry.NewPackageRef("default", config.Input.Name, config.Input.Version),
		},
	}
	for _, opt := range opts {
		opt.set(driver)
	}

	return driver, nil
}

func (d *Driver) DownloadPackages(ctx context.Context) error {
	for _, listener := range d.listeners {
		listener.BeforeDownload()
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
		listener.AfterDownload(err)
	}
	return err
}

func (d *Driver) LoadTransforms() ([]*transform.Transform, error) {
	for _, listener := range d.listeners {
		listener.BeforeLoadTransform()
	}

	var errs []error
	transforms := make([]*transform.Transform, 0, len(d.transformConfigs))
	for _, t := range d.transformConfigs {
		t, err := transform.New(d.mode, t, transform.WithFuncs(Funcs), transform.WithReporter(d.reporter))
		if err != nil {
			errs = append(errs, err)
			continue
		}
		transforms = append(transforms, t)
	}

	err := errors.Join(errs...)
	for _, listener := range d.listeners {
		listener.AfterLoadTransform(err)
	}
	return transforms, err
}

func (d *Driver) LoadConformanceModule() error {
	for _, listener := range d.listeners {
		listener.BeforeLoadConformance()
	}
	loader := loader.New(d.cache,
		loader.WithModule(d.module),
		loader.WithWorkers(d.parallel),
		loader.WithListeners(toListeners[loader.Listener](d.listeners)...),
	)
	err := loader.Load(d.explicitPackages...)
	for _, listener := range d.listeners {
		listener.AfterLoadConformance(err)
	}
	return err
}

func (d *Driver) LoadModel() (*model.Model, error) {
	for _, listener := range d.listeners {
		listener.BeforeLoadModel()
	}
	model := model.NewModel(d.module)
	err := model.DefineAllTypes()
	for _, listener := range d.listeners {
		listener.AfterLoadModel(err)
	}
	return model, err
}

func (d *Driver) Transform(ctx context.Context, model *model.Model, transforms []*transform.Transform) error {
	for _, listener := range d.listeners {
		listener.BeforeTransform()
	}
	runner := task.NewRunner(d.parallel)

	for _, t := range transforms {
		jobs, err := job.New(model, d.outputPath, t)
		if err != nil {
			return err
		}
		for _, job := range jobs {
			runner.Add(task.Func(func(ctx context.Context) error {
				for _, listener := range d.listeners {
					listener.OnTransform(job.OutputPath())
				}

				err := job.Execute(ctx)
				return err
			}))
		}
	}

	jobs, err := runner.Run(ctx)
	for _, listener := range d.listeners {
		listener.AfterTransform(jobs, err)
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
