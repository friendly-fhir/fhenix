package config_test

import (
	"os"
	"testing"

	"github.com/friendly-fhir/fhenix/internal/config"
	"github.com/friendly-fhir/fhenix/internal/template"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"gopkg.in/yaml.v3"
)

func nodeFrom(t *testing.T, input string) *yaml.Node {
	t.Helper()

	var node yaml.Node
	if err := yaml.Unmarshal([]byte(input), &node); err != nil {
		t.Fatalf("yaml.Unmarshal: %v", err)
	}

	return &node
}

func TestVersion_UnmarshalYAML(t *testing.T) {
	testCases := []struct {
		name    string
		input   *yaml.Node
		want    config.Version
		wantErr error
	}{
		{
			name:  "valid",
			input: nodeFrom(t, `"1"`),
			want:  1,
		}, {
			name:    "invalid",
			input:   nodeFrom(t, `"2"`),
			wantErr: cmpopts.AnyError,
		}, {
			name:    "bad input",
			input:   nodeFrom(t, `sfgdsg`),
			wantErr: cmpopts.AnyError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var got config.Version
			err := tc.input.Decode(&got)
			if got, want := err, tc.wantErr; !cmp.Equal(got, want, cmpopts.EquateErrors()) {
				t.Errorf("UnmarshalYAML(%s) = %v, want = %v", tc.input.Value, got, want)
			}

			if got != tc.want {
				t.Errorf("UnmarshalYAML(%s) = %v, want = %v", tc.input.Value, got, tc.want)
			}
		})
	}
}

func TestType_UnmarshalYAML(t *testing.T) {
	testCases := []struct {
		name    string
		input   *yaml.Node
		want    config.Type
		wantErr error
	}{
		{
			name:  "structure-definition",
			input: nodeFrom(t, `"StructureDefinition"`),
			want:  config.TypeStructureDefinition,
		}, {
			name:  "value-set",
			input: nodeFrom(t, `"ValueSet"`),
			want:  config.TypeValueSet,
		}, {
			name:  "code-system",
			input: nodeFrom(t, `"CodeSystem"`),
			want:  config.TypeCodeSystem,
		}, {
			name:    "invalid",
			input:   nodeFrom(t, `"Invalid"`),
			wantErr: cmpopts.AnyError,
		}, {
			name:    "bad input",
			input:   nodeFrom(t, `sfgdsg`),
			wantErr: cmpopts.AnyError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var got config.Type
			err := tc.input.Decode(&got)
			if got, want := err, tc.wantErr; !cmp.Equal(got, want, cmpopts.EquateErrors()) {
				t.Errorf("UnmarshalYAML(%s) = %v, want = %v", tc.input.Value, got, want)
			}

			if got != tc.want {
				t.Errorf("UnmarshalYAML(%s) = %v, want = %v", tc.input.Value, got, tc.want)
			}
		})
	}
}

func TestCondition_UnmarshalYAML(t *testing.T) {
	strNode := nodeFrom(t, `"alias"`)
	alias := &yaml.Node{Kind: yaml.AliasNode, Alias: strNode}
	testCases := []struct {
		name    string
		input   *yaml.Node
		wantErr error
	}{
		{
			name:  "scalar",
			input: nodeFrom(t, `"scalar"`),
		}, {
			name:  "alias",
			input: alias,
		}, {
			name:    "invalid",
			input:   nodeFrom(t, `{}`),
			wantErr: cmpopts.AnyError,
		}, {
			name:    "invalid sequence",
			input:   nodeFrom(t, `[{}]`),
			wantErr: cmpopts.AnyError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var got config.Condition
			err := tc.input.Decode(&got)
			if got, want := err, tc.wantErr; !cmp.Equal(got, want, cmpopts.EquateErrors()) {
				t.Errorf("UnmarshalYAML(%s) = %v, want = %v", tc.input.Value, got, want)
			}
		})
	}
}

func TestCondition_Evaluate(t *testing.T) {
	testCases := []struct {
		name      string
		input     *config.Condition
		data      interface{}
		wantValue bool
	}{
		{
			name:      "empty",
			input:     config.NewCondition(nil),
			data:      nil,
			wantValue: true,
		}, {
			name: "true",
			input: config.NewCondition(
				template.MustParse("template", `{{ . }}`),
			),
			data:      true,
			wantValue: true,
		}, {
			name: "false",
			input: config.NewCondition(
				template.MustParse("template", `{{ . }}`),
			),
			data:      false,
			wantValue: false,
		}, {
			name: "42",
			input: config.NewCondition(
				template.MustParse("template", `{{ . }}`),
			),
			data:      42,
			wantValue: true,
		}, {
			name: "empty string",
			input: config.NewCondition(
				template.MustParse("template", `{{ . }}`),
			),
			data:      "",
			wantValue: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got, want := tc.input.Evaluate(tc.data), tc.wantValue; got != want {
				t.Errorf("Conditions.Evaluate(%v) = %v, want = %v", tc.data, got, want)
			}
		})
	}
}

func TestConfig(t *testing.T) {
	file, err := os.Open("testdata/config.yaml")
	if err != nil {
		t.Fatalf("os.Open: %v", err)
	}
	defer file.Close()

	var cfg config.Config
	if err := yaml.NewDecoder(file).Decode(&cfg); err != nil {
		t.Fatalf("yaml.NewDecoder.Decode: %v", err)
	}
}
