package model_test

import (
	"testing"

	"github.com/friendly-fhir/fhenix/internal/model"
	"github.com/friendly-fhir/fhenix/internal/model/raw"
	"github.com/google/go-cmp/cmp"
)

func TestFHIRCache_StructureDefinitions(t *testing.T) {
	sut := model.NewFHIRCache("http://example.com")
	sd1 := &raw.StructureDefinition{URL: "sd1"}
	sd2 := &raw.StructureDefinition{URL: "http://example.com/StructureDefinition/sd2"}
	want := []*raw.StructureDefinition{sd1, sd2}
	sut.AddStructureDefinition(nil, "", sd1)
	sut.AddStructureDefinition(nil, "", sd2)

	entries := sut.StructureDefinitions()
	var got []*raw.StructureDefinition
	for _, e := range entries {
		got = append(got, e.Definition)
	}

	if len(got) != 2 {
		t.Fatalf("FHIRCache.StructureDefinitions() = %d entries, want 2 entries", len(got))
	}

	if !cmp.Equal(got, want) {
		t.Errorf("FHIRCache.StructureDefinitions() = %v, want %v", got, want)
	}
}

func TestFHIRCache_GetStructureDefinition(t *testing.T) {
	sd := &raw.StructureDefinition{URL: "http://example.com/StructureDefinition/sd"}
	sut := model.NewFHIRCache("http://example.com")
	sut.AddStructureDefinition(nil, "", sd)

	testCases := []struct {
		name string
		url  string
		want *raw.StructureDefinition
	}{
		{
			name: "Exact match",
			url:  "http://example.com/StructureDefinition/sd",
			want: sd,
		}, {
			name: "Root match",
			url:  "sd",
			want: sd,
		}, {
			name: "Not found",
			url:  "http://example.com/StructureDefinition/missing",
			want: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := sut.GetStructureDefinition(tc.url)

			if got != nil && got.Definition != tc.want {
				t.Errorf("FHIRCache.GetStructureDefinition(%q) = %v, want %v", tc.url, got, tc.want)
			}
		})
	}
}

func TestFHIRCache_LookupStructureDefinition(t *testing.T) {
	sd := &raw.StructureDefinition{URL: "http://example.com/StructureDefinition/sd"}
	sut := model.NewFHIRCache("http://example.com")
	sut.AddStructureDefinition(nil, "", sd)

	testCases := []struct {
		name   string
		url    string
		want   *raw.StructureDefinition
		wantOK bool
	}{
		{
			name:   "Exact match",
			url:    "http://example.com/StructureDefinition/sd",
			want:   sd,
			wantOK: true,
		}, {
			name:   "Root match",
			url:    "sd",
			want:   sd,
			wantOK: true,
		}, {
			name:   "Not found",
			url:    "http://example.com/StructureDefinition/missing",
			want:   nil,
			wantOK: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := sut.LookupStructureDefinition(tc.url)

			if ok != tc.wantOK {
				t.Fatalf("FHIRCache.LookupStructureDefinition(%q) = %t, want %t", tc.url, ok, tc.wantOK)
			}
			if got != nil && got.Definition != tc.want {
				t.Errorf("FHIRCache.LookupStructureDefinition(%q) = %v, want %v", tc.url, got, tc.want)
			}
		})
	}
}
