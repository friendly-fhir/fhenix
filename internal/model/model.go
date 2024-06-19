package model

import (
	"slices"
	"strings"

	"github.com/friendly-fhir/fhenix/internal/fhirig"
)

type Model struct {
	pkg   *fhirig.Package
	base  string
	types map[string]*Type
	codes map[string]*CodeSystem
}

func (m *Model) BaseIG() string {
	return m.base
}

func (m *Model) LookupType(url string) (*Type, bool) {
	t, ok := m.types[url]
	return t, ok
}

func (m *Model) Type(url string) *Type {
	return m.types[url]
}

func (m *Model) CodeSystem(url string) *CodeSystem {
	return m.codes[url]
}

func (m *Model) LookupCodeSystem(url string) (*CodeSystem, bool) {
	cs, ok := m.codes[url]
	return cs, ok
}

func (m *Model) CodeSystems() []*CodeSystem {
	result := make([]*CodeSystem, 0, len(m.codes))
	for _, cs := range m.codes {
		result = append(result, cs)
	}
	slices.SortFunc(result, func(lhs, rhs *CodeSystem) int {
		return strings.Compare(lhs.Name, rhs.Name)
	})
	result = slices.DeleteFunc(result, func(cs *CodeSystem) bool {
		return cs.Status != "active"
	})
	return result
}

// Types returns a list of all the types defined in the modeled IG.
// This will not include foreign types defined in other IGs.
func (m *Model) Types() []*Type {
	result := slices.DeleteFunc(m.AllTypes(), func(t *Type) bool {
		return t.Package != m.pkg.Name() || t.Version != m.pkg.Version()
	})
	result = slices.DeleteFunc(result, func(t *Type) bool {
		return t.Status != "active"
	})
	return result
}

// AllTypes returns a list of all the types defined in the modeled IG.
// This will include foreign types defined in other IGs.
func (m *Model) AllTypes() []*Type {
	result := make([]*Type, 0, len(m.types))
	for _, t := range m.types {
		result = append(result, t)
	}

	slices.SortFunc(result, func(lhs, rhs *Type) int {
		return strings.Compare(lhs.Name, rhs.Name)
	})
	return result
}

func (m *Model) AbstractTypes() []*Type {
	var result []*Type
	for _, t := range m.types {
		if t.IsAbstract {
			result = append(result, t)
		}
	}
	return result
}

func (m *Model) ConcreteTypes() []*Type {
	var result []*Type
	for _, t := range m.types {
		if !t.IsAbstract {
			result = append(result, t)
		}
	}
	return result
}

func (m *Model) DirectChildren(base *Type) []*Type {
	var result []*Type
	for _, t := range m.types {
		if t.Base == base {
			result = append(result, t)
		}
	}
	return result
}
