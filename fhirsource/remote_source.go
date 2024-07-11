package fhirsource

import (
	"context"

	"github.com/friendly-fhir/fhenix/internal/fhirig"
)

// PackageCache is a mechanism for caching FHIR packages.
type PackageCache interface {
	// Dependencies returns the dependencies of the specified package.
	Dependencies(pkg *Package) ([]*Package, error)

	// FetchAndGet fetches the specified package and returns the files.
	FetchAndGet(ctx context.Context, pkg *Package) ([]string, error)
}

// NewRemoteSource constructs a new remote Source from the specified package.
func NewRemoteSource(pkg *Package, cache PackageCache) Source {
	return &remoteSource{
		pkg:   pkg,
		cache: cache,
	}
}

type Package = fhirig.Package

type remoteSource struct {
	pkg   *fhirig.Package
	cache PackageCache
}

func (rs *remoteSource) Bundles(ctx context.Context) ([]*Bundle, error) {
	dependencies, err := rs.cache.Dependencies(rs.pkg)
	if err != nil {
		return nil, err
	}
	packages := []*Package{rs.pkg}
	packages = append(packages, dependencies...)

	var result []*Bundle
	for _, pkg := range packages {
		files, err := rs.cache.FetchAndGet(ctx, pkg)
		if err != nil {
			return nil, err
		}
		result = append(result, &Bundle{
			Package: pkg,
			Files:   files,
		})
	}
	return result, nil
}

var _ Source = (*remoteSource)(nil)
