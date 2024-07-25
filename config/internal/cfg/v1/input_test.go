package cfg_test

import (
	"strings"
	"testing"

	rootcfg "github.com/friendly-fhir/fhenix/config/internal/cfg"
	"github.com/friendly-fhir/fhenix/config/internal/cfg/v1"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"gopkg.in/yaml.v3"
)

func lines(lines ...string) string {
	return strings.Join(lines, "\n")
}

func zeroIfNil[T any](t *T) *T {
	if t != nil {
		return t
	}
	var result T
	return &result
}

func TestInput(t *testing.T) {
	testCases := []struct {
		name    string
		input   string
		want    *cfg.Input
		wantErr error
	}{
		{
			name: "valid input",
			input: lines(
				`packages:`,
				`  - name: hl7.fhir.r4.core`,
				`    version: "4.0.1"`,
			),
			want: &cfg.Input{
				Packages: []*cfg.InputPackage{
					{
						Name:    "hl7.fhir.r4.core",
						Version: "4.0.1",
					},
				},
			},
		}, {
			name:    "missing package",
			input:   lines("packages:"),
			wantErr: rootcfg.ErrMissingField,
		}, {
			name: "invalid package",
			input: lines(
				`packages:`,
				`  - name: hl7.fhir.r4.core`,
				`    version: "hello.world"`,
			),
			wantErr: rootcfg.ErrInvalidField,
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
			var got cfg.Input
			err := yaml.Unmarshal([]byte(tc.input), &got)

			if got, want := err, tc.wantErr; !cmp.Equal(got, want, cmpopts.EquateErrors()) {
				t.Fatalf("Input.UnmarshalYAML(...) = %v, want %v", got, want)
			}
			if want := zeroIfNil(tc.want); !cmp.Equal(&got, want) {
				t.Errorf("Input.UnmarshalYAML() = %v, want %v", got, want)
			}
		})
	}
}

func TestInputPackage(t *testing.T) {
	testCases := []struct {
		name    string
		input   string
		want    *cfg.InputPackage
		wantErr error
	}{
		{
			name: "valid input package",
			input: lines(
				`name: hl7.fhir.r4.core`,
				`version: "4.0.1"`,
			),
			want: &cfg.InputPackage{
				Name:    "hl7.fhir.r4.core",
				Version: "4.0.1",
			},
		}, {
			name: "missing name",
			input: lines(
				`version: "4.0.1"`,
			),
			wantErr: rootcfg.ErrMissingField,
		}, {
			name: "missing version",
			input: lines(
				`name: hl7.fhir.r4.core`,
			),
			wantErr: rootcfg.ErrMissingField,
		}, {
			name: "invalid name",
			input: lines(
				`name: "https://fhir.core/r4"`,
				`version: "4.0.1"`,
			),
			wantErr: rootcfg.ErrInvalidField,
		}, {
			name: "invalid version",
			input: lines(
				`name: hl7.fhir.r4.core`,
				`version: "hello.world"`,
			),
			wantErr: rootcfg.ErrInvalidField,
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
			var got cfg.InputPackage
			err := yaml.Unmarshal([]byte(tc.input), &got)

			if got, want := err, tc.wantErr; !cmp.Equal(got, want, cmpopts.EquateErrors()) {
				t.Fatalf("InputPackage.UnmarshalYAML(...) = %v, want %v", got, want)
			}
			if want := zeroIfNil(tc.want); !cmp.Equal(&got, want) {
				t.Errorf("InputPackage.UnmarshalYAML() = %v, want %v", got, want)
			}
		})
	}
}
