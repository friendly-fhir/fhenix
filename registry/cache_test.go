package registry_test

import (
	"context"
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/friendly-fhir/fhenix/registry"
	"github.com/friendly-fhir/fhenix/registry/registrytest"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestCache_Fetch(t *testing.T) {
	testErr := errors.New("test error")
	client := registrytest.NewFakeClient()
	client.SetOK("test.package", "1.0.0", goodArchive)
	client.SetOK("test.package", "2.0.0", goodArchive)
	client.SetError("fail.package", "1.0.0", testErr)

	testCases := []struct {
		name     string
		registry string
		pkg      string
		version  string
		wantErr  error
	}{
		{
			name:     "good package",
			registry: "test",
			pkg:      "test.package",
			version:  "1.0.0",
		}, {
			name:     "bad package registry",
			registry: "bad",
			pkg:      "test.package",
			version:  "1.0.0",
			wantErr:  cmpopts.AnyError,
		}, {
			name:     "bad package name",
			registry: "test",
			pkg:      "fail.package",
			version:  "1.0.0",
			wantErr:  testErr,
		}, {
			name:     "cached package",
			registry: "test",
			pkg:      "test.package",
			version:  "2.0.0",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cache := registry.NewCache(t.TempDir())
			cache.AddClient("test", client.Client)
			if err := cache.ForceFetch(context.Background(), "test", "test.package", "2.0.0"); err != nil {
				t.Fatalf("Cache.ForceFetch() error = %v", err)
			}

			err := cache.Fetch(context.Background(), tc.registry, tc.pkg, tc.version)

			if got, want := err, tc.wantErr; !cmp.Equal(got, want, cmpopts.EquateErrors()) {
				t.Fatalf("Cache.Fetch() error = %v, want %v", got, want)
			}
		})
	}
}

func TestCache_ForceFetch(t *testing.T) {
	testErr := errors.New("test error")
	client := registrytest.NewFakeClient()
	client.SetOK("test.package", "1.0.0", goodArchive)
	client.SetError("fail.package", "1.0.0", testErr)

	testCases := []struct {
		name     string
		registry string
		pkg      string
		version  string
		wantErr  error
	}{
		{
			name:     "good package",
			registry: "test",
			pkg:      "test.package",
			version:  "1.0.0",
		}, {
			name:     "bad package registry",
			registry: "bad",
			pkg:      "test.package",
			version:  "1.0.0",
			wantErr:  cmpopts.AnyError,
		}, {
			name:     "bad package name",
			registry: "test",
			pkg:      "fail.package",
			version:  "1.0.0",
			wantErr:  testErr,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cache := registry.NewCache(t.TempDir())
			cache.AddClient("test", client.Client)

			err := cache.ForceFetch(context.Background(), tc.registry, tc.pkg, tc.version)

			if got, want := err, tc.wantErr; !cmp.Equal(got, want, cmpopts.EquateErrors()) {
				t.Fatalf("Cache.ForceFetch() error = %v, want %v", got, want)
			}
		})
	}
}

func readManifest(t *testing.T, path string) *registry.PackageManifest {
	t.Helper()
	var manifest registry.PackageManifest

	file, err := os.Open(path)
	if err != nil {
		t.Fatalf("ReadManifest() error = %v", err)
	}
	defer file.Close()
	if err := json.NewDecoder(file).Decode(&manifest); err != nil {
		t.Fatalf("ReadManifest() error = %v", err)
	}
	return &manifest
}

func TestCache_Get(t *testing.T) {
	testErr := errors.New("test error")
	client := registrytest.NewFakeClient()
	client.SetOK("test.package", "1.0.0", goodArchive)
	client.SetOK("test.package", "2.0.0", goodArchive)
	client.SetError("fail.package", "1.0.0", testErr)

	testCases := []struct {
		name     string
		registry string
		pkg      string
		version  string
		want     *registry.PackageManifest
		wantErr  error
	}{
		{
			name:     "package is not fetched",
			registry: "test",
			pkg:      "test.package",
			version:  "1.0.0",
			wantErr:  cmpopts.AnyError,
		}, {
			name:     "package is cached",
			registry: "test",
			pkg:      "test.package",
			version:  "2.0.0",
			want:     readManifest(t, filepath.Join("testdata", "test-package", "package.json")),
		}, {
			name:     "bad package registry",
			registry: "bad",
			pkg:      "test.package",
			version:  "1.0.0",
			wantErr:  cmpopts.AnyError,
		}, {
			name:     "bad package name",
			registry: "test",
			pkg:      "fail.package",
			version:  "1.0.0",
			wantErr:  fs.ErrNotExist,
		}, {
			name:    "No registry name",
			pkg:     "test.package",
			version: "1.0.0",
			wantErr: cmpopts.AnyError,
		}, {
			name:     "No package name",
			registry: "test",
			version:  "1.0.0",
			wantErr:  cmpopts.AnyError,
		}, {
			name:     "No version",
			registry: "test",
			pkg:      "test.package",
			wantErr:  cmpopts.AnyError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			cache := registry.NewCache(dir)
			cache.AddClient("test", client.Client)
			if err := cache.ForceFetch(context.Background(), "test", "test.package", "2.0.0"); err != nil {
				t.Fatalf("Cache.ForceFetch() error = %v", err)
			}

			got, err := cache.Get(tc.registry, tc.pkg, tc.version)

			var want *registry.Package
			if tc.want != nil {
				want = &registry.Package{
					Path:     cache.CacheDir(tc.registry, tc.pkg, tc.version),
					Manifest: tc.want,
				}
			}
			if got, want := err, tc.wantErr; !cmp.Equal(got, want, cmpopts.EquateErrors()) {
				t.Fatalf("Cache.Get(%q, %q, %q) error = %v, want %v", tc.registry, tc.pkg, tc.version, got, want)
			}
			if got, want := got, want; !cmp.Equal(got, want) {
				t.Errorf("Cache.Get(%q, %q, %q) = %v, want %v", tc.registry, tc.pkg, tc.version, got, want)
			}
		})
	}
}

func TestCache_Get_InvokesListener(t *testing.T) {
	const (
		registryName = "test"
		pkg          = "test.package"
		version      = "1.0.0"
	)
	client := registrytest.NewFakeClient()
	client.SetOK(pkg, version, goodArchive)

	listener := registrytest.NewCacheListener()

	cache := registry.NewCache(t.TempDir())
	cache.AddClient(registryName, client.Client)
	cache.AddListener(listener)

	if err := cache.Fetch(context.Background(), registryName, pkg, version); err != nil {
		t.Fatalf("Cache.ForceFetch() error = %v", err)
	}

	if got, want := listener.CacheHits(registryName, pkg, version), 0; got != want {
		t.Errorf("listener.CacheHits(%q, %q, %q) = %d, want %d", registryName, pkg, version, got, want)
	}
	if got, want := listener.FetchCalls(registryName, pkg, version), 1; got != want {
		t.Errorf("listener.FetchCalls(%q, %q, %q) = %d, want %d", registryName, pkg, version, got, want)
	}
	if got, want := listener.FetchBytes(registryName, pkg, version), int64(len(goodArchive)); got != want {
		t.Errorf("listener.FetchBytes(%q, %q, %q) = %d, want %d", registryName, pkg, version, got, want)
	}
	if got, want := listener.UnpackCalls(registryName, pkg, version), 2; got != want {
		t.Errorf("listener.UnpackCalls(%q, %q, %q) = %d, want %d", registryName, pkg, version, got, want)
	}

}

func TestCache_GetOrFetch(t *testing.T) {
	testErr := errors.New("test error")
	client := registrytest.NewFakeClient()
	client.SetOK("test.package", "1.0.0", goodArchive)
	client.SetOK("test.package", "2.0.0", goodArchive)
	client.SetError("fail.package", "1.0.0", testErr)

	testCases := []struct {
		name     string
		registry string
		pkg      string
		version  string
		want     *registry.PackageManifest
		wantErr  error
	}{
		{
			name:     "fetch package that is not fetched",
			registry: "test",
			pkg:      "test.package",
			version:  "1.0.0",
			want:     readManifest(t, filepath.Join("testdata", "test-package", "package.json")),
		}, {
			name:     "use cached package",
			registry: "test",
			pkg:      "test.package",
			version:  "2.0.0",
			want:     readManifest(t, filepath.Join("testdata", "test-package", "package.json")),
		}, {
			name:     "bad package registry",
			registry: "bad",
			pkg:      "test.package",
			version:  "1.0.0",
			wantErr:  cmpopts.AnyError,
		}, {
			name:     "bad package name",
			registry: "test",
			pkg:      "fail.package",
			version:  "1.0.0",
			wantErr:  testErr,
		}, {
			name:    "No registry name",
			pkg:     "test.package",
			version: "1.0.0",
			wantErr: cmpopts.AnyError,
		}, {
			name:     "No package name",
			registry: "test",
			version:  "1.0.0",
			wantErr:  cmpopts.AnyError,
		}, {
			name:     "No version",
			registry: "test",
			pkg:      "test.package",
			wantErr:  cmpopts.AnyError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			cache := registry.NewCache(dir)
			cache.AddClient("test", client.Client)
			if err := cache.ForceFetch(context.Background(), "test", "test.package", "2.0.0"); err != nil {
				t.Fatalf("Cache.ForceFetch() error = %v", err)
			}

			got, err := cache.GetOrFetch(context.Background(), tc.registry, tc.pkg, tc.version)

			var want *registry.Package
			if tc.want != nil {
				want = &registry.Package{
					Path:     cache.CacheDir(tc.registry, tc.pkg, tc.version),
					Manifest: tc.want,
				}
			}
			if got, want := err, tc.wantErr; !cmp.Equal(got, want, cmpopts.EquateErrors()) {
				t.Fatalf("Cache.GetOrFetch() error = %v, want %v", got, want)
			}
			if got, want := got, want; !cmp.Equal(got, want) {
				t.Errorf("Cache.GetOrFetch() = %v, want %v", got, want)
			}
		})
	}
}

func TestCache_Delete(t *testing.T) {
	client := registrytest.NewFakeClient()
	client.SetOK("test.package", "2.0.0", goodArchive)

	testCases := []struct {
		name     string
		registry string
		pkg      string
		version  string
		wantErr  error
	}{
		{
			name:     "use cached package",
			registry: "test",
			pkg:      "test.package",
			version:  "2.0.0",
		}, {
			name:     "bad package registry",
			registry: "bad",
			pkg:      "test.package",
			version:  "1.0.0",
			wantErr:  cmpopts.AnyError,
		}, {
			name:    "No registry name",
			pkg:     "test.package",
			version: "1.0.0",
			wantErr: cmpopts.AnyError,
		}, {
			name:     "No package name",
			registry: "test",
			version:  "1.0.0",
			wantErr:  cmpopts.AnyError,
		}, {
			name:     "No version",
			registry: "test",
			pkg:      "test.package",
			wantErr:  cmpopts.AnyError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			cache := registry.NewCache(dir)
			cache.AddClient("test", client.Client)
			if err := cache.ForceFetch(context.Background(), "test", "test.package", "2.0.0"); err != nil {
				t.Fatalf("Cache.ForceFetch() error = %v", err)
			}

			err := cache.Delete(tc.registry, tc.pkg, tc.version)

			if got, want := err, tc.wantErr; !cmp.Equal(got, want, cmpopts.EquateErrors()) {
				t.Fatalf("Cache.Delete() error = %v, want %v", got, want)
			}
		})
	}
}

func setEnv(name, value string) func(t *testing.T) {
	return func(t *testing.T) {
		t.Helper()
		t.Setenv(name, value)
	}
}

func setHome(value string) func(t *testing.T) {
	return func(t *testing.T) {
		if runtime.GOOS == "windows" {
			t.Setenv("USERPROFILE", value)
		} else {
			t.Setenv("HOME", value)
		}
	}
}

func TestDefaultCache(t *testing.T) {
	dir := t.TempDir()
	testCases := []struct {
		name  string
		setup func(*testing.T)
		want  string
	}{
		{
			name:  "FHIR_CACHE",
			setup: setEnv("FHIR_CACHE", dir),
			want:  dir,
		}, {
			name:  "User home directory",
			setup: setHome(dir),
			want:  filepath.Join(dir, ".fhir"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup(t)

			cache := registry.DefaultCache()

			if got, want := cache.Root(), tc.want; got != want {
				t.Errorf("DefaultCache().Root() = %q, want %q", got, want)
			}
		})
	}
}
