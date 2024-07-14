/*
Package fhirsource provides a wrapper for retrieving FHIR source definitions
from the fhenix configuration format.

DEPRECATED: This package is deprecated and will be removed in a future release.
*/
package fhirsource

import (
	"context"

	"github.com/friendly-fhir/fhenix/config"
	"github.com/friendly-fhir/fhenix/internal/fhirig"
)

// Listener is a mechanism for listening to fetch events from remote package
// sources.
type Listener interface {
	// OnFetchStart is called when a package fetch is started.
	OnFetchStart(pkg *Package)

	// OnFetchEnd is called when a package fetch is completed.
	OnFetchEnd(pkg *Package, err error)

	// OnCacheHit is called when a package fetch is completed and the cache was
	// hit.
	OnCacheHit(pkg *Package)
}

type BaseListener struct{}

func (l *BaseListener) OnFetchStart(*Package) {}

func (l *BaseListener) OnFetchEnd(*Package, error) {}

func (l *BaseListener) OnCacheHit(*Package) {}

// New constructs a new [Source] from the specified configuration.
func New(cfg *config.Package, listener Listener) (Source, error) {
	pkg := fhirig.NewPackage(cfg.Name, cfg.Version)
	if cfg.Path != "" {
		return NewLocalSource(pkg, cfg.Path), nil
	}
	if listener == nil {
		listener = &fhirig.BaseListener{}
	}
	cache := fhirig.NewSystemCache()
	cache.Listener = listener
	return NewRemoteSource(pkg, cache), nil
}

// Source is a mechanism for retrieving JSON FHIR definitions.
//
// This abstracts local from remote sources, so that a path may exist on disk
// or by a simplify.net registry package.
type Source interface {
	// Bundles returns a list of the defined package bundles in the fhir source.
	Bundles(ctx context.Context) ([]*Bundle, error)
}

// Bundle is a collection of files that are part of a FHIR package.
type Bundle struct {
	Package *fhirig.Package
	Files   []string
}
