package registry_test

import (
	"context"
	"embed"
	"errors"
	"testing"
	"time"

	"github.com/friendly-fhir/fhenix/pkg/registry"
	"github.com/friendly-fhir/fhenix/pkg/registry/registrytest"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

var (
	//go:embed testdata/leaf-package
	leafPackage embed.FS

	//go:embed testdata/dependent-package-one
	dependentPackageOne embed.FS

	//go:embed testdata/dependent-package-two
	dependentPackageTwo embed.FS
)

func TestDownloader(t *testing.T) {
	t.Parallel()

	const (
		registryName = "test"
		version      = "1.0.0"
	)
	testErr := errors.New("error")
	client := registrytest.NewFakeClient()
	client.SetGzipTarball("test.package", version, goodArchive)
	client.SetTarballFS("dependent.package.one", version, dependentPackageOne)
	client.SetTarballFS("dependent.package.two", version, dependentPackageTwo)
	client.SetTarballFS("leaf.package", version, leafPackage)
	client.SetError("bad.package", version, testErr)

	testCases := []struct {
		name                string
		packages            []string
		includeDependencies bool
		wantErr             error
		force               bool
		wantPackages        []string
	}{
		{
			name:         "no packages specified returns successfully",
			packages:     nil,
			wantErr:      nil,
			wantPackages: nil,
		}, {
			name:                "valid package, no dependency tracing",
			packages:            []string{"test.package"},
			includeDependencies: false,

			wantErr:      nil,
			wantPackages: []string{"test.package"},
		}, {
			name: "package does not exist",
			packages: []string{
				"test.package",
				"bad.package",
			},
			includeDependencies: false,
			wantErr:             cmpopts.AnyError,
		}, {
			name: "Package has dependencies",
			packages: []string{
				"dependent.package.one",
			},
			includeDependencies: true,
			wantErr:             nil,
			wantPackages: []string{
				"dependent.package.one",
				"leaf.package",
			},
		}, {
			name: "Packages have shared dependency",
			packages: []string{
				"dependent.package.one",
				"dependent.package.two",
			},
			includeDependencies: true,
			wantErr:             nil,
			wantPackages: []string{
				"dependent.package.one",
				"dependent.package.two",
				"leaf.package",
			},
		}, {
			name: "Force fetch",
			packages: []string{
				"dependent.package.one",
				"dependent.package.two",
			},
			includeDependencies: true,
			force:               true,
			wantErr:             nil,
			wantPackages: []string{
				"dependent.package.one",
				"dependent.package.two",
				"leaf.package",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			cache := registry.NewCache(dir)
			cache.AddClient(registryName, client.Client)

			downloader := registry.NewDownloader(cache).Force(tc.force)

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			for _, pkg := range tc.packages {
				downloader.Add(registryName, pkg, version, tc.includeDependencies)
			}

			err := downloader.Start(ctx)

			if got, want := err, tc.wantErr; !cmp.Equal(got, want, cmpopts.EquateErrors()) {
				t.Errorf("Downloader.Start() = error %v; want %v", err, tc.wantErr)
			}

			if len(tc.wantPackages) > 0 {
				for _, pkg := range tc.wantPackages {
					if _, err := cache.Get(registryName, pkg, version); err != nil {
						t.Errorf("cache.Get(%q, %q, %q) = error %v; want nil", registryName, pkg, version, err)
					}
				}
			}
		})
	}
}
