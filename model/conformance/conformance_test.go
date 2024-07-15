package conformance_test

import (
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/friendly-fhir/fhenix/model/conformance"
	"github.com/friendly-fhir/fhenix/model/conformance/definition"
	"github.com/friendly-fhir/fhenix/registry"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func mustParseJSON[T any](t *testing.T, bytes []byte) *T {
	t.Helper()

	var out T
	if err := json.Unmarshal(bytes, &out); err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}
	return &out
}

func mustReadJSON[T any](t *testing.T, path string) *T {
	t.Helper()

	bytes, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}
	return mustParseJSON[T](t, bytes)
}

func TestModuleParseFile(t *testing.T) {
	testCases := []struct {
		name    string
		path    string
		url     string
		wantErr error
	}{
		{
			name: "structure definition",
			path: "testdata/structure-definition.json",
			url:  mustReadJSON[definition.StructureDefinition](t, "testdata/structure-definition.json").GetURL().GetValue(),
		}, {
			name: "code system",
			path: "testdata/code-system.json",
			url:  mustReadJSON[definition.CodeSystem](t, "testdata/code-system.json").GetURL().GetValue(),
		}, {
			name: "value set",
			path: "testdata/value-set.json",
			url:  mustReadJSON[definition.ValueSets](t, "testdata/value-set.json").GetURL().GetValue(),
		}, {
			name: "concept map",
			path: "testdata/concept-map.json",
			url:  mustReadJSON[definition.ConceptMap](t, "testdata/concept-map.json").GetURL().GetValue(),
		}, {
			name:    "invalid path",
			path:    "testdata/invalid.json",
			wantErr: fs.ErrNotExist,
		}, {
			name:    "not a resource",
			path:    "testdata/not-resource.json",
			wantErr: cmpopts.AnyError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			module := conformance.DefaultModule()
			pkg := registry.NewPackageRef("default", "hl7.fhir.r4.core", "4.0.1")

			err := module.ParseFile(tc.path, pkg)
			if got, want := err, tc.wantErr; !cmp.Equal(got, want, cmpopts.EquateErrors()) {
				t.Fatalf("Module.ParseFile(%q, %q) = error %v, want %v", tc.path, tc.url, got, want)
			}
			if got, want := module.Contains(tc.url), (err == nil); got != want {
				t.Fatalf("Module.Contains(%q) = %v, want %v", tc.url, got, want)
			}
		})
	}
}

func TestModuleLookupCanonical(t *testing.T) {
	sd := mustReadJSON[definition.StructureDefinition](t, "testdata/structure-definition.json")
	cs := mustReadJSON[definition.CodeSystem](t, "testdata/code-system.json")
	vs := mustReadJSON[definition.ValueSets](t, "testdata/value-set.json")
	cm := mustReadJSON[definition.ConceptMap](t, "testdata/concept-map.json")

	module := conformance.DefaultModule()
	pkg := registry.NewPackageRef("default", "hl7.fhir.r4.core", "4.0.1")
	module.AddDefinition(sd, &conformance.Source{Package: pkg})
	module.AddDefinition(cs, &conformance.Source{Package: pkg})
	module.AddDefinition(vs, &conformance.Source{Package: pkg})
	module.AddDefinition(cm, &conformance.Source{Package: pkg})

	testCases := []struct {
		name   string
		url    string
		want   definition.Canonical
		wantOK bool
	}{
		{
			name:   "structure definition",
			url:    sd.GetURL().GetValue(),
			want:   sd,
			wantOK: true,
		}, {
			name:   "structure definition base relative",
			url:    filepath.Base(sd.GetURL().GetValue()),
			want:   sd,
			wantOK: true,
		}, {
			name:   "code system",
			url:    cs.GetURL().GetValue(),
			want:   cs,
			wantOK: true,
		}, {
			name:   "code system base relative",
			url:    filepath.Base(cs.GetURL().GetValue()),
			want:   cs,
			wantOK: true,
		}, {
			name:   "value set",
			url:    vs.GetURL().GetValue(),
			want:   vs,
			wantOK: true,
		}, {
			name:   "value set base relative",
			url:    filepath.Base(vs.GetURL().GetValue()),
			want:   vs,
			wantOK: true,
		}, {
			name:   "concept map",
			url:    cm.GetURL().GetValue(),
			want:   cm,
			wantOK: true,
		}, {
			name:   "concept map base relative",
			url:    filepath.Base(cm.GetURL().GetValue()),
			want:   cm,
			wantOK: true,
		}, {
			name:   "unknown",
			url:    "http://example.com/unknown",
			want:   nil,
			wantOK: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			got, ok := module.LookupCanonical(tc.url)
			if got, want := ok, tc.wantOK; got != want {
				t.Fatalf("Module.LookupCanonical(%q) ok = %v, want %v", tc.want.GetURL().GetValue(), got, want)
			}
			if got, want := got, tc.want; !cmp.Equal(got, want) {
				t.Errorf("Module.LookupCanonical(%q) = %v, want %v", tc.want.GetURL().GetValue(), got, want)
			}
		})
	}
}

func TestModuleLookupSource(t *testing.T) {
	sd := mustReadJSON[definition.StructureDefinition](t, "testdata/structure-definition.json")
	cs := mustReadJSON[definition.CodeSystem](t, "testdata/code-system.json")
	vs := mustReadJSON[definition.ValueSets](t, "testdata/value-set.json")
	cm := mustReadJSON[definition.ConceptMap](t, "testdata/concept-map.json")
	pkg := registry.NewPackageRef("default", "hl7.fhir.r4.core", "4.0.1")
	sdSource := &conformance.Source{Package: pkg}
	csSource := &conformance.Source{Package: pkg}
	vsSource := &conformance.Source{Package: pkg}
	cmSource := &conformance.Source{Package: pkg}
	module := conformance.DefaultModule()
	module.AddDefinition(sd, sdSource)
	module.AddDefinition(cs, csSource)
	module.AddDefinition(vs, vsSource)
	module.AddDefinition(cm, cmSource)

	testCases := []struct {
		name   string
		url    string
		want   *conformance.Source
		wantOK bool
	}{
		{
			name:   "structure definition",
			url:    sd.GetURL().GetValue(),
			want:   sdSource,
			wantOK: true,
		}, {
			name:   "structure definition base relative",
			url:    filepath.Base(sd.GetURL().GetValue()),
			want:   sdSource,
			wantOK: true,
		}, {
			name:   "code system",
			url:    cs.GetURL().GetValue(),
			want:   csSource,
			wantOK: true,
		}, {
			name:   "code system base relative",
			url:    filepath.Base(cs.GetURL().GetValue()),
			want:   csSource,
			wantOK: true,
		}, {
			name:   "value set",
			url:    vs.GetURL().GetValue(),
			want:   vsSource,
			wantOK: true,
		}, {
			name:   "value set base relative",
			url:    filepath.Base(vs.GetURL().GetValue()),
			want:   vsSource,
			wantOK: true,
		}, {
			name:   "concept map",
			url:    cm.GetURL().GetValue(),
			want:   cmSource,
			wantOK: true,
		}, {
			name:   "concept map base relative",
			url:    filepath.Base(cm.GetURL().GetValue()),
			want:   cmSource,
			wantOK: true,
		}, {
			name:   "unknown",
			url:    "http://example.com/unknown",
			want:   nil,
			wantOK: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := module.LookupSource(tc.url)

			if got, want := ok, tc.wantOK; got != want {
				t.Fatalf("Module.LookupSource(%q) ok = %v, want %v", tc.url, got, want)
			}
			if got, want := got, tc.want; !cmp.Equal(got, want) {
				t.Errorf("Module.LookupSource(%q) = %v, want %v", tc.url, got, want)
			}
		})
	}
}

func TestModuleCanonical(t *testing.T) {
	sd := mustReadJSON[definition.StructureDefinition](t, "testdata/structure-definition.json")
	cs := mustReadJSON[definition.CodeSystem](t, "testdata/code-system.json")
	vs := mustReadJSON[definition.ValueSets](t, "testdata/value-set.json")
	cm := mustReadJSON[definition.ConceptMap](t, "testdata/concept-map.json")

	module := conformance.DefaultModule()
	pkg := registry.NewPackageRef("default", "hl7.fhir.r4.core", "4.0.1")
	module.AddDefinition(sd, &conformance.Source{Package: pkg})
	module.AddDefinition(cs, &conformance.Source{Package: pkg})
	module.AddDefinition(vs, &conformance.Source{Package: pkg})
	module.AddDefinition(cm, &conformance.Source{Package: pkg})

	testCases := []struct {
		name string
		url  string
		want definition.Canonical
	}{
		{
			name: "structure definition",
			url:  sd.GetURL().GetValue(),
			want: sd,
		}, {
			name: "code system",
			url:  cs.GetURL().GetValue(),
			want: cs,
		}, {
			name: "value set",
			url:  vs.GetURL().GetValue(),
			want: vs,
		}, {
			name: "concept map",
			url:  cm.GetURL().GetValue(),
			want: cm,
		}, {
			name: "unknown",
			url:  "http://example.com/unknown",
			want: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := module.Canonical(tc.url)

			if got, want := got, tc.want; !cmp.Equal(got, want) {
				t.Errorf("Module.Canonical(%q) = %v, want %v", tc.url, got, want)
			}
		})
	}
}

func TestModuleSource(t *testing.T) {
	sd := mustReadJSON[definition.StructureDefinition](t, "testdata/structure-definition.json")
	cs := mustReadJSON[definition.CodeSystem](t, "testdata/code-system.json")
	vs := mustReadJSON[definition.ValueSets](t, "testdata/value-set.json")
	cm := mustReadJSON[definition.ConceptMap](t, "testdata/concept-map.json")
	pkg := registry.NewPackageRef("default", "hl7.fhir.r4.core", "4.0.1")
	sdSource := &conformance.Source{Package: pkg}
	csSource := &conformance.Source{Package: pkg}
	vsSource := &conformance.Source{Package: pkg}
	cmSource := &conformance.Source{Package: pkg}
	module := conformance.DefaultModule()
	module.AddDefinition(sd, sdSource)
	module.AddDefinition(cs, csSource)
	module.AddDefinition(vs, vsSource)
	module.AddDefinition(cm, cmSource)

	testCases := []struct {
		name string
		url  string
		want *conformance.Source
	}{
		{
			name: "structure definition",
			url:  sd.GetURL().GetValue(),
			want: sdSource,
		}, {
			name: "code system",
			url:  cs.GetURL().GetValue(),
			want: csSource,
		}, {
			name: "value set",
			url:  vs.GetURL().GetValue(),
			want: vsSource,
		}, {
			name: "concept map",
			url:  cm.GetURL().GetValue(),
			want: cmSource,
		}, {
			name: "unknown",
			url:  "http://example.com/unknown",
			want: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := module.Source(tc.url)

			if got, want := got, tc.want; !cmp.Equal(got, want) {
				t.Errorf("Module.Source(%q) = %v, want %v", tc.url, got, want)
			}
		})
	}
}

func TestModuleSourceOf(t *testing.T) {
	sd := mustReadJSON[definition.StructureDefinition](t, "testdata/structure-definition.json")
	cs := mustReadJSON[definition.CodeSystem](t, "testdata/code-system.json")
	vs := mustReadJSON[definition.ValueSets](t, "testdata/value-set.json")
	cm := mustReadJSON[definition.ConceptMap](t, "testdata/concept-map.json")
	pkg := registry.NewPackageRef("default", "hl7.fhir.r4.core", "4.0.1")
	sdSource := &conformance.Source{Package: pkg}
	csSource := &conformance.Source{Package: pkg}
	vsSource := &conformance.Source{Package: pkg}
	cmSource := &conformance.Source{Package: pkg}
	module := conformance.DefaultModule()
	module.AddDefinition(sd, sdSource)
	module.AddDefinition(cs, csSource)
	module.AddDefinition(vs, vsSource)
	module.AddDefinition(cm, cmSource)

	testCases := []struct {
		name      string
		canonical definition.Canonical
		want      *conformance.Source
	}{
		{
			name:      "structure definition",
			canonical: sd,
			want:      sdSource,
		}, {
			name:      "code system",
			canonical: cs,
			want:      csSource,
		}, {
			name:      "value set",
			canonical: vs,
			want:      vsSource,
		}, {
			name:      "concept map",
			canonical: cm,
			want:      cmSource,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := module.SourceOf(tc.canonical)

			if got, want := got, tc.want; !cmp.Equal(got, want) {
				t.Errorf("Module.SourceOf(%v) = %v, want %v", tc.canonical, got, want)
			}
		})
	}
}

func TestModuleStructureDefinitions(t *testing.T) {
	sd := mustReadJSON[definition.StructureDefinition](t, "testdata/structure-definition.json")
	cs := mustReadJSON[definition.CodeSystem](t, "testdata/code-system.json")
	vs := mustReadJSON[definition.ValueSets](t, "testdata/value-set.json")
	pkg := registry.NewPackageRef("default", "hl7.fhir.r4.core", "4.0.1")
	sdSource := &conformance.Source{Package: pkg}
	csSource := &conformance.Source{Package: pkg}
	vsSource := &conformance.Source{Package: pkg}
	module := conformance.DefaultModule()
	module.AddDefinition(sd, sdSource)
	module.AddDefinition(cs, csSource)
	module.AddDefinition(vs, vsSource)

	want := []*definition.StructureDefinition{sd}
	got := module.StructureDefinitions()

	if !cmp.Equal(got, want) {
		t.Errorf("Module.StructureDefinitions() = %v, want %v", got, want)
	}
}

func TestModuleValueSets(t *testing.T) {
	sd := mustReadJSON[definition.StructureDefinition](t, "testdata/structure-definition.json")
	cs := mustReadJSON[definition.CodeSystem](t, "testdata/code-system.json")
	vs := mustReadJSON[definition.ValueSets](t, "testdata/value-set.json")
	pkg := registry.NewPackageRef("default", "hl7.fhir.r4.core", "4.0.1")
	sdSource := &conformance.Source{Package: pkg}
	csSource := &conformance.Source{Package: pkg}
	vsSource := &conformance.Source{Package: pkg}
	module := conformance.DefaultModule()
	module.AddDefinition(sd, sdSource)
	module.AddDefinition(cs, csSource)
	module.AddDefinition(vs, vsSource)

	want := []*definition.ValueSets{vs}
	got := module.ValueSets()

	if !cmp.Equal(got, want) {
		t.Errorf("Module.ValueSets() = %v, want %v", got, want)
	}
}

func TestModuleCodeSystems(t *testing.T) {
	sd := mustReadJSON[definition.StructureDefinition](t, "testdata/structure-definition.json")
	cs := mustReadJSON[definition.CodeSystem](t, "testdata/code-system.json")
	vs := mustReadJSON[definition.ValueSets](t, "testdata/value-set.json")
	pkg := registry.NewPackageRef("default", "hl7.fhir.r4.core", "4.0.1")
	sdSource := &conformance.Source{Package: pkg}
	csSource := &conformance.Source{Package: pkg}
	vsSource := &conformance.Source{Package: pkg}
	module := conformance.DefaultModule()
	module.AddDefinition(sd, sdSource)
	module.AddDefinition(cs, csSource)
	module.AddDefinition(vs, vsSource)

	want := []*definition.CodeSystem{cs}
	got := module.CodeSystems()

	if !cmp.Equal(got, want) {
		t.Errorf("Module.CodeSystems() = %v, want %v", got, want)
	}
}

func TestModuleConceptMaps(t *testing.T) {
	sd := mustReadJSON[definition.StructureDefinition](t, "testdata/structure-definition.json")
	cs := mustReadJSON[definition.CodeSystem](t, "testdata/code-system.json")
	vs := mustReadJSON[definition.ValueSets](t, "testdata/value-set.json")
	cm := mustReadJSON[definition.ConceptMap](t, "testdata/concept-map.json")
	pkg := registry.NewPackageRef("default", "hl7.fhir.r4.core", "4.0.1")
	sdSource := &conformance.Source{Package: pkg}
	csSource := &conformance.Source{Package: pkg}
	vsSource := &conformance.Source{Package: pkg}
	cmSource := &conformance.Source{Package: pkg}
	module := conformance.DefaultModule()
	module.AddDefinition(sd, sdSource)
	module.AddDefinition(cs, csSource)
	module.AddDefinition(vs, vsSource)
	module.AddDefinition(cm, cmSource)

	want := []*definition.ConceptMap{cm}
	got := module.ConceptMaps()

	if !cmp.Equal(got, want) {
		t.Errorf("Module.ConceptMaps() = %v, want %v", got, want)
	}
}

func TestModuleAll(t *testing.T) {
	sd := mustReadJSON[definition.StructureDefinition](t, "testdata/structure-definition.json")
	cs := mustReadJSON[definition.CodeSystem](t, "testdata/code-system.json")
	vs := mustReadJSON[definition.ValueSets](t, "testdata/value-set.json")
	cm := mustReadJSON[definition.ConceptMap](t, "testdata/concept-map.json")
	pkg := registry.NewPackageRef("default", "hl7.fhir.r4.core", "4.0.1")
	sdSource := &conformance.Source{Package: pkg}
	csSource := &conformance.Source{Package: pkg}
	vsSource := &conformance.Source{Package: pkg}
	cmSource := &conformance.Source{Package: pkg}
	module := conformance.DefaultModule()
	module.AddDefinition(sd, sdSource)
	module.AddDefinition(cs, csSource)
	module.AddDefinition(vs, vsSource)
	module.AddDefinition(cm, cmSource)

	want := []definition.Canonical{cs, sd, vs, cm}
	got := module.All()

	if !cmp.Equal(got, want, cmpopts.SortSlices(sortCanonical)) {
		t.Errorf("Module.All() = %v, want %v", got, want)
	}
}

func sortCanonical(lhs, rhs definition.Canonical) bool {
	return lhs.GetURL().GetValue() < rhs.GetURL().GetValue()
}

func TestModuleFilterStructureDefinitions(t *testing.T) {
	sd := mustReadJSON[definition.StructureDefinition](t, "testdata/structure-definition.json")
	cs := mustReadJSON[definition.CodeSystem](t, "testdata/code-system.json")
	vs := mustReadJSON[definition.ValueSets](t, "testdata/value-set.json")
	pkg := registry.NewPackageRef("default", "hl7.fhir.r4.core", "4.0.1")
	sdSource := &conformance.Source{Package: pkg}
	csSource := &conformance.Source{Package: pkg}
	vsSource := &conformance.Source{Package: pkg}
	module := conformance.DefaultModule()
	module.AddDefinition(sd, sdSource)
	module.AddDefinition(cs, csSource)
	module.AddDefinition(vs, vsSource)

	testCases := []struct {
		name string
		pkg  registry.PackageRef
		want []*definition.StructureDefinition
	}{
		{
			name: "structure definition is from package",
			pkg:  pkg,
			want: []*definition.StructureDefinition{sd},
		}, {
			name: "structure definition is not from package",
			pkg:  registry.NewPackageRef("default", "hl7.fhir.r4.core", "4.0.2"),
			want: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := module.FilterStructureDefinitions(tc.pkg)

			if got, want := got, tc.want; !cmp.Equal(got, want, cmpopts.EquateEmpty()) {
				t.Errorf("Module.FilterStructureDefinitions(%v) = %v, want %v", tc.pkg, got, want)
			}
		})
	}
}

func TestModuleFilterValueSets(t *testing.T) {
	sd := mustReadJSON[definition.StructureDefinition](t, "testdata/structure-definition.json")
	cs := mustReadJSON[definition.CodeSystem](t, "testdata/code-system.json")
	vs := mustReadJSON[definition.ValueSets](t, "testdata/value-set.json")
	pkg := registry.NewPackageRef("default", "hl7.fhir.r4.core", "4.0.1")
	sdSource := &conformance.Source{Package: pkg}
	csSource := &conformance.Source{Package: pkg}
	vsSource := &conformance.Source{Package: pkg}
	module := conformance.DefaultModule()
	module.AddDefinition(sd, sdSource)
	module.AddDefinition(cs, csSource)
	module.AddDefinition(vs, vsSource)

	testCases := []struct {
		name string
		pkg  registry.PackageRef
		want []*definition.ValueSets
	}{
		{
			name: "value set is from package",
			pkg:  pkg,
			want: []*definition.ValueSets{vs},
		}, {
			name: "value set is not from package",
			pkg:  registry.NewPackageRef("default", "hl7.fhir.r4.core", "4.0.2"),
			want: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := module.FilterValueSets(tc.pkg)

			if got, want := got, tc.want; !cmp.Equal(got, want, cmpopts.EquateEmpty()) {
				t.Errorf("Module.FilterValueSets(%v) = %v, want %v", tc.pkg, got, want)
			}
		})
	}
}

func TestModuleFilterCodeSystems(t *testing.T) {
	sd := mustReadJSON[definition.StructureDefinition](t, "testdata/structure-definition.json")
	cs := mustReadJSON[definition.CodeSystem](t, "testdata/code-system.json")
	vs := mustReadJSON[definition.ValueSets](t, "testdata/value-set.json")
	pkg := registry.NewPackageRef("default", "hl7.fhir.r4.core", "4.0.1")
	sdSource := &conformance.Source{Package: pkg}
	csSource := &conformance.Source{Package: pkg}
	vsSource := &conformance.Source{Package: pkg}
	module := conformance.DefaultModule()
	module.AddDefinition(sd, sdSource)
	module.AddDefinition(cs, csSource)
	module.AddDefinition(vs, vsSource)

	testCases := []struct {
		name string
		pkg  registry.PackageRef
		want []*definition.CodeSystem
	}{
		{
			name: "code system is from package",
			pkg:  pkg,
			want: []*definition.CodeSystem{cs},
		}, {
			name: "code system is not from package",
			pkg:  registry.NewPackageRef("default", "hl7.fhir.r4.core", "4.0.2"),
			want: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := module.FilterCodeSystems(tc.pkg)

			if got, want := got, tc.want; !cmp.Equal(got, want, cmpopts.EquateEmpty()) {
				t.Errorf("Module.FilterCodeSystems(%v) = %v, want %v", tc.pkg, got, want)
			}
		})
	}
}

func TestModuleFilterConceptMaps(t *testing.T) {
	sd := mustReadJSON[definition.StructureDefinition](t, "testdata/structure-definition.json")
	cs := mustReadJSON[definition.CodeSystem](t, "testdata/code-system.json")
	vs := mustReadJSON[definition.ValueSets](t, "testdata/value-set.json")
	cm := mustReadJSON[definition.ConceptMap](t, "testdata/concept-map.json")
	pkg := registry.NewPackageRef("default", "hl7.fhir.r4.core", "4.0.1")
	sdSource := &conformance.Source{Package: pkg}
	csSource := &conformance.Source{Package: pkg}
	vsSource := &conformance.Source{Package: pkg}
	cmSource := &conformance.Source{Package: pkg}
	module := conformance.DefaultModule()
	module.AddDefinition(sd, sdSource)
	module.AddDefinition(cs, csSource)
	module.AddDefinition(vs, vsSource)
	module.AddDefinition(cm, cmSource)

	testCases := []struct {
		name string
		pkg  registry.PackageRef
		want []*definition.ConceptMap
	}{
		{
			name: "concept map is from package",
			pkg:  pkg,
			want: []*definition.ConceptMap{cm},
		}, {
			name: "concept map is not from package",
			pkg:  registry.NewPackageRef("default", "hl7.fhir.r4.core", "4.0.2"),
			want: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := module.FilterConceptMaps(tc.pkg)

			if got, want := got, tc.want; !cmp.Equal(got, want, cmpopts.EquateEmpty()) {
				t.Errorf("Module.FilterConceptMaps(%v) = %v, want %v", tc.pkg, got, want)
			}
		})
	}
}

func TestModuleFilterAll(t *testing.T) {
	sd := mustReadJSON[definition.StructureDefinition](t, "testdata/structure-definition.json")
	cs := mustReadJSON[definition.CodeSystem](t, "testdata/code-system.json")
	vs := mustReadJSON[definition.ValueSets](t, "testdata/value-set.json")
	cm := mustReadJSON[definition.ConceptMap](t, "testdata/concept-map.json")
	pkg := registry.NewPackageRef("default", "hl7.fhir.r4.core", "4.0.1")
	sdSource := &conformance.Source{Package: pkg}
	csSource := &conformance.Source{Package: pkg}
	vsSource := &conformance.Source{Package: pkg}
	cmSource := &conformance.Source{Package: pkg}
	module := conformance.DefaultModule()
	module.AddDefinition(sd, sdSource)
	module.AddDefinition(cs, csSource)
	module.AddDefinition(vs, vsSource)
	module.AddDefinition(cm, cmSource)

	testCases := []struct {
		name string
		pkg  registry.PackageRef
		want []definition.Canonical
	}{
		{
			name: "definition is from package",
			pkg:  pkg,
			want: []definition.Canonical{cs, sd, vs, cm},
		}, {
			name: "definition is not from package",
			pkg:  registry.NewPackageRef("default", "hl7.fhir.r4.core", "4.0.2"),
			want: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := module.FilterAll(tc.pkg)

			if got, want := got, tc.want; !cmp.Equal(got, want, cmpopts.SortSlices(sortCanonical)) {
				t.Errorf("Module.FilterAll(%v) = %v, want %v", tc.pkg, got, want)
			}
		})
	}
}

func TestModuleLookupStructureDefinition(t *testing.T) {
	sd := mustReadJSON[definition.StructureDefinition](t, "testdata/structure-definition.json")
	cs := mustReadJSON[definition.CodeSystem](t, "testdata/code-system.json")
	vs := mustReadJSON[definition.ValueSets](t, "testdata/value-set.json")
	pkg := registry.NewPackageRef("default", "hl7.fhir.r4.core", "4.0.1")
	sdSource := &conformance.Source{Package: pkg}
	csSource := &conformance.Source{Package: pkg}
	vsSource := &conformance.Source{Package: pkg}
	module := conformance.DefaultModule()
	module.AddDefinition(sd, sdSource)
	module.AddDefinition(cs, csSource)
	module.AddDefinition(vs, vsSource)

	testCases := []struct {
		name string
		url  string
		want *definition.StructureDefinition
	}{
		{
			name: "structure definition",
			url:  sd.GetURL().GetValue(),
			want: sd,
		}, {
			name: "structure definition base relative",
			url:  filepath.Base(sd.GetURL().GetValue()),
			want: sd,
		}, {
			name: "code system",
			url:  cs.GetURL().GetValue(),
			want: nil,
		}, {
			name: "value set",
			url:  vs.GetURL().GetValue(),
			want: nil,
		}, {
			name: "unknown",
			url:  "http://example.com/unknown",
			want: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := module.LookupStructureDefinition(tc.url)

			if got, want := ok, tc.want != nil; got != want {
				t.Fatalf("Module.LookupStructureDefinition(%q) ok = %v, want %v", tc.url, got, want)
			}
			if got, want := got, tc.want; !cmp.Equal(got, want) {
				t.Errorf("Module.LookupStructureDefinition(%q) = %v, want %v", tc.url, got, want)
			}
		})
	}
}

func TestModuleLookupValueSet(t *testing.T) {
	sd := mustReadJSON[definition.StructureDefinition](t, "testdata/structure-definition.json")
	cs := mustReadJSON[definition.CodeSystem](t, "testdata/code-system.json")
	vs := mustReadJSON[definition.ValueSets](t, "testdata/value-set.json")
	pkg := registry.NewPackageRef("default", "hl7.fhir.r4.core", "4.0.1")
	sdSource := &conformance.Source{Package: pkg}
	csSource := &conformance.Source{Package: pkg}
	vsSource := &conformance.Source{Package: pkg}
	module := conformance.DefaultModule()
	module.AddDefinition(sd, sdSource)
	module.AddDefinition(cs, csSource)
	module.AddDefinition(vs, vsSource)

	testCases := []struct {
		name string
		url  string
		want *definition.ValueSets
	}{
		{
			name: "structure definition",
			url:  sd.GetURL().GetValue(),
			want: nil,
		}, {
			name: "code system",
			url:  cs.GetURL().GetValue(),
			want: nil,
		}, {
			name: "value set",
			url:  vs.GetURL().GetValue(),
			want: vs,
		}, {
			name: "unknown",
			url:  "http://example.com/unknown",
			want: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := module.LookupValueSet(tc.url)

			if got, want := ok, tc.want != nil; got != want {
				t.Fatalf("Module.LookupValueSet(%q) ok = %v, want %v", tc.url, got, want)
			}
			if got, want := got, tc.want; !cmp.Equal(got, want) {
				t.Errorf("Module.LookupValueSet(%q) = %v, want %v", tc.url, got, want)
			}
		})
	}
}

func TestModuleLookupCodeSystem(t *testing.T) {
	sd := mustReadJSON[definition.StructureDefinition](t, "testdata/structure-definition.json")
	cs := mustReadJSON[definition.CodeSystem](t, "testdata/code-system.json")
	vs := mustReadJSON[definition.ValueSets](t, "testdata/value-set.json")
	pkg := registry.NewPackageRef("default", "hl7.fhir.r4.core", "4.0.1")
	sdSource := &conformance.Source{Package: pkg}
	csSource := &conformance.Source{Package: pkg}
	vsSource := &conformance.Source{Package: pkg}
	module := conformance.DefaultModule()
	module.AddDefinition(sd, sdSource)
	module.AddDefinition(cs, csSource)
	module.AddDefinition(vs, vsSource)

	testCases := []struct {
		name string
		url  string
		want *definition.CodeSystem
	}{
		{
			name: "structure definition",
			url:  sd.GetURL().GetValue(),
			want: nil,
		}, {
			name: "code system",
			url:  cs.GetURL().GetValue(),
			want: cs,
		}, {
			name: "value set",
			url:  vs.GetURL().GetValue(),
			want: nil,
		}, {
			name: "unknown",
			url:  "http://example.com/unknown",
			want: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := module.LookupCodeSystem(tc.url)

			if got, want := ok, tc.want != nil; got != want {
				t.Fatalf("Module.LookupCodeSystem(%q) ok = %v, want %v", tc.url, got, want)
			}
			if got, want := got, tc.want; !cmp.Equal(got, want) {
				t.Errorf("Module.LookupCodeSystem(%q) = %v, want %v", tc.url, got, want)
			}
		})
	}
}

func TestModuleLookupConceptMap(t *testing.T) {
	sd := mustReadJSON[definition.StructureDefinition](t, "testdata/structure-definition.json")
	cs := mustReadJSON[definition.CodeSystem](t, "testdata/code-system.json")
	vs := mustReadJSON[definition.ValueSets](t, "testdata/value-set.json")
	cm := mustReadJSON[definition.ConceptMap](t, "testdata/concept-map.json")
	pkg := registry.NewPackageRef("default", "hl7.fhir.r4.core", "4.0.1")
	sdSource := &conformance.Source{Package: pkg}
	csSource := &conformance.Source{Package: pkg}
	vsSource := &conformance.Source{Package: pkg}
	cmSource := &conformance.Source{Package: pkg}
	module := conformance.DefaultModule()
	module.AddDefinition(sd, sdSource)
	module.AddDefinition(cs, csSource)
	module.AddDefinition(vs, vsSource)
	module.AddDefinition(cm, cmSource)

	testCases := []struct {
		name string
		url  string
		want *definition.ConceptMap
	}{
		{
			name: "structure definition",
			url:  sd.GetURL().GetValue(),
			want: nil,
		}, {
			name: "code system",
			url:  cs.GetURL().GetValue(),
			want: nil,
		}, {
			name: "value set",
			url:  vs.GetURL().GetValue(),
			want: nil,
		}, {
			name: "concept map",
			url:  cm.GetURL().GetValue(),
			want: cm,
		}, {
			name: "unknown",
			url:  "http://example.com/unknown",
			want: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := module.LookupConceptMap(tc.url)

			if got, want := ok, tc.want != nil; got != want {
				t.Fatalf("Module.LookupConceptMap(%q) ok = %v, want %v", tc.url, got, want)
			}
			if got, want := got, tc.want; !cmp.Equal(got, want) {
				t.Errorf("Module.LookupConceptMap(%q) = %v, want %v", tc.url, got, want)
			}
		})
	}
}
