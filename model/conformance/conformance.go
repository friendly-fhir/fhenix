/*
Package conformance provides the FHIR Conformance Module definition.

Types within this package are raw FHIR definitions that contain relevant
source information with it.
*/
package conformance

import (
	"fmt"
	"io"
	"slices"
	"strings"

	"github.com/friendly-fhir/fhenix/internal/fhirig"
	"github.com/friendly-fhir/fhenix/model/conformance/definition"
)

// Source is the source information for a FHIR definition.
type Source struct {
	Package *fhirig.Package
	File    string
}

// Module is the raw conformance module for the FHIR definitions.
type Module struct {
	base                 string
	structureDefinitions []*definition.StructureDefinition
	valueSets            []*definition.ValueSets
	codeSystems          []*definition.CodeSystem

	// source maps the canonical URL of the definition to the Source definition.
	source map[string]*source
}

type source struct {
	Canonical definition.Canonical
	Source    *Source
}

// NewModule constructs a new conformance module.
func NewModule(base string) *Module {
	base = strings.TrimSuffix(base, "/")
	return &Module{
		base:   base,
		source: map[string]*source{},
	}
}

// DefaultModule returns a new conformance module with the FHIR canonical URL.
func DefaultModule() *Module {
	return NewModule("http://hl7.org/fhir")
}

// Base returns the base URL of the conformance module.
func (m *Module) Base() string {
	return m.base
}

// ParseFile parses a file and adds the definitions to the conformance module.
func (m *Module) ParseFile(file string, pkg *fhirig.Package) error {
	canonical, err := definition.FromFile(file)
	if err != nil {
		return err
	}
	m.AddDefinition(canonical, &Source{
		Package: pkg,
		File:    file,
	})
	return nil
}

// ParseReader parses a reader and adds the definitions to the conformance
// module.
func (m *Module) ParseReader(reader io.Reader, pkg *fhirig.Package) error {
	canonical, err := definition.FromReader(reader)
	if err != nil {
		return err
	}
	m.AddDefinition(canonical, &Source{
		Package: pkg,
	})
	return nil
}

// ParseJSON parses a JSON byte slice and adds the definitions to the
// conformance module.
func (m *Module) ParseJSON(data []byte, pkg *fhirig.Package) error {
	canonical, err := definition.FromJSON(data)
	if err != nil {
		return err
	}
	m.AddDefinition(canonical, &Source{
		Package: pkg,
	})
	return nil
}

// StructureDefinitions returns all the structure definitions in the conformance
// module.
func (m *Module) StructureDefinitions() []*definition.StructureDefinition {
	result := append([]*definition.StructureDefinition(nil), m.structureDefinitions...)
	slices.SortFunc(result, sortURL)
	return result
}

// ValueSets returns all the value sets in the conformance module.
func (m *Module) ValueSets() []*definition.ValueSets {
	result := append([]*definition.ValueSets(nil), m.valueSets...)
	slices.SortFunc(result, sortURL)
	return result
}

// CodeSystems returns all the code systems in the conformance module.
func (m *Module) CodeSystems() []*definition.CodeSystem {
	result := append([]*definition.CodeSystem(nil), m.codeSystems...)
	slices.SortFunc(result, sortURL)
	return result
}

// All returns all the canonical definitions in the conformance module, sorted
// by URL.
func (m *Module) All() []definition.Canonical {
	var result []definition.Canonical
	for _, src := range m.source {
		result = append(result, src.Canonical)
	}
	slices.SortFunc(result, sortURL)
	return result
}

func sortURL[T definition.Canonical](a, b T) int {
	return strings.Compare(a.GetURL().GetValue(), b.GetURL().GetValue())
}

// FilterStructureDefinitions returns the structure definitions that are from
// the given package.
func (m *Module) FilterStructureDefinitions(pkg *fhirig.Package) []*definition.StructureDefinition {
	var result []*definition.StructureDefinition
	for _, def := range m.structureDefinitions {
		if src := m.SourceOf(def); src != nil && src.Package.String() == pkg.String() {
			result = append(result, def)
		}
	}
	slices.SortFunc(result, sortURL)
	return result
}

// FilterValueSets returns the value sets that are from the given package.
func (m *Module) FilterValueSets(pkg *fhirig.Package) []*definition.ValueSets {
	var result []*definition.ValueSets
	for _, def := range m.valueSets {
		if src := m.SourceOf(def); src != nil && src.Package.String() == pkg.String() {
			result = append(result, def)
		}
	}
	slices.SortFunc(result, sortURL)
	return result
}

// FilterCodeSystems returns the code systems that are from the given package.
func (m *Module) FilterCodeSystems(pkg *fhirig.Package) []*definition.CodeSystem {
	var result []*definition.CodeSystem
	for _, def := range m.codeSystems {
		if src := m.SourceOf(def); src != nil && src.Package.String() == pkg.String() {
			result = append(result, def)
		}
	}
	slices.SortFunc(result, sortURL)
	return result
}

func (m *Module) FilterAll(pkg *fhirig.Package) []definition.Canonical {
	var result []definition.Canonical
	for _, src := range m.source {
		if src.Source.Package.String() == pkg.String() {
			result = append(result, src.Canonical)
		}
	}
	slices.SortFunc(result, sortURL)
	return result
}

// AddDefinition adds a new definition to the conformance module.
func (m *Module) AddDefinition(canonical definition.Canonical, src *Source) {
	m.source[canonical.GetURL().GetValue()] = &source{
		Canonical: canonical,
		Source:    src,
	}
	switch def := canonical.(type) {
	case *definition.StructureDefinition:
		m.structureDefinitions = append(m.structureDefinitions, def)
	case *definition.ValueSets:
		m.valueSets = append(m.valueSets, def)
	case *definition.CodeSystem:
		m.codeSystems = append(m.codeSystems, def)
	}
}

// SourceOf returns the source information for the specific canonical definition.
func (m *Module) SourceOf(canonical definition.Canonical) *Source {
	return m.Source(canonical.GetURL().GetValue())
}

// Source returns the source information for the given URL.
func (m *Module) Source(url string) *Source {
	source, _ := m.LookupSource(url)
	return source
}

// Canonical returns the canonical definition for the given URL.
func (m *Module) Canonical(url string) definition.Canonical {
	canonical, _ := m.LookupCanonical(url)
	return canonical
}

// LookupSource returns the source information for the given URL.
func (m *Module) LookupSource(url string) (*Source, bool) {
	src, ok := m.lookup(url)
	if !ok {
		return nil, false
	}
	return src.Source, ok
}

// LookupCanonical returns the canonical definition for the given URL.
func (m *Module) LookupCanonical(url string) (definition.Canonical, bool) {
	src, ok := m.lookup(url)
	if !ok {
		return nil, false
	}
	return src.Canonical, ok
}

// LookupStructureDefinition returns the structure definition for the given URL.
func (m *Module) LookupStructureDefinition(url string) (*definition.StructureDefinition, bool) {
	src, ok := m.lookup(url)
	if !ok {
		return nil, false
	}
	sd, ok := src.Canonical.(*definition.StructureDefinition)
	return sd, ok
}

// LookupValueSet returns the value set for the given URL.
func (m *Module) LookupValueSet(url string) (*definition.ValueSets, bool) {
	src, ok := m.lookup(url)
	if !ok {
		return nil, false
	}
	vs, ok := src.Canonical.(*definition.ValueSets)
	return vs, ok
}

// LookupCodeSystem returns the code system for the given URL.
func (m *Module) LookupCodeSystem(url string) (*definition.CodeSystem, bool) {
	src, ok := m.lookup(url)
	if !ok {
		return nil, false
	}
	cs, ok := src.Canonical.(*definition.CodeSystem)
	return cs, ok
}

func (m *Module) lookup(url string) (*source, bool) {
	src, ok := m.source[url]
	if !ok {
		src, ok = m.source[fmt.Sprintf("%s/StructureDefinition/%s", m.base, url)]
	}
	if !ok {
		src, ok = m.source[fmt.Sprintf("%s/CodeSystem/%s", m.base, url)]
	}
	if !ok {
		src, ok = m.source[fmt.Sprintf("%s/ValueSet/%s", m.base, url)]
	}
	if !ok {
		src, ok = m.source[fmt.Sprintf("%s/%s", m.base, url)]
	}
	return src, ok
}

// Contains returns true if the conformance module contains the given URL.
func (m *Module) Contains(url string) bool {
	_, ok := m.source[url]
	return ok
}
