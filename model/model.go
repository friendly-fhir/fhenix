package model

import (
	"cmp"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/friendly-fhir/fhenix/internal/fhirig"
	"github.com/friendly-fhir/fhenix/model/conformance"
	"github.com/friendly-fhir/fhenix/model/conformance/definition"
	fhir "github.com/friendly-fhir/go-fhir/r4/core"
)

type Model struct {
	module  *conformance.Module
	types   *TypeSet
	defined bool
}

func NewModel(module *conformance.Module) *Model {
	ts := NewTypeSet(module.Base() + "/StructureDefinition")
	return &Model{
		module: module,
		types:  ts,
	}
}

func (m *Model) DefineType(url string) error {
	if m.defined {
		return nil
	}
	if _, ok := m.types.Lookup(url); ok {
		return nil
	}

	entry, ok := m.module.LookupStructureDefinition(url)
	if !ok {
		return fmt.Errorf("structure definition %q not found", url)
	}
	src := m.module.SourceOf(entry)

	return m.typeFromStructureDef(src.Package, src.File, entry)
}

func (m *Model) DefineAllTypes() error {
	if m.defined {
		return nil
	}
	var errs []error
	for _, sd := range m.module.StructureDefinitions() {
		errs = append(errs, m.DefineType(sd.GetURL().GetValue()))
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

func (m *Model) typeFromStructureDef(pkg *fhirig.Package, file string, sd *definition.StructureDefinition) error {
	t := &Type{
		Source: &TypeSource{
			Package:             pkg,
			File:                file,
			StructureDefinition: sd,
		},
		Name: sd.GetName().GetValue(),
		// Short:       sd.GetShort().GetValue(),
		// Comment:     sd.GetComment().GetValue(),
		Description: sd.GetDescription().GetValue(),

		URL:        sd.GetURL().GetValue(),
		Kind:       TypeKind(sd.GetKind().GetValue()),
		IsAbstract: sd.GetAbstract().GetValue(),
	}

	m.types.Add(t)

	if sd.GetBaseDefinition().GetValue() != "" {
		base, err := m.Type(sd.GetBaseDefinition().GetValue())
		if err != nil {
			return err
		}
		t.Base = base
		base.Derived = append(base.Derived, t)
	}
	if err := m.typeFromElements(t, sd.GetSnapshot().GetElement()); err != nil {
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

func (m *Model) fieldFromElement(t *Type, field *Field, elem *fhir.ElementDefinition) error {
	if len(elem.Type) == 0 {
		return fmt.Errorf("element %q has no types", elem.GetPath().GetValue())
	}
	if len(elem.Type) == 1 {
		return m.scalarFieldFromType(t, field, elem.Type[0])
	}

	for _, t := range elem.Type {
		ty, err := m.Type(t.GetCode().GetValue())
		if err != nil {
			return err
		}
		field.Alternatives = append(field.Alternatives, ty)
	}
	return nil
}

func (m *Model) scalarFieldFromType(t *Type, field *Field, rawType *fhir.ElementDefinitionType) error {
	var builtin Builtin
	if err := builtin.FromType(rawType); err == nil {
		field.Builtin = &builtin
		return nil
	}
	base, err := m.Type(rawType.GetCode().GetValue())
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
	if ty, err := m.Type(rawType.GetCode().GetValue()); err == nil {
		field.Type = ty
		return nil
	}
	return fmt.Errorf("unknown type: %q", rawType.GetCode().GetValue())
}

func (m *Model) typeFromElements(t *Type, elems []*fhir.ElementDefinition) error {
	elems = slices.DeleteFunc(elems, func(elem *fhir.ElementDefinition) bool {
		return !strings.HasPrefix(elem.GetPath().GetValue(), t.Name+".")
	})
	// These should already be sorted -- but this technically is not a requirement
	// in the FHIR specification. Better safe than sorry.
	slices.SortFunc(elems, func(lhs, rhs *fhir.ElementDefinition) int {
		return cmp.Compare(lhs.GetPath().GetValue(), rhs.GetPath().GetValue())
	})

	for _, elem := range elems {
		if len(elem.Type) == 0 {
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
		name := m.fieldname(elem.GetPath().GetValue())
		field := &Field{
			Name:            name,
			Path:            elem.GetPath().GetValue(),
			Short:           elem.GetShort().GetValue(),
			Comment:         elem.GetComment().GetValue(),
			Definition:      elem.GetDefinition().GetValue(),
			Cardinality:     cardinality,
			BaseCardinality: baseCardinality,
		}
		if err := m.fieldFromElement(t, field, elem); err != nil {
			return err
		}
		if t.Package() == "hl7.fhir.r4.core" && elem.GetPath().GetValue() == "unsignedInt.value" || elem.GetPath().GetValue() == "positiveInt.value" {
			field.Builtin = &Builtin{
				Name: m.fieldpath(elem.GetPath().GetValue()),
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
