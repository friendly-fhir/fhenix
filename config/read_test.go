package config_test

import (
	"io/fs"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/friendly-fhir/fhenix/config"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func thisdir(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	return filepath.Dir(file)
}

func TestFromFile(t *testing.T) {
	testCases := []struct {
		name    string
		input   string
		want    *config.Config
		wantErr error
	}{
		{
			name:  "valid config",
			input: "testdata/v1.yaml",
			want: &config.Config{
				Mode:      config.Mode("text"),
				OutputDir: thisdir(t),
				Transforms: []*config.Transform{
					{
						Include: []*config.TransformFilter{
							{
								Name: ".*",
								Type: "StructureDefinition",
							},
						},
						Funcs: map[string]string{
							"receiver": filepath.Join(thisdir(t), "testdata", "templates", "funcs", "receiver.tmpl"),
						},
						Templates: map[string]string{
							"header": filepath.Join(thisdir(t), "testdata", "templates", "header.tmpl"),
							"footer": filepath.Join(thisdir(t), "testdata", "templates", "footer.tmpl"),
							"type":   filepath.Join(thisdir(t), "testdata", "templates", "type.tmpl"),
							"custom": filepath.Join(thisdir(t), "testdata", "templates", "custom.tmpl"),
						},
					},
				},
				Input: &config.Package{
					Name:                "hl7.fhir.r4.core",
					Version:             "4.0.1",
					IncludeDependencies: true,
				},
			},
		}, {
			name:    "file does not exist",
			input:   "testdata/does-not-exist.yaml",
			wantErr: fs.ErrNotExist,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := config.FromFile(tc.input)
			if got, want := err, tc.wantErr; !cmp.Equal(got, want, cmpopts.EquateErrors()) {
				t.Fatalf("FromFile(%q) = %v, want %v", tc.input, got, want)
			}

			if want := tc.want; !cmp.Equal(got, want, cmpopts.EquateEmpty()) {
				diff := cmp.Diff(got, want, cmpopts.EquateEmpty())
				t.Errorf("FromFile(%q) (-got +want)\n%s", tc.input, diff)
			}
		})
	}
}
