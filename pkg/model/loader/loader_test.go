package loader_test

import (
	"context"
	_ "embed"
	"testing"

	"github.com/friendly-fhir/fhenix/pkg/model/conformance"
	"github.com/friendly-fhir/fhenix/pkg/model/loader"
	"github.com/friendly-fhir/fhenix/pkg/registry"
	"github.com/friendly-fhir/fhenix/pkg/registry/registrytest"
)

var (
	//go:embed testdata/leaf-package/package.tar.gz
	leafArchive []byte

	//go:embed testdata/dependent-package-one/package.tar.gz
	dependentArchiveOne []byte

	//go:embed testdata/dependent-package-two/package.tar.gz
	dependentArchiveTwo []byte
)

type TestListener struct {
	loader.BaseListener
}

func (TestListener) OnLoad(ref registry.PackageRef) {

}

func TestLoader_Load(t *testing.T) {
	const (
		registryName = "test"
		version      = "1.0.0"
	)
	var listener TestListener
	client := registrytest.NewFakeClient()
	client.SetGzipTarball("dependent.package.one", version, dependentArchiveOne)
	client.SetGzipTarball("dependent.package.two", version, dependentArchiveTwo)
	client.SetGzipTarball("leaf.package", version, leafArchive)

	cache := registry.NewCache(t.TempDir())
	cache.AddClient(registryName, client.Client)

	downloader := registry.NewDownloader(cache)
	downloader.Add(registryName, "dependent.package.one", version, true)
	downloader.Add(registryName, "dependent.package.two", version, true)
	if err := downloader.Start(context.Background()); err != nil {
		t.Fatalf("downloader.Start() = %v, want nil", err)
	}

	module := conformance.DefaultModule()

	loader := loader.New(cache,
		loader.WithModule(module),
		loader.WithWorkers(4),
		loader.WithListeners(listener),
	)

	ref1 := registry.NewPackageRef(registryName, "dependent.package.one", version)
	ref2 := registry.NewPackageRef(registryName, "dependent.package.two", version)

	err := loader.Load(ref1, ref2)

	if err != nil {
		t.Errorf("loader.Load() = %v, want nil", err)
	}
}
