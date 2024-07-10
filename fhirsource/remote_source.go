package fhirsource

import (
	"context"

	"github.com/friendly-fhir/fhenix/internal/fhirig"
)

type Package = fhirig.Package

type RemoteListener interface {
	OnFetchStart(pkg *Package)
	OnFetchEnd(pkg *Package, err error)
	OnCacheHit(pkg *Package)
}

type remoteSource struct {
	pkg      *fhirig.Package
	listener RemoteListener
}

func (rs *remoteSource) Definitions(ctx context.Context) ([]string, error) {
	cache := fhirig.NewSystemCache()
	cache.Listener = rs.listener
	if !cache.Has(rs.pkg) {
		if err := cache.Fetch(ctx, rs.pkg); err != nil {
			return nil, err
		}
	}
	return cache.Get(rs.pkg)
}
