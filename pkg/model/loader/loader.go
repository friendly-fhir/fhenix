package loader

import (
	"context"
	"sync"

	"github.com/friendly-fhir/fhenix/internal/task"
	"github.com/friendly-fhir/fhenix/pkg/model/conformance"
	"github.com/friendly-fhir/fhenix/pkg/registry"
)

type Listener interface {
	BeforeLoadPackage(ref registry.PackageRef)
	AfterLoadPackage(ref registry.PackageRef, err error)
	listener()
}

type BaseListener struct{}

func (BaseListener) BeforeLoadPackage(ref registry.PackageRef) {}

func (BaseListener) AfterLoadPackage(ref registry.PackageRef, err error) {}

func (BaseListener) listener() {}

type Option interface {
	set(*Loader)
}

type option func(*Loader)

func (o option) set(l *Loader) {
	o(l)
}

// WithModule returns an [Option] for the [Loader] that will set the conformance
// module to use for loading.
func WithModule(module *conformance.Module) Option {
	return option(func(l *Loader) {
		l.module = module
	})
}

// WithWorkers returns an [Option] for the [Loader] that will set the number of
// parallel worker threads.
func WithWorkers(parallel int) Option {
	return option(func(l *Loader) {
		l.parallel = parallel
	})
}

// WithListeners returns an [Option] for the [Loader] that will set the listeners
// to notify when a package is loaded.
func WithListeners(listeners ...Listener) Option {
	return option(func(l *Loader) {
		l.listeners = append(l.listeners, listeners...)
	})
}

// Loader is a FHIR definition loader that loads definitions from a registry.
type Loader struct {
	m         sync.Mutex
	module    *conformance.Module
	cache     *registry.Cache
	parallel  int
	loaded    *sync.Map
	listeners []Listener
}

// New constructs a new loader with the given cache and options.
func New(cache *registry.Cache, opts ...Option) *Loader {
	result := &Loader{
		cache:    cache,
		module:   conformance.DefaultModule(),
		parallel: 0,
		loaded:   &sync.Map{},
	}
	for _, opt := range opts {
		opt.set(result)
	}
	return result
}

// Load loads the FHIR definitions for the given package reference.
func (l *Loader) Load(refs ...registry.PackageRef) error {
	runner := task.NewRunner(l.parallel)
	for _, ref := range refs {
		runner.Add(l.load(ref))
	}
	_, err := runner.Run(context.Background())
	return err
}

func (l *Loader) load(ref registry.PackageRef) task.Task {
	return task.Func(func(ctx context.Context) error {
		if _, ok := l.loaded.LoadOrStore(ref, struct{}{}); ok {
			return nil
		}

		pkg, err := l.cache.Get(ref.Parts())
		if err != nil {
			return err
		}
		for _, listener := range l.listeners {
			listener.BeforeLoadPackage(ref)
		}
		l.m.Lock()
		err = l.module.FromPackage(pkg)
		l.m.Unlock()
		for _, listener := range l.listeners {
			listener.AfterLoadPackage(ref, err)
		}
		if err != nil {
			return err
		}

		runner := task.RunnerFromContext(ctx)
		for name, version := range pkg.Dependencies() {
			runner.Add(l.load(registry.NewPackageRef(ref.Registry(), name, version)))
		}
		return nil
	})
}

// Module returns the conformance module used by the loader.
func (l *Loader) Module() *conformance.Module {
	return l.module
}
