package model

import (
	"slices"
	"strings"
	"sync"

	"github.com/friendly-fhir/fhenix/pkg/registry"
)

type TypeSet struct {
	base  string
	types *sync.Map
}

func NewTypeSet(base string, entries ...*Type) *TypeSet {
	ts := &TypeSet{
		base:  base,
		types: &sync.Map{},
	}
	if !strings.HasSuffix(ts.base, "/") {
		ts.base += "/"
	}
	for _, t := range entries {
		ts.Add(t)
	}
	return ts
}

func DefaultTypeSet() *TypeSet {
	return NewTypeSet("http://hl7.org/fhir/StructureDefinition")
}

// Base returns the base URL of the structure definition IG in the typeset.
func (ts *TypeSet) Base() string {
	return strings.TrimSuffix(ts.base, "/")
}

func (ts *TypeSet) Add(t *Type) {
	ts.types.Store(ts.join(t.URL), t)
}

func (ts *TypeSet) Lookup(url string) (*Type, bool) {
	t, ok := ts.types.Load(url)
	if !ok {
		t, ok = ts.types.Load(ts.join(url))
	}
	if !ok {
		return nil, false
	}
	return t.(*Type), true
}

func (ts *TypeSet) Get(url string) *Type {
	t, _ := ts.Lookup(url)
	return t
}

func (ts *TypeSet) TypesMatching(condition func(*Type) bool) *TypeSet {
	result := NewTypeSet(ts.base)
	ts.types.Range(func(_, value any) bool {
		t := value.(*Type)
		if condition(t) {
			result.types.Store(t.URL, t)
		}
		return true
	})
	return result
}

func (ts *TypeSet) DefinedInPackage(pkg registry.PackageRef) *TypeSet {
	return ts.TypesMatching(func(t *Type) bool {
		return t.Source.Package == pkg
	})
}

func (ts *TypeSet) All() []*Type {
	var all []*Type
	ts.types.Range(func(_, value any) bool {
		all = append(all, value.(*Type))
		return true
	})
	slices.SortFunc(all, func(lhs, rhs *Type) int {
		return strings.Compare(lhs.Name, rhs.Name)
	})
	return all
}

func (ts *TypeSet) Resources() *TypeSet {
	return ts.TypesMatching(func(t *Type) bool { return t.Kind == TypeKindResource })
}

func (ts *TypeSet) ComplexTypes() *TypeSet {
	return ts.TypesMatching(func(t *Type) bool { return t.Kind == TypeKindComplexType })
}

func (ts *TypeSet) PrimitiveTypes() *TypeSet {
	return ts.TypesMatching(func(t *Type) bool { return t.Kind == TypeKindPrimitive })
}

func (ts *TypeSet) InBase() *TypeSet {
	return ts.TypesMatching(func(t *Type) bool { return strings.HasPrefix(t.URL, ts.base) })
}

func (ts *TypeSet) join(url string) string {
	if strings.Contains(url, "/") {
		return url
	}
	return ts.base + url
}
