package model_test

import (
	"testing"

	"github.com/friendly-fhir/fhenix/internal/fhirig"
	"github.com/friendly-fhir/fhenix/internal/model"
	"github.com/google/go-cmp/cmp"
)

func TestDefaultTypeSet(t *testing.T) {
	want := "http://hl7.org/fhir/StructureDefinition"
	sut := model.DefaultTypeSet()

	got := sut.Base()

	if got != want {
		t.Errorf("DefaultTypeSet().Base() = %q, want %q", got, want)
	}
}

func TestTypeSet_Add(t *testing.T) {
	base := "http://example.com"
	ts := model.NewTypeSet(base)

	ty := &model.Type{Name: "Patient", URL: "http://example.com/Patient"}

	ts.Add(ty)

	t.Run("Increases Count", func(t *testing.T) {
		if got, want := len(ts.All()), 1; got != want {
			t.Errorf("TypeSet.Get(...) = %d entries, want %d entries", got, want)
		}
	})
	t.Run("Adds type to set", func(t *testing.T) {
		got, ok := ts.Lookup(ty.URL)
		if !ok {
			t.Fatalf("TypeSet.Get(...): lookup was nil, want non-nil")
		}

		if want := ty; got != want {
			t.Errorf("TypeSet.Get(...): lookup was %v, want %v", got, want)
		}
	})
}

func TestTypeSet_Lookup(t *testing.T) {
	ty := &model.Type{Name: "Patient", URL: "http://example.com/Patient"}
	testCases := []struct {
		name   string
		url    string
		want   *model.Type
		wantOK bool
	}{
		{
			name:   "Exact match",
			url:    "http://example.com/Patient",
			want:   ty,
			wantOK: true,
		}, {
			name:   "Relative match",
			url:    "Patient",
			want:   ty,
			wantOK: true,
		}, {
			name:   "No match",
			url:    "http://example.com/Encounter",
			want:   nil,
			wantOK: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ts := model.NewTypeSet("http://example.com", ty)
			got, ok := ts.Lookup(tc.url)

			if got, want := ok, tc.wantOK; got != want {
				t.Errorf("TypeSet.Lookup(%q): gotOK = %v, want %v", tc.url, got, want)
			}
			if got != tc.want {
				t.Errorf("TypeSet.Lookup(%q): got = %v, want %v", tc.url, got, tc.want)
			}
		})
	}
}

func TestTypeSet_Get(t *testing.T) {
	ty := &model.Type{Name: "Patient", URL: "http://example.com/Patient"}
	testCases := []struct {
		name string
		url  string
		want *model.Type
	}{
		{
			name: "Exact match",
			url:  "http://example.com/Patient",
			want: ty,
		}, {
			name: "Relative match",
			url:  "Patient",
			want: ty,
		}, {
			name: "No match",
			url:  "http://example.com/Encounter",
			want: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ts := model.NewTypeSet("http://example.com", ty)
			got := ts.Get(tc.url)

			if got != tc.want {
				t.Errorf("TypeSet.Get(%q): got = %v, want %v", tc.url, got, tc.want)
			}
		})
	}
}

func TestTypeSet_All(t *testing.T) {
	ty1 := &model.Type{Name: "Patient", URL: "http://example.com/Patient"}
	ty2 := &model.Type{Name: "Encounter", URL: "http://example.com/Encounter"}
	want := []*model.Type{ty2, ty1}
	ts := model.NewTypeSet("http://example.com", ty1, ty2)

	got := ts.All()

	if len(got) != len(want) {
		t.Fatalf("TypeSet.All(): got %d entries, want %d entries", len(got), len(want))
	}

	if !cmp.Equal(got, want) {
		t.Errorf("TypeSet.All() = %v, want %v", got, want)
	}
}

func TestTypeSet_Resources(t *testing.T) {
	ty1 := &model.Type{Name: "Patient", URL: "http://example.com/Patient", Kind: "resource"}
	ty2 := &model.Type{Name: "Encounter", URL: "http://example.com/Encounter", Kind: "complex-type"}
	want := []*model.Type{ty1}
	ts := model.NewTypeSet("http://example.com", ty1, ty2)

	got := ts.Resources().All()

	if len(got) != len(want) {
		t.Fatalf("TypeSet.Resources(): got %d entries, want %d entries", len(got), len(want))
	}

	if !cmp.Equal(got, want) {
		t.Errorf("TypeSet.Resources() = %v, want %v", got, want)
	}
}

func TestTypeSet_PrimitiveTypes(t *testing.T) {
	ty1 := &model.Type{Name: "Patient", URL: "http://example.com/Patient", Kind: "primitive-type"}
	ty2 := &model.Type{Name: "Encounter", URL: "http://example.com/Encounter", Kind: "complex-type"}
	want := []*model.Type{ty1}
	ts := model.NewTypeSet("http://example.com", ty1, ty2)

	got := ts.PrimitiveTypes().All()

	if len(got) != len(want) {
		t.Fatalf("TypeSet.PrimitiveTypes(): got %d entries, want %d entries", len(got), len(want))
	}

	if !cmp.Equal(got, want) {
		t.Errorf("TypeSet.PrimitiveTypes() = %v, want %v", got, want)
	}
}

func TestTypeSet_ComplexTypes(t *testing.T) {
	ty1 := &model.Type{Name: "Patient", URL: "http://example.com/Patient", Kind: "complex-type"}
	ty2 := &model.Type{Name: "Encounter", URL: "http://example.com/Encounter", Kind: "resource"}
	want := []*model.Type{ty1}
	ts := model.NewTypeSet("http://example.com", ty1, ty2)

	got := ts.ComplexTypes().All()

	if len(got) != len(want) {
		t.Fatalf("TypeSet.ComplexTypes(): got %d entries, want %d entries", len(got), len(want))
	}

	if !cmp.Equal(got, want) {
		t.Errorf("TypeSet.ComplexTypes() = %v, want %v", got, want)
	}
}

func TestTypeSet_InBase(t *testing.T) {
	ty1 := &model.Type{Name: "Patient", URL: "http://example.com/Patient"}
	ty2 := &model.Type{Name: "Encounter", URL: "http://example.com/Encounter"}
	ty3 := &model.Type{Name: "Observation", URL: "http://second-example.com/Observation"}
	want := []*model.Type{ty2, ty1}
	ts := model.NewTypeSet("http://example.com", ty1, ty2, ty3)

	got := ts.InBase().All()

	if len(got) != len(want) {
		t.Fatalf("TypeSet.InBase(): got %d entries, want %d entries", len(got), len(want))
	}

	if !cmp.Equal(got, want) {
		t.Errorf("TypeSet.InBase() = %v, want %v", got, want)
	}
}

func TestTypeSet_DefinedInPackage(t *testing.T) {
	pkg1 := fhirig.NewPackage("hl7.fhir.r4.core", "4.0.1")
	pkg2 := fhirig.NewPackage("hl7.fhir.us.core", "6.0.0")
	ty1 := &model.Type{Name: "Patient", URL: "http://example.com/Patient", Source: &model.TypeSource{Package: pkg1}}
	ty2 := &model.Type{Name: "Encounter", URL: "http://example.com/Encounter", Source: &model.TypeSource{Package: pkg1}}
	ty3 := &model.Type{Name: "Observation", URL: "http://second-example.com/Observation", Source: &model.TypeSource{Package: pkg2}}
	want := []*model.Type{ty2, ty1}
	ts := model.NewTypeSet("http://example.com", ty1, ty2, ty3)

	got := ts.DefinedInPackage(pkg1).All()

	if len(got) != len(want) {
		t.Fatalf("TypeSet.DefinedInPackage(...): got %d entries, want %d entries", len(got), len(want))
	}

	if !cmp.Equal(got, want) {
		t.Errorf("TypeSet.DefinedInPackage(...): got = %v, want %v", got, want)
	}
}
