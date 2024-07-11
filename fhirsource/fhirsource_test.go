package fhirsource_test

import (
	"context"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
	"testing"

	"github.com/friendly-fhir/fhenix/fhirsource"
	"github.com/friendly-fhir/fhenix/internal/fhirig"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

type result[T any] struct {
	Value T
	Err   error
}

type FakePackageCache struct {
	fetch        map[string]*result[[]string]
	dependencies map[string]*result[[]*fhirig.Package]
}

func NewFakePackageCache() *FakePackageCache {
	return &FakePackageCache{
		fetch:        map[string]*result[[]string]{},
		dependencies: map[string]*result[[]*fhirig.Package]{},
	}
}

func (fpc *FakePackageCache) SetFetch(pkg *fhirig.Package, files []string, err error) {
	fpc.fetch[pkg.String()] = &result[[]string]{
		Value: files,
		Err:   err,
	}
}

func (fpc *FakePackageCache) SetDependencies(pkg *fhirig.Package, deps []*fhirig.Package, err error) {
	fpc.dependencies[pkg.String()] = &result[[]*fhirig.Package]{
		Value: deps,
		Err:   err,
	}
}

func (fpc *FakePackageCache) Dependencies(pkg *fhirig.Package) ([]*fhirig.Package, error) {
	got, ok := fpc.dependencies[pkg.String()]
	if !ok || got == nil {
		return nil, fmt.Errorf("dependencies: %q requested but not configured", pkg.String())
	}
	return got.Value, got.Err
}

func (fpc *FakePackageCache) FetchAndGet(_ context.Context, pkg *fhirig.Package) ([]string, error) {
	got, ok := fpc.fetch[pkg.String()]
	if !ok || got == nil {
		return nil, fmt.Errorf("fetch: %q requested but not configured", pkg.String())
	}
	return got.Value, got.Err
}

func TestLocalSourceBundles(t *testing.T) {
	testCases := []struct {
		name    string
		pkg     *fhirig.Package
		path    string
		want    []*fhirsource.Bundle
		wantErr error
	}{
		{
			name: "valid local source",
			path: "testdata",
			pkg:  fhirig.NewPackage("test", "1.0.0"),
			want: []*fhirsource.Bundle{
				{
					Package: fhirig.NewPackage("test", "1.0.0"),
					Files: []string{
						filepath.Join("testdata", "file-a.json"),
						filepath.Join("testdata", "file-b.json"),
						filepath.Join("testdata", "dir", "file-c.json"),
					},
				},
			},
		}, {
			name:    "invalid local source",
			path:    "testdata/missing",
			pkg:     fhirig.NewPackage("test", "1.0.0"),
			wantErr: fs.ErrNotExist,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sut := fhirsource.NewLocalSource(tc.pkg, tc.path)

			got, err := sut.Bundles(context.Background())

			if got, want := err, tc.wantErr; !cmp.Equal(got, want, cmpopts.EquateErrors()) {
				t.Fatalf("LocalSource.Bundles(...) = error %v; want %v", got, tc.want)
			}
			less := func(a, b string) bool { return strings.Compare(a, b) < 0 }
			if got, want := got, tc.want; !cmp.Equal(got, want, cmpopts.SortSlices(less)) {
				t.Errorf("LocalSource.Bundles(...) = %v; want %v", got, tc.want)
			}
		})
	}
}

func TestRemoteSourceBundles(t *testing.T) {
	package1 := &fhirsource.Bundle{
		Package: fhirig.NewPackage("package-1", "1.0.0"),
		Files: []string{
			filepath.Join("testdata", "file-a.json"),
			filepath.Join("testdata", "file-b.json"),
			filepath.Join("testdata", "dir", "file-c.json"),
		},
	}
	package2 := &fhirsource.Bundle{
		Package: fhirig.NewPackage("package-2", "2.0.0"),
		Files: []string{
			filepath.Join("testdata", "file-d.json"),
			filepath.Join("testdata", "file-e.json"),
			filepath.Join("testdata", "dir", "file-f.json"),
		},
	}
	package3 := &fhirsource.Bundle{
		Package: fhirig.NewPackage("package-3", "3.0.0"),
		Files: []string{
			filepath.Join("testdata", "file-g.json"),
		},
	}
	err := fmt.Errorf("test error")

	testCases := []struct {
		name    string
		input   *fhirig.Package
		setup   func() *FakePackageCache
		want    []*fhirsource.Bundle
		wantErr error
	}{
		{
			name:  "valid remote source",
			input: package1.Package,
			setup: func() *FakePackageCache {
				fake := NewFakePackageCache()
				fake.SetDependencies(package1.Package, nil, nil)
				fake.SetFetch(package1.Package, package1.Files, nil)
				return fake
			},
			want: []*fhirsource.Bundle{package1},
		}, {
			name:  "valid remote source with dependencies",
			input: package1.Package,
			setup: func() *FakePackageCache {
				fake := NewFakePackageCache()
				fake.SetDependencies(package1.Package, []*fhirig.Package{package2.Package, package3.Package}, nil)
				fake.SetFetch(package1.Package, package1.Files, nil)
				fake.SetFetch(package2.Package, package2.Files, nil)
				fake.SetFetch(package3.Package, package3.Files, nil)
				return fake
			},
			want: []*fhirsource.Bundle{package1, package2, package3},
		}, {
			name:  "dependencies fails",
			input: package1.Package,
			setup: func() *FakePackageCache {
				fake := NewFakePackageCache()
				fake.SetDependencies(package1.Package, nil, err)
				return fake
			},
			wantErr: err,
		}, {
			name:  "fetch fails",
			input: package1.Package,
			setup: func() *FakePackageCache {
				fake := NewFakePackageCache()
				fake.SetDependencies(package1.Package, nil, nil)
				fake.SetFetch(package1.Package, nil, err)
				return fake
			},
			wantErr: err,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fake := tc.setup()
			sut := fhirsource.NewRemoteSource(tc.input, fake)

			got, err := sut.Bundles(context.Background())

			if got, want := err, tc.wantErr; !cmp.Equal(got, want, cmpopts.EquateErrors()) {
				t.Fatalf("RemoteSource.Bundles(...) = error %v; want %v", got, tc.want)
			}

			less := func(a, b string) bool { return strings.Compare(a, b) < 0 }
			if got, want := got, tc.want; !cmp.Equal(got, want, cmpopts.SortSlices(less)) {
				t.Errorf("RemoteSource.Bundles(...) = %v; want %v", got, tc.want)
			}
		})
	}
}
