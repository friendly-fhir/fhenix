package config_test

import (
	"testing"

	"github.com/friendly-fhir/fhenix/internal/config"
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

func TestConditions_UnmarshalYAML(t *testing.T) {
	strNode := nodeFrom(t, `"alias"`)
	alias := &yaml.Node{Kind: yaml.AliasNode, Alias: strNode}
	testCases := []struct {
		name    string
		input   *yaml.Node
		want    config.Conditions
		wantErr error
	}{
		{
			name:  "scalar",
			input: nodeFrom(t, `"scalar"`),
			want:  config.Conditions{"scalar"},
		}, {
			name:  "sequence",
			input: nodeFrom(t, `["sequence"]`),
			want:  config.Conditions{"sequence"},
		}, {
			name:  "alias",
			input: alias,
			want:  config.Conditions{"alias"},
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
			var got config.Conditions
			err := tc.input.Decode(&got)
			if got, want := err, tc.wantErr; !cmp.Equal(got, want, cmpopts.EquateErrors()) {
				t.Errorf("UnmarshalYAML(%s) = %v, want = %v", tc.input.Value, got, want)
			}

			if !cmp.Equal(got, tc.want, cmpopts.SortSlices(sortStrings)) {
				t.Errorf("UnmarshalYAML(%s) = %v, want = %v", tc.input.Value, got, tc.want)
			}
		})
	}
}

func sortStrings(lhs, rhs string) bool {
	return lhs < rhs
}
