package registry_test

import (
	stdcmp "cmp"
	"io/fs"
	"path/filepath"
	"testing"

	"github.com/friendly-fhir/fhenix/pkg/registry"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestNewPackage(t *testing.T) {
	testCases := []struct {
		name    string
		path    string
		wantErr error
	}{
		{
			name: "good path",
			path: filepath.Join("testdata", "test-package"),
		}, {
			name:    "invalid package.json definition",
			path:    filepath.Join("testdata", "malformed-package"),
			wantErr: cmpopts.AnyError,
		}, {
			name:    "missing package.json",
			path:    "testdata",
			wantErr: fs.ErrNotExist,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := registry.NewPackage(tc.path)

			if got, want := err, tc.wantErr; !cmp.Equal(got, want, cmpopts.EquateErrors()) {
				t.Fatalf("NewPackage() error = %v, want %v", got, want)
			}
		})
	}
}

func TestNewPackage_Files(t *testing.T) {
	pkg, err := registry.NewPackage(filepath.Join("testdata", "test-package"))
	if err != nil {
		t.Fatalf("NewPackage() error = %v", err)
	}
	want := []string{
		filepath.Join("testdata", "test-package", "StructureDefinition-foo.json"),
	}

	files, err := pkg.Files()
	if err != nil {
		t.Fatalf("Package.Files() error = %v", err)
	}

	if got := files; !cmp.Equal(got, want, cmpopts.SortSlices(stdcmp.Less[string])) {
		t.Errorf("Package.Files() = %v, want %v", got, want)
	}
}

func TestNewPackage_Name(t *testing.T) {
	testCases := []struct {
		name     string
		manifest *registry.PackageManifest
		want     string
	}{
		{
			name: "good manifest",
			manifest: &registry.PackageManifest{
				Name: "test-package",
			},
			want: "test-package",
		}, {
			name: "nil manifest",
			want: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			pkg := &registry.Package{
				Manifest: tc.manifest,
			}

			if got, want := pkg.Name(), tc.want; got != want {
				t.Errorf("Package.Name() = %q, want %q", got, want)
			}
		})
	}
}

func TestNewPackage_Version(t *testing.T) {
	testCases := []struct {
		name     string
		manifest *registry.PackageManifest
		want     string
	}{
		{
			name: "good manifest",
			manifest: &registry.PackageManifest{
				Version: "1.0.0",
			},
			want: "1.0.0",
		}, {
			name: "nil manifest",
			want: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			pkg := &registry.Package{
				Manifest: tc.manifest,
			}

			if got, want := pkg.Version(), tc.want; got != want {
				t.Errorf("Package.Version() = %q, want %q", got, want)
			}
		})
	}
}

func TestNewPackage_Dependencies(t *testing.T) {
	testCases := []struct {
		name     string
		manifest *registry.PackageManifest
		want     map[string]string
	}{
		{
			name: "good manifest",
			manifest: &registry.PackageManifest{
				Dependencies: map[string]string{
					"foo": "1.0.0",
				},
			},
			want: map[string]string{
				"foo": "1.0.0",
			},
		}, {
			name: "nil manifest",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			pkg := &registry.Package{
				Manifest: tc.manifest,
			}

			if got, want := pkg.Dependencies(), tc.want; !cmp.Equal(got, want, cmpopts.EquateEmpty()) {
				t.Errorf("Package.Dependencies() = %q, want %q", got, want)
			}
		})
	}
}

func TestNewPackage_FHIRVersions(t *testing.T) {
	testCases := []struct {
		name     string
		manifest *registry.PackageManifest
		want     []string
	}{
		{
			name: "good manifest",
			manifest: &registry.PackageManifest{
				FHIRVersionList: []string{
					"4.0.0",
				},
			},
			want: []string{
				"4.0.0",
			},
		}, {
			name: "nil manifest",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			pkg := &registry.Package{
				Manifest: tc.manifest,
			}

			if got, want := pkg.FHIRVersionList(), tc.want; !cmp.Equal(got, want, cmpopts.EquateEmpty()) {
				t.Errorf("Package.FHIRVersionList() = %q, want %q", got, want)
			}
		})
	}
}
