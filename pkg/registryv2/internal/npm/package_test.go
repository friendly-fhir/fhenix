package npm_test

import (
	"io"
	"io/fs"
	"os"
	"strings"
	"testing"

	"github.com/friendly-fhir/fhenix/pkg/registryv2/internal/npm"
	"github.com/friendly-fhir/fhenix/pkg/registryv2/internal/npm/npmtest"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestReadManifest(t *testing.T) {
	testCases := []struct {
		name    string
		content string
		want    *npm.PackageManifest
		wantErr error
	}{
		{
			name: "valid manifest",
			content: `{
				"name": "foo",
				"version": "1.0.0"
		  }`,
			want: &npm.PackageManifest{
				Name:    "foo",
				Version: "1.0.0",
			},
		}, {
			name:    "invalid manifest",
			content: `{`,
			wantErr: cmpopts.AnyError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := npm.ReadManifest(strings.NewReader(tc.content))

			if got, want := err, tc.wantErr; !cmp.Equal(got, want, cmpopts.EquateErrors()) {
				t.Fatalf("ReadManifest() err = %v, want %v", got, want)
			}
			if got, want := got, tc.want; !cmp.Equal(got, want) {
				t.Fatalf("ReadManifest() = %v, want %v", got, want)
			}
		})
	}
}

func TestBasePackageVisitor_VisitManifest(t *testing.T) {
	sut := &npm.BasePackageVisitor{}

	err := sut.VisitManifest(&npm.PackageManifest{
		Name:    "foo",
		Version: "1.0.0",
	})

	if err != nil {
		t.Fatalf("VisitManifest() err = %v, want nil", err)
	}
}

func TestBasePackageVisitor_VisitFile(t *testing.T) {
	sut := &npm.BasePackageVisitor{}

	err := sut.VisitFile("foo", nil, nil)

	if err != nil {
		t.Fatalf("VisitFile() err = %v, want nil", err)
	}
}

func TestVisitPackageStream(t *testing.T) {
	reader := npmtest.SimplePackageTarball()
	defer reader.Close()
	wantManifest := &npm.PackageManifest{
		Name:         "package",
		Version:      "1.0.0",
		License:      "ISC",
		Dependencies: map[string]string{},
	}
	wantFiles := npmtest.FSContents(t, npmtest.SimplePackageFS())
	v := &testPackageVisitor{
		files: make(map[string]string),
	}

	err := npm.VisitPackageStream(reader, v)

	if got, want := err, error(nil); !cmp.Equal(got, want, cmpopts.EquateErrors()) {
		t.Fatalf("VisitPackageStream() err = %v, want %v", got, want)
	}
	if diff := cmp.Diff(v.manifest, wantManifest); diff != "" {
		t.Fatalf("VisitPackageStream() manifest (-got +want):\n%s", diff)
	}
	if diff := cmp.Diff(v.files, wantFiles); diff != "" {
		t.Fatalf("VisitPackageStream() files (-got +want):\n%s", diff)
	}
}

func TestUnpackPackageVisitor(t *testing.T) {
	tmp := t.TempDir()
	reader := npmtest.SimplePackageTarball()
	defer reader.Close()
	wantFiles := npmtest.FSContents(t, npmtest.SimplePackageFS())
	sut := &npm.UnpackPackageVisitor{
		Destination: tmp,
	}

	err := npm.VisitPackageStream(reader, sut)
	got := npmtest.FSContents(t, os.DirFS(tmp))

	if got, want := err, error(nil); !cmp.Equal(got, want, cmpopts.EquateErrors()) {
		t.Fatalf("UnpackPackageVisitor() err = %v, want %v", got, want)
	}
	if diff := cmp.Diff(got, wantFiles); diff != "" {
		t.Fatalf("UnpackPackageVisitor() files (-got +want):\n%s", diff)
	}
}

type testPackageVisitor struct {
	manifest *npm.PackageManifest
	files    map[string]string

	npm.BasePackageVisitor
}

func (v *testPackageVisitor) VisitManifest(manifest *npm.PackageManifest) error {
	v.manifest = manifest
	return nil
}

func (v *testPackageVisitor) VisitFile(path string, _ fs.FileInfo, r io.Reader) error {
	b, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	v.files[path] = string(b)
	return nil
}
