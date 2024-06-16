package fhirig_test

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"testing"
	"testing/iotest"

	"github.com/friendly-fhir/fhenix/internal/fhirig"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestParsePackage(t *testing.T) {
	testCases := []struct {
		name    string
		input   string
		want    *fhirig.Package
		wantErr error
	}{
		{
			name:  "valid package",
			input: "hl7.fhir.r4.core@4.0.1",
			want:  fhirig.NewPackage("hl7.fhir.r4.core", "4.0.1"),
		}, {
			name:    "invalid package",
			input:   "hl7.fhir.r4.core",
			wantErr: cmpopts.AnyError,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := fhirig.ParsePackage(tc.input)

			if got, want := err, tc.wantErr; !cmp.Equal(got, want, cmpopts.EquateErrors()) {
				t.Fatalf("ParsePackage(%q) = %v; want %v", tc.input, got, want)
			}
			if want := tc.want; !cmp.Equal(got, want) {
				t.Errorf("ParsePackage(%q) = %v; want %v", tc.input, got, want)
			}
		})
	}
}

func TestPackageCache_Clear(t *testing.T) {
	testCases := []struct {
		name    string
		pkg     *fhirig.Package
		wantErr error
	}{
		{
			name:    "package exists and deletes successfully",
			pkg:     fhirig.NewPackage("hl7.fhir.r4.core", "4.0.1"),
			wantErr: nil,
		}, {
			name:    "package does not exist",
			pkg:     fhirig.NewPackage("noexist.pkg.name", "420.6.9"),
			wantErr: nil,
		},
	}

	sut := &fhirig.PackageCache{
		Root: "testdata",
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			makeDummyFiles(t, "testdata/hl7.fhir.r4.core/4.0.1/package.tar.gz")

			err := sut.Clear(tc.pkg)

			if got, want := err, tc.wantErr; !cmp.Equal(got, want, cmpopts.EquateErrors()) {
				t.Errorf("PackageCache.Clear(%q) = %v; want %v", tc.pkg.String(), got, want)
			}
		})
	}
}

func TestPackageCache_Has(t *testing.T) {
	testCases := []struct {
		name string
		pkg  *fhirig.Package
		want bool
	}{
		{
			name: "package exists",
			pkg:  fhirig.NewPackage("hl7.fhir.r4.core", "4.0.1"),
			want: true,
		}, {
			name: "package does not exist",
			pkg:  fhirig.NewPackage("noexist.pkg.name", "420.6.9"),
			want: false,
		},
	}

	sut := &fhirig.PackageCache{
		Root: "testdata",
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			makeDummyFiles(t, "testdata/hl7.fhir.r4.core/4.0.1/package.tar.gz")

			got := sut.Has(tc.pkg)

			if got != tc.want {
				t.Errorf("PackageCache.Has(%q) = %v; want %v", tc.pkg.String(), got, tc.want)
			}
		})
	}
}

func TestPackageCache_Get(t *testing.T) {
	testCases := []struct {
		name    string
		pkg     *fhirig.Package
		wantErr error
	}{
		{
			name:    "package exists",
			pkg:     fhirig.NewPackage("hl7.fhir.r4.core", "4.0.1"),
			wantErr: nil,
		}, {
			name:    "package does not exist",
			pkg:     fhirig.NewPackage("noexist.pkg.name", "420.6.9"),
			wantErr: os.ErrNotExist,
		},
	}

	sut := &fhirig.PackageCache{
		Root: "testdata",
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			makeDummyFiles(t, "testdata/hl7.fhir.r4.core/4.0.1/package.tar.gz")

			_, err := sut.Get(tc.pkg)

			if got, want := err, tc.wantErr; !cmp.Equal(got, want, cmpopts.EquateErrors()) {
				t.Errorf("PackageCache.Get(%q) = %v; want %v", tc.pkg.String(), got, want)
			}
		})
	}
}

func TestPackageCache_Path(t *testing.T) {
	testCases := []struct {
		name    string
		pkg     *fhirig.Package
		want    string
		wantErr error
	}{
		{
			name:    "package exists",
			pkg:     fhirig.NewPackage("hl7.fhir.r4.core", "4.0.1"),
			want:    filepath.Join("testdata", "hl7.fhir.r4.core", "4.0.1"),
			wantErr: nil,
		}, {
			name:    "package does not exist",
			pkg:     fhirig.NewPackage("noexist.pkg.name", "420.6.9"),
			wantErr: os.ErrNotExist,
		},
	}

	sut := &fhirig.PackageCache{
		Root: "testdata",
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			makeDummyFiles(t, "testdata/hl7.fhir.r4.core/4.0.1/package.tar.gz")

			got, err := sut.Path(tc.pkg)

			if got, want := err, tc.wantErr; !cmp.Equal(got, want, cmpopts.EquateErrors()) {
				t.Fatalf("PackageCache.Path(%q) = %v; want %v", tc.pkg.String(), got, want)
			}
			if want := tc.want; got != want {
				t.Errorf("PackageCache.Path(%q) = %v; want %v", tc.pkg.String(), got, want)
			}
		})
	}
}
func makeDummyFiles(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatalf("failed to create directory: %v", err)
	}
	file, err := os.Create(path)
	if err != nil {
		t.Fatalf("failed to create file: %v", err)
	}
	defer file.Close()
}

func TestPackageCache_Fetch(t *testing.T) {
	testError := errors.New("error")
	testCases := []struct {
		name    string
		pkg     *fhirig.Package
		fetcher fhirig.Fetcher
		wantErr error
	}{
		{
			name:    "package exists",
			pkg:     fhirig.NewPackage("hl7.fhir.r4.core", "4.0.1"),
			wantErr: nil,
		}, {
			name: "package does not exist and fails to fetch",
			pkg:  fhirig.NewPackage("noexist.pkg.name", "420.6.9"),
			fetcher: fhirig.FetcherFunc(func(url string) (io.ReadCloser, error) {
				return nil, testError
			}),
			wantErr: testError,
		}, {
			name: "package does not exist and fetches bad gzip file",
			pkg:  fhirig.NewPackage("noexist.pkg.name", "420.6.9"),
			fetcher: fhirig.FetcherFunc(func(url string) (io.ReadCloser, error) {
				return io.NopCloser(iotest.ErrReader(testError)), nil
			}),
			wantErr: testError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			makeDummyFiles(t, "testdata/hl7.fhir.r4.core/4.0.1/package.tar.gz")
			sut := &fhirig.PackageCache{
				Root:    "testdata",
				Fetcher: tc.fetcher,
			}

			err := sut.Fetch(context.Background(), tc.pkg)

			if got, want := err, tc.wantErr; !cmp.Equal(got, want, cmpopts.EquateErrors()) {
				t.Errorf("PackageCache.Fetch(%q) = %v; want %v", tc.pkg.String(), got, want)
			}
		})
	}
}

func TestRealFetch(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	sut := &fhirig.PackageCache{
		Root: "testdata/real",
	}
	pkg := fhirig.NewPackage("hl7.fhir.us.core", "6.1.0")
	err := sut.Fetch(context.Background(), pkg)
	if err != nil {
		t.Logf("PackageCache.Fetch(%q) = %v; want nil", pkg.String(), err)
	}
}
