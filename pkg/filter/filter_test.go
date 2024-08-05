package filter_test

import (
	"testing"

	"github.com/friendly-fhir/fhenix/pkg/config"
	"github.com/friendly-fhir/fhenix/pkg/filter"
	"github.com/friendly-fhir/fhenix/pkg/model"
	"github.com/friendly-fhir/fhenix/pkg/registry"
)

func TestFilterMatchesType(t *testing.T) {
	testCases := []struct {
		name string
		cfg  *config.TransformFilter
		t    *model.Type
		want bool
	}{
		{
			name: "Empty filter matches nothing",
			cfg:  &config.TransformFilter{},
			t:    &model.Type{},
			want: false,
		}, {
			name: "Filter matches type by name",
			cfg:  &config.TransformFilter{Name: "Patient"},
			t:    &model.Type{Name: "Patient"},
			want: true,
		}, {
			name: "Filter does not match type by name",
			cfg:  &config.TransformFilter{Name: "Patient"},
			t:    &model.Type{Name: "Observation"},
			want: false,
		}, {
			name: "Filter matches type by pattern",
			cfg:  &config.TransformFilter{Name: "P.*"},
			t:    &model.Type{Name: "Patient"},
			want: true,
		}, {
			name: "Filter does not match type by pattern",
			cfg:  &config.TransformFilter{Name: "P.*"},
			t:    &model.Type{Name: "Observation"},
			want: false,
		}, {
			name: "Filter matches type by source",
			cfg:  &config.TransformFilter{Source: "file.go"},
			t:    &model.Type{Source: &model.TypeSource{File: "file.go"}},
			want: true,
		}, {
			name: "Filter does not match type by source",
			cfg:  &config.TransformFilter{Source: "file.go"},
			t:    &model.Type{Source: &model.TypeSource{File: "file_test.go"}},
			want: false,
		}, {
			name: "Filter matches type by source pattern",
			cfg:  &config.TransformFilter{Source: ".*\\.go"},
			t:    &model.Type{Source: &model.TypeSource{File: "file.go"}},
			want: true,
		}, {
			name: "Filter does not match type by source pattern",
			cfg:  &config.TransformFilter{Source: ".*\\.go"},
			t:    &model.Type{Source: &model.TypeSource{File: "file_test.py"}},
			want: false,
		}, {
			name: "Nil filter",
			cfg:  nil,
			t:    &model.Type{},
			want: false,
		}, {
			name: "Filter does not match type",
			cfg:  &config.TransformFilter{Type: "CodeSystem"},
			t:    &model.Type{Name: "Patient"},
			want: false,
		}, {
			name: "Filter matches by URL",
			cfg:  &config.TransformFilter{URL: "http://example.com"},
			t:    &model.Type{URL: "http://example.com"},
			want: true,
		}, {
			name: "Filter does not match by URL",
			cfg:  &config.TransformFilter{URL: "http://example.com"},
			t:    &model.Type{URL: "http://example.org"},
			want: false,
		}, {
			name: "Filter matches by package",
			cfg:  &config.TransformFilter{Package: "hl7.fhir.core.r4"},
			t: &model.Type{
				Source: &model.TypeSource{
					Package: registry.NewPackageRef("default", "hl7.fhir.core.r4", "4.0.1"),
				},
			},
			want: true,
		}, {
			name: "Filter does not match by package",
			cfg:  &config.TransformFilter{Package: "hl7.fhir.core.r4"},
			t: &model.Type{
				Source: &model.TypeSource{
					Package: registry.NewPackageRef("default", "hl7.fhir.core.r5", "5.0.0"),
				},
			},
			want: false,
		}, {
			name: "Filter matches by condition",
			cfg:  &config.TransformFilter{Condition: `{{ eq .Kind "primitive-type" }}`},
			t:    &model.Type{Kind: "primitive-type"},
			want: true,
		}, {
			name: "Filter does not match by condition",
			cfg:  &config.TransformFilter{Condition: `{{ eq .Kind "primitive-type" }}`},
			t:    &model.Type{Kind: "complex-type"},
			want: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			filter := filter.New(tc.cfg)

			got := filter.MatchesType(tc.t)

			if got != tc.want {
				t.Errorf("Filter.MatchesType(%s) = %v, want = %v", tc.t.Name, got, tc.want)
			}
		})
	}
}
