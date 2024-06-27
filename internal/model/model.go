package model

import (
	"cmp"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/friendly-fhir/fhenix/internal/fhirig"
	"github.com/friendly-fhir/fhenix/internal/model/raw"
)

type Model struct {
	cache   *FHIRCache
	types   *TypeSet
	defined bool
}

func NewModel(cache *FHIRCache) *Model {
	ts := NewTypeSet(cache.Base() + "/StructureDefinition")
	return &Model{
		cache: cache,
		types: ts,
	}
}

func (m *Model) DefineType(url string) error {
	if m.defined {
		return nil
	}
	if _, ok := m.types.Lookup(url); ok {
		return nil
	}

	entry, ok := m.cache.LookupStructureDefinition(url)
	if !ok {
		return fmt.Errorf("structure definition %q not found", url)
	}

	return m.typeFromStructureDef(entry.Package, entry.File, entry.Definition)
}

func (m *Model) DefineAllTypes() error {
	if m.defined {
		return nil
	}
	var errs []error
	for _, sd := range m.cache.StructureDefinitions() {
		errs = append(errs, m.DefineType(sd.Definition.URL))
	}
	if err := errors.Join(errs...); err != nil {
		return err
	}
	m.defined = true
	return nil
}

func (m *Model) Types() *TypeSet {
	_ = m.DefineAllTypes()
	return m.types
}

func (m *Model) CodeSystems() []*CodeSystem {
	return nil
}

func (m *Model) Type(url string) (*Type, error) {
	err := m.DefineType(url)
	if err != nil {
		return nil, err
	}
	result, ok := m.types.Lookup(url)
	if !ok {
		return nil, fmt.Errorf("unknown type: %q", url)
	}
	return result, nil
}

func (m *Model) typeFromStructureDef(pkg *fhirig.Package, file string, sd *raw.StructureDefinition) error {
	t := &Type{
		Source: &TypeSource{
			Package:             pkg,
			File:                file,
			StructureDefinition: sd,
		},
		Name:        sd.Name,
		Short:       sd.Short,
		Comment:     sd.Comment,
		Description: sd.Description,

		URL:        sd.URL,
		Kind:       TypeKind(sd.Kind),
		IsAbstract: sd.Abstract,
	}

	m.types.Add(t)

	if sd.BaseDefinition != "" {
		base, err := m.Type(sd.BaseDefinition)
		if err != nil {
			return err
		}
		t.Base = base
		base.Derived = append(base.Derived, t)
	}
	if err := m.typeFromElements(t, sd.Snapshot.Element); err != nil {
		return err
	}

	return nil
}

func (m *Model) fieldname(path string) string {
	parts := strings.Split(path, ".")
	return strings.ReplaceAll(parts[len(parts)-1], "[x]", "")
}

func (m *Model) fieldpath(path string) string {
	parts := strings.Split(path, ".")
	if len(parts) <= 1 {
		return ""
	}
	return strings.Join(parts[:len(parts)-1], ".")
}

func (m *Model) fieldFromElement(t *Type, field *Field, elem *raw.ElementDefinition) error {
	if len(elem.Types) == 0 {
		return fmt.Errorf("element %q has no types", elem.Path)
	}
	if len(elem.Types) == 1 {
		return m.scalarFieldFromType(t, field, &elem.Types[0])
	}

	for _, t := range elem.Types {
		ty, err := m.Type(t.Code)
		if err != nil {
			return err
		}
		field.Alternatives = append(field.Alternatives, ty)
	}
	return nil
}

func (m *Model) scalarFieldFromType(t *Type, field *Field, rawType *raw.Type) error {
	var builtin Builtin
	if err := builtin.FromType(rawType); err == nil {
		field.Builtin = &builtin
		return nil
	}
	base, err := m.Type(rawType.Code)
	if err != nil {
		return err
	}

	// When we find an abstract type, it means we are defining a new type in-place.
	if base.IsAbstract && base.Kind != TypeKindResource {
		ty := &Type{
			Source:      t.Source,
			Name:        field.Path,
			Short:       field.Short,
			Comment:     field.Comment,
			Description: field.Definition,
			Kind:        TypeKindBackbone,
			Base:        base,
		}
		t.SubTypes = append(t.SubTypes, ty)
		field.Type = ty
		return nil
	}
	if ty, err := m.Type(rawType.Code); err == nil {
		field.Type = ty
		return nil
	}
	return fmt.Errorf("unknown type: %q", rawType.Code)
}

func (m *Model) typeFromElements(t *Type, elems []*raw.ElementDefinition) error {
	elems = slices.DeleteFunc(elems, func(elem *raw.ElementDefinition) bool {
		return !strings.HasPrefix(elem.Path, t.Name+".")
	})
	// These should already be sorted -- but this technically is not a requirement
	// in the FHIR specification. Better safe than sorry.
	slices.SortFunc(elems, func(lhs, rhs *raw.ElementDefinition) int {
		return cmp.Compare(lhs.Path, rhs.Path)
	})

	for _, elem := range elems {
		if len(elem.Types) == 0 {
			continue
		}
		var cardinality Cardinality
		if err := cardinality.FromElementDefinition(elem); err != nil {
			return err
		}
		baseCardinality := cardinality
		if elem.Base != nil {
			if err := baseCardinality.FromBaseElementDefinition(elem); err != nil {
				return err
			}
		}
		name := m.fieldname(elem.Path)
		field := &Field{
			Name:            name,
			Path:            elem.Path,
			Short:           elem.Short,
			Comment:         elem.Comment,
			Definition:      elem.Definition,
			Cardinality:     cardinality,
			BaseCardinality: baseCardinality,
		}
		if err := m.fieldFromElement(t, field, elem); err != nil {
			return err
		}
		if t.Package() == "hl7.fhir.r4.core" && elem.Path == "unsignedInt.value" || elem.Path == "positiveInt.value" {
			field.Builtin = &Builtin{
				Name: m.fieldpath(elem.Path),
			}
		}
		m.addField(t, field)
	}
	return nil
}

func (m *Model) addField(t *Type, f *Field) {
	for _, backbones := range t.SubTypes {
		if m.fieldpath(f.Path) == backbones.Name {
			backbones.Fields = append(backbones.Fields, f)
			return
		}
	}
	t.Fields = append(t.Fields, f)
}
