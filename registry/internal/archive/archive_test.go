package archive_test

import (
	"bytes"
	stdcmp "cmp"
	_ "embed"
	"io"
	"strings"
	"testing"

	"github.com/friendly-fhir/fhenix/registry/internal/archive"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

var (
	//go:embed testdata/good-package.tar.gz
	goodArchive []byte

	//go:embed testdata/broken-package.tar.gz
	brokenArchive []byte

	//go:embed testdata/uncompressed-package.tar
	uncompressedArchive []byte
)

func TestVisitFiles(t *testing.T) {
	testCases := []struct {
		name      string
		contents  []byte
		filter    func(s string) bool
		transform func(s string) string
		want      []string
		wantErr   error
	}{
		{
			name:     "good package",
			contents: goodArchive,
			want: []string{
				"good-package/package.json",
				"good-package/StructureDefinition-foo.json",
			},
		}, {
			name:     "good package with filter",
			contents: goodArchive,
			filter: func(s string) bool {
				return s == "good-package/package.json"
			},
			want: []string{
				"good-package/package.json",
			},
		}, {
			name:     "good package with transform",
			contents: goodArchive,
			transform: func(s string) string {
				return strings.TrimPrefix(s, "good-package/")
			},
			want: []string{
				"package.json",
				"StructureDefinition-foo.json",
			},
		}, {
			name:     "broken package",
			contents: brokenArchive,
			wantErr:  cmpopts.AnyError,
		}, {
			name:     "uncompressed package",
			contents: uncompressedArchive,
			wantErr:  cmpopts.AnyError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var opts []archive.Option
			if tc.filter != nil {
				opts = append(opts, archive.Filter(tc.filter))
			}
			if tc.transform != nil {
				opts = append(opts, archive.Transform(tc.transform))
			}
			sut := archive.New(bytes.NewReader(tc.contents), opts...)

			var got []string
			err := sut.Unpack(archive.VisitorFunc(func(s string, _ io.Reader) error {
				got = append(got, s)
				return nil
			}))

			if got, want := err, tc.wantErr; !cmp.Equal(got, want, cmpopts.EquateErrors()) {
				t.Fatalf("VisitFiles() = error %v, want %v", got, want)
			}
			if want := tc.want; !cmp.Equal(got, want, cmpopts.SortSlices(stdcmp.Less[string])) {
				t.Errorf("VisitFiles() = %v, want %v", got, want)
			}
		})
	}
}
