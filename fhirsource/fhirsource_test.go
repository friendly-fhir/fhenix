package fhirsource_test

import (
	"context"
	"io/fs"
	"path/filepath"
	"strings"
	"testing"

	"github.com/friendly-fhir/fhenix/config"
	"github.com/friendly-fhir/fhenix/fhirsource"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestSourceDefinitions(t *testing.T) {
	testCases := []struct {
		name    string
		input   *config.Package
		want    []string
		wantErr error
	}{
		{
			name:  "valid local source",
			input: &config.Package{Path: "testdata"},
			want: []string{
				filepath.Join("testdata", "file-a.json"),
				filepath.Join("testdata", "file-b.json"),
				filepath.Join("testdata", "dir", "file-c.json"),
			},
		}, {
			name:    "invalid local source",
			input:   &config.Package{Path: "testdata/missing"},
			wantErr: fs.ErrNotExist,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sut := fhirsource.New(tc.input, nil)

			got, err := sut.Definitions(context.Background())

			if got, want := err, tc.wantErr; !cmp.Equal(got, want, cmpopts.EquateErrors()) {
				t.Fatalf("Source.Definitions(...) = error %v; want %v", got, tc.want)
			}
			less := func(a, b string) bool { return strings.Compare(a, b) < 0 }
			if got, want := got, tc.want; !cmp.Equal(got, want, cmpopts.SortSlices(less)) {
				t.Errorf("Source.Definitions(...) = %v; want %v", got, tc.want)
			}
		})
	}
}
