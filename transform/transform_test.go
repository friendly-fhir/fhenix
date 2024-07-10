package transform_test

import (
	"io/fs"
	"testing"

	"github.com/friendly-fhir/fhenix/config"
	"github.com/friendly-fhir/fhenix/model"
	"github.com/friendly-fhir/fhenix/transform"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestNew(t *testing.T) {
	testCases := []struct {
		name      string
		mode      config.Mode
		input     *config.Transform
		wantValue bool
		wantErr   error
	}{
		{
			name:    "bad mode returns error",
			mode:    config.Mode("invalid"),
			input:   &config.Transform{},
			wantErr: cmpopts.AnyError,
		}, {
			name: "bad function path",
			mode: config.Mode("text"),
			input: &config.Transform{
				Funcs: map[string]string{
					"no_exist": "testdata/does-not-exist.tmpl",
				},
			},
			wantErr: fs.ErrNotExist,
		}, {
			name: "bad function definition",
			mode: config.Mode("text"),
			input: &config.Transform{
				Funcs: map[string]string{
					"bad": "testdata/bad.tmpl",
				},
			},
			wantErr: cmpopts.AnyError,
		}, {
			name: "bad template path",
			mode: config.Mode("text"),
			input: &config.Transform{
				Templates: map[string]string{
					"no_exist": "testdata/does-not-exist.tmpl",
				},
			},
			wantErr: fs.ErrNotExist,
		}, {
			name: "bad template definition",
			mode: config.Mode("text"),
			input: &config.Transform{
				Templates: map[string]string{
					"bad": "testdata/bad.tmpl",
				},
			},
			wantErr: cmpopts.AnyError,
		}, {
			name: "bad output path",
			mode: config.Mode("text"),
			input: &config.Transform{
				OutputPath: "{{- . | calls-bad-function-def }}",
			},
			wantErr: cmpopts.AnyError,
		}, {
			name: "valid text mode",
			mode: config.Mode("text"),
			input: &config.Transform{
				Funcs: map[string]string{
					"first": "testdata/first.tmpl",
				},
				Templates: map[string]string{
					"header": "testdata/header.tmpl",
					"footer": "testdata/footer.tmpl",
					"main":   "testdata/body.tmpl",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := transform.New(tc.mode, tc.input)

			if got, want := err, tc.wantErr; !cmp.Equal(got, want, cmpopts.EquateErrors()) {
				t.Fatalf("New() = error %v, want %v", got, want)
			}
			if got == nil && tc.wantValue {
				t.Errorf("New() = nil, want non-nil")
			}
		})
	}
}

func TestTransformCanTransform(t *testing.T) {
	testCases := []struct {
		name   string
		config *config.Transform
		input  any
		want   bool
	}{
		{
			name: "type matches one include filter",
			config: &config.Transform{
				Include: []*config.TransformFilter{
					{Type: "CodeSystem"},
					{Type: "StructureDefinition"},
				},
			},
			input: &model.Type{},
			want:  true,
		}, {
			name: "type matches one exclude filter",
			config: &config.Transform{
				Exclude: []*config.TransformFilter{
					{Type: "CodeSystem"},
					{Type: "StructureDefinition"},
				},
			},
			input: &model.Type{},
			want:  false,
		}, {
			name: "type matches both include and exclude filter",
			config: &config.Transform{
				Include: []*config.TransformFilter{
					{Type: "StructureDefinition"},
				},
				Exclude: []*config.TransformFilter{
					{Type: "StructureDefinition"},
				},
			},
			input: &model.Type{},
			want:  false,
		}, {
			name: "type matches include filter and not exclude filter",
			config: &config.Transform{
				Include: []*config.TransformFilter{
					{Type: "StructureDefinition"},
				},
				Exclude: []*config.TransformFilter{
					{Type: "CodeSystem"},
				},
			},
			input: &model.Type{},
			want:  true,
		}, {
			name: "type matches no filter",
			config: &config.Transform{
				Include: []*config.TransformFilter{
					{Name: "Some Name"},
				},
			},
			input: &model.Type{
				Name: "Not that name",
			},
			want: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			transform, err := transform.New(config.Mode("text"), tc.config)
			if err != nil {
				t.Fatalf("New() = %v", err)
			}

			got := transform.CanTransform(tc.input)

			if got, want := got, tc.want; got != want {
				t.Errorf("CanTransform() = %v, want %v", got, want)
			}
		})
	}
}

func TestTransformOutputPath(t *testing.T) {
	testCases := []struct {
		name    string
		config  *config.Transform
		input   any
		want    string
		wantErr error
	}{
		{
			name: "no templated values",
			config: &config.Transform{
				OutputPath: "testdata/output.tmpl",
			},
			input: &model.Type{},
			want:  "testdata/output.tmpl",
		}, {
			name: "with templated values",
			config: &config.Transform{
				OutputPath: "testdata/{{ .Name | lowercase }}.output",
			},
			input: &model.Type{
				Name: "SomeName",
			},
			want: "testdata/somename.output",
		}, {
			name: "bad template invocation",
			config: &config.Transform{
				OutputPath: "{{- .ThisFieldDoesntExist }}",
			},
			input: &model.Type{
				Name: "SomeName",
			},
			wantErr: cmpopts.AnyError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			transform, err := transform.New(config.Mode("text"), tc.config)
			if err != nil {
				t.Fatalf("New() = %v", err)
			}

			got, err := transform.OutputPath(tc.input)

			if got, want := err, tc.wantErr; !cmp.Equal(got, want, cmpopts.EquateErrors()) {
				t.Fatalf("OutputPath() = error %v, want %v", got, want)
			}
			if got, want := got, tc.want; got != want {
				t.Errorf("OutputPath() = %v, want %v", got, want)
			}
		})
	}
}
