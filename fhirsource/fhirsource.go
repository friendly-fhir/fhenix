/*
Package fhirsource provides a wrapper for retrieving FHIR source definitions
from the fhenix configuration format.
*/
package fhirsource

import (
	"context"

	"github.com/friendly-fhir/fhenix/config"
	"github.com/friendly-fhir/fhenix/internal/fhirig"
)

// Source is a mechanism for retrieving JSON FHIR definitions.
//
// This abstracts local from remote sources, so that a path may exist on disk
// or by a simplify.net registry package.
type Source interface {
	// Definitions returns a list of file paths to JSON FHIR definitions.
	Definitions(ctx context.Context) ([]string, error)
}

// New constructs a new Source from the configuration.
func New(config *config.Package, listener RemoteListener) Source {
	if config.Path != "" {
		return NewLocalSource(config.Path)
	}
	return NewPackageSource(config.Name, config.Version, listener)
}

// NewLocalSource constructs a new local Source from the specified path.
func NewLocalSource(path string) Source {
	return localSource(path)
}

// NewPackageSource constructs a new remote Source from the specified package.
func NewPackageSource(name, version string, listener RemoteListener) Source {
	if listener == nil {
		listener = &fhirig.BaseListener{}
	}
	return &remoteSource{
		pkg:      fhirig.NewPackage(name, version),
		listener: listener,
	}
}
