package model

import (
	"cmp"
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/friendly-fhir/fhenix/internal/fhirig"
	"github.com/friendly-fhir/fhenix/internal/model/raw"
)

// ModelBuilder is a helper type that assists in building a model from the
// raw definitions.
type ModelBuilder struct {
	pkg *fhirig.Package

	structureDefinitions map[string]*definition[raw.StructureDefinition]
	codeSystems          map[string]*definition[raw.CodeSystem]
	valueSets            map[string]*definition[raw.ValueSet]

	// base is the base URL of the FHIR IG that is being processed.
	base string
}

// definition is a helper type that is used to store the package and the
type definition[T any] struct {
	Package    *fhirig.Package
	Definition *T
}

// NewModelBuilder creates a new model builder instance.
func NewModelBuilder(pkg *fhirig.Package) *ModelBuilder {
	return &ModelBuilder{
		pkg: pkg,

		structureDefinitions: map[string]*definition[raw.StructureDefinition]{},
		codeSystems:          map[string]*definition[raw.CodeSystem]{},
		valueSets:            map[string]*definition[raw.ValueSet]{},
	}
}

func sortedkeys[T cmp.Ordered, V any](m map[T]V) []T {
	keys := make([]T, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	return keys
}

func filterKind(m map[string]*definition[raw.StructureDefinition], kind string) map[string]*definition[raw.StructureDefinition] {
	result := map[string]*definition[raw.StructureDefinition]{}
	for key, t := range m {
		if t.Definition.Kind == kind {
			result[key] = t
		}
	}
	return result
}

func filterDerivation(m map[string]*definition[raw.StructureDefinition], derivation string) map[string]*definition[raw.StructureDefinition] {
	result := map[string]*definition[raw.StructureDefinition]{}
	for key, t := range m {
		if t.Definition.Derivation == derivation {
			result[key] = t
		}
	}
	return result
}

// Build converts this builder into a full Model of the FHIR IG.
func (m *ModelBuilder) Build() (*Model, error) {
	model := Model{
		types: map[string]*Type{},
		codes: map[string]*CodeSystem{},
		base:  m.baseIG(),
		pkg:   m.pkg,
	}
	sds := filterDerivation(m.structureDefinitions, "specialization")
	for _, url := range sortedkeys(sds) {
		if _, err := m.typeFromStructureDef("", &model, url); err != nil {
			return nil, err
		}
	}
	// for _, url := range sortedkeys(filterKind(sds, "complex-type")) {
	// 	if _, err := m.typeFromStructureDef("", &model, url); err != nil {
	// 		return nil, err
	// 	}
	// }
	// for _, url := range sortedkeys(filterKind(sds, "resource")) {
	// 	if _, err := m.typeFromStructureDef("", &model, url); err != nil {
	// 		return nil, err
	// 	}
	// }
	for _, url := range sortedkeys(m.codeSystems) {
		if _, err := m.typeFromCodeSystem(&model, url); err != nil {
			return nil, err
		}
	}
	return &model, nil
}

func (m *ModelBuilder) typeFromStructureDef(prefix string, model *Model, url string) (*Type, error) {
	const (
		fhirpathPrefix = "http://hl7.org/fhirpath/System."
	)
	if t, ok := model.LookupType(url); ok {
		return t, nil
	}
	if t, ok := model.LookupType(m.baseStructureDefinition() + url); ok {
		return t, nil
	}

	if ty, ok := strings.CutPrefix(url, fhirpathPrefix); ok {
		return &Type{
			Name: strings.ToLower(ty),
			Kind: "builtin",
		}, nil
	}

	sd, ok := m.structureDefinitions[url]
	if !ok {
		url = m.baseStructureDefinition() + url
		sd, ok = m.structureDefinitions[url]
	}
	if !ok {
		return nil, fmt.Errorf("unable to find structure definition for %s", url)
	}
	def := &Type{
		Name:        sd.Definition.Name,
		Base:        nil,
		Kind:        sd.Definition.Kind,
		Status:      sd.Definition.Status,
		IsAbstract:  sd.Definition.Abstract,
		Package:     sd.Package.Name(),
		Version:     sd.Package.Version(),
		Short:       sd.Definition.Short,
		Derivation:  sd.Definition.Derivation,
		Description: sd.Definition.Description,
		Comment:     sd.Definition.Comment,
		URL:         url,
		Fields:      map[string]*Field{},
	}
	model.types[sd.Definition.URL] = def

	if sd.Definition.BaseDefinition != "" {
		if base, ok := model.LookupType(sd.Definition.BaseDefinition); ok {
			def.Base = base
		} else {
			base, err := m.typeFromStructureDef(prefix+"-", model, sd.Definition.BaseDefinition)
			if err != nil {
				return nil, err
			}
			def.Base = base
		}
		def.Base.DerivedTypes = append(def.Base.DerivedTypes, def)
	}
	if err := m.buildFields(prefix, model, def, sd.Definition); err != nil {
		return nil, err
	}
	return def, nil
}

type FieldPath string

func (f FieldPath) LastField() string {
	parts := strings.Split(string(f), ".")
	if len(parts) <= 1 {
		return string(f)
	}
	return parts[len(parts)-1]
}

func (f FieldPath) Parent() FieldPath {
	parts := strings.Split(string(f), ".")
	if len(parts) <= 1 {
		return ""
	}
	return FieldPath(strings.Join(parts[:len(parts)-1], "."))
}

func (f FieldPath) Length() int {
	return strings.Count(string(f), ".")
}

func (m *ModelBuilder) buildFields(prefix string, model *Model, t *Type, sd *raw.StructureDefinition) error {
	types := map[FieldPath]*Type{}
	for _, elem := range sd.Snapshot.Element {
		if elem.Path == t.Name {
			types[FieldPath(elem.Path)] = t
			continue
		}
		path := FieldPath(elem.Path)
		parent, field := path.Parent(), path.LastField()
		if len(elem.Types) == 0 {
			continue
		}

		var err error
		var ty *Type
		if elem.Types[0].Code == "BackboneElement" {
			backbone, err := m.typeFromStructureDef(prefix+"-", model, elem.Types[0].Code)
			if err != nil {
				return err
			}

			ty = &Type{
				Name:           elem.Path,
				Package:        t.Package,
				Version:        t.Version,
				Status:         t.Status,
				Kind:           "complex-type",
				Base:           backbone,
				Representation: ReprStruct,
				Short:          elem.Short,
				Description:    elem.Definition,
				Derivation:     t.Derivation,
				Comment:        elem.Comment,
				IsAbstract:     false,
				Fields:         map[string]*Field{},
			}
		} else {
			const (
				typeURL = "http://hl7.org/fhir/StructureDefinition/structuredefinition-fhir-type"
				extURL  = "http://hl7.org/fhir/StructureDefinition/regex"
			)
			if ext := elem.Types[0].Extension(typeURL); ext != nil {
				ty = &Type{
					Name: ext.ValueURL,
					Kind: "builtin",
				}
				if ext := elem.Types[0].Extension(extURL); ext != nil {
					ty.Regex, err = regexp.Compile(ext.ValueCode)
					if err != nil {
						return fmt.Errorf("unable to compile regular expression %s: %w", ext.ValueCode, err)
					}
				}
			} else {
				ty, err = m.typeFromStructureDef(prefix+"-", model, elem.Types[0].Code)
				if err != nil {
					return err
				}
			}
		}
		types[path] = ty
		var cardinality Cardinality
		cardinality.Min = elem.Min
		if elem.Max == "*" {
			cardinality.Max = Unbound
		} else {
			m, err := strconv.Atoi(elem.Max)
			if err != nil {
				return fmt.Errorf("unable to parse max cardinality %s: %w", elem.Max, err)
			}
			cardinality.Max = MaxCardinality(m)
		}
		types[parent].Fields[field] = &Field{
			Name:        path.LastField(),
			Path:        elem.Path,
			Type:        ty,
			Cardinality: cardinality,
			Definition:  elem.Definition,
			Comment:     elem.Comment,
		}
	}
	return nil
}

func (m *ModelBuilder) typeFromCodeSystem(model *Model, url string) (*CodeSystem, error) {
	if cs, ok := model.LookupCodeSystem(url); ok {
		return cs, nil
	}
	cs, ok := m.codeSystems[url]
	if !ok {
		return nil, fmt.Errorf("unable to find code system for %s", url)
	}
	def := &CodeSystem{
		Package:     cs.Package.Name(),
		Version:     cs.Package.Version(),
		Description: cs.Definition.Description,
		URL:         cs.Definition.URL,
		Name:        cs.Definition.Name,
		Title:       cs.Definition.Title,
		Status:      cs.Definition.Status,
		Codes:       make([]Code, 0, len(cs.Definition.Concept)),
	}
	for _, code := range cs.Definition.Concept {
		def.Codes = append(def.Codes, Code{
			Value:      code.Code,
			Display:    code.Display,
			Definition: code.Definition,
		})
	}
	model.codes[cs.Definition.URL] = def
	return def, nil
}

func (m *ModelBuilder) baseIG() string {
	if m.base == "" {
		return "http://hl7.org/fhir/"
	}
	return m.base
}

func (m *ModelBuilder) baseStructureDefinition() string {
	return m.baseIG() + "StructureDefinition/"
}

func (m *ModelBuilder) AddStructureDefinition(pkg *fhirig.Package, sd *raw.StructureDefinition) {
	m.structureDefinitions[sd.URL] = &definition[raw.StructureDefinition]{
		Package:    pkg,
		Definition: sd,
	}
}

func (m *ModelBuilder) AddCodeSystem(pkg *fhirig.Package, cs *raw.CodeSystem) {
	m.codeSystems[cs.URL] = &definition[raw.CodeSystem]{
		Package:    pkg,
		Definition: cs,
	}
}

func (m *ModelBuilder) AddValueSet(pkg *fhirig.Package, vs *raw.ValueSet) {
	m.valueSets[vs.URL] = &definition[raw.ValueSet]{
		Package:    pkg,
		Definition: vs,
	}
}
