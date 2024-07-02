package cfg_test

import (
	"testing"

	rootcfg "github.com/friendly-fhir/fhenix/internal/cfg"
	"github.com/friendly-fhir/fhenix/internal/cfg/v1"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"gopkg.in/yaml.v3"
)

func TestTransform(t *testing.T) {
	testCases := []struct {
		name    string
		input   string
		want    *cfg.Transform
		wantErr error
	}{
		{
			name: "valid transform",
			input: lines(
				`templates:`,
				`  main: "main"`,
				`  header: "header"`,
				`  footer: "footer"`,
				`  extra: "extra"`,
				`include:`,
				`  - name: "filter"`,
				`    condition: '{{ .Field | eq "value" }}'`,
			),
			want: &cfg.Transform{
				Templates: &cfg.TransformTemplates{
					Main:   "main",
					Header: "header",
					Footer: "footer",
					Partials: map[string]string{
						"extra": "extra",
					},
				},
				Include: []*cfg.TransformFilter{
					{
						Name:      "filter",
						Condition: `{{ .Field | eq "value" }}`,
					},
				},
			},
		}, {
			name: "Bad function name specified",
			input: lines(
				`funcs:`,
				`  "invalid name": "func.tmpl"`,
			),
			wantErr: rootcfg.ErrInvalidField,
		}, {
			name: "Bad output-path template",
			input: lines(
				`output-path: "{{ .Field | invalid-func }}"`,
			),
			wantErr: cmpopts.AnyError,
		}, {
			name: "invalid type",
			input: lines(
				`"hello"`,
			),
			wantErr: cmpopts.AnyError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var got cfg.Transform
			err := yaml.Unmarshal([]byte(tc.input), &got)

			if got, want := err, tc.wantErr; !cmp.Equal(got, want, cmpopts.EquateErrors()) {
				t.Fatalf("Transform.UnmarshalYAML(...) = %v, want %v", got, want)
			}
			if want := zeroIfNil(tc.want); !cmp.Equal(&got, want) {
				t.Errorf("Transform.UnmarshalYAML() = %v, want %v", got, want)
			}
		})
	}
}

func TestTransformTemplates(t *testing.T) {
	testCases := []struct {
		name    string
		input   string
		want    *cfg.TransformTemplates
		wantErr error
	}{
		{
			name: "main is set",
			input: lines(
				`main: "main"`,
			),
			want: &cfg.TransformTemplates{
				Main: "main",
			},
		}, {
			name: "header is set",
			input: lines(
				`header: "header"`,
			),
			want: &cfg.TransformTemplates{
				Header: "header",
			},
		}, {
			name: "footer is set",
			input: lines(
				`footer: "footer"`,
			),
			want: &cfg.TransformTemplates{
				Footer: "footer",
			},
		}, {
			name: "multiple set, extras become partials",
			input: lines(
				`main: "main"`,
				`header: "header"`,
				`footer: "footer"`,
				`extra: "extra"`,
			),
			want: &cfg.TransformTemplates{
				Main:   "main",
				Header: "header",
				Footer: "footer",
				Partials: map[string]string{
					"extra": "extra",
				},
			},
		}, {
			name: "invalid type",
			input: lines(
				`"hello"`,
			),
			wantErr: cmpopts.AnyError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var got cfg.TransformTemplates
			err := yaml.Unmarshal([]byte(tc.input), &got)

			if got, want := err, tc.wantErr; !cmp.Equal(got, want, cmpopts.EquateErrors()) {
				t.Fatalf("TransformTemplates.UnmarshalYAML(...) = %v, want %v", got, want)
			}
			if want := zeroIfNil(tc.want); !cmp.Equal(&got, want, cmpopts.EquateEmpty()) {
				t.Errorf("TransformTemplates.UnmarshalYAML() = %v, want %v", got, want)
			}
		})
	}
}

func TestTransformFilter(t *testing.T) {
	testCases := []struct {
		name    string
		input   string
		want    *cfg.TransformFilter
		wantErr error
	}{
		{
			name: "valid filter",
			input: lines(
				`name: "filter"`,
			),
			want: &cfg.TransformFilter{
				Name: "filter",
			},
		}, {
			name: "valid filter with condition",
			input: lines(
				`name: "filter"`,
				`condition: "{{ .Field | eq \"value\" }}"`,
			),
			want: &cfg.TransformFilter{
				Name:      "filter",
				Condition: `{{ .Field | eq "value" }}`,
			},
		}, {
			name: "valid regex name",
			input: lines(
				`name: "[a-z]+"`,
			),
			want: &cfg.TransformFilter{
				Name: "[a-z]+",
			},
		}, {
			name: "empty",
			// using a non-existent field so that the check for missing fields is triggered
			input:   "some-imaginary-field: value",
			wantErr: rootcfg.ErrMissingField,
		}, {
			name: "name is invalid regex",
			input: lines(
				`name: "["`,
			),
			wantErr: cmpopts.AnyError,
		}, {
			name: "condition is invalid template",
			input: lines(
				`condition: "{{ .Field | invalid-func }}"`,
			),
			wantErr: cmpopts.AnyError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var got cfg.TransformFilter
			err := yaml.Unmarshal([]byte(tc.input), &got)

			if got, want := err, tc.wantErr; !cmp.Equal(got, want, cmpopts.EquateErrors()) {
				t.Fatalf("TransformFilter.UnmarshalYAML(...) = %v, want %v", got, want)
			}
			if want := zeroIfNil(tc.want); !cmp.Equal(&got, want) {
				t.Errorf("TransformFilter.UnmarshalYAML() = %v, want %v", got, want)
			}
		})
	}
}
