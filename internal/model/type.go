package model

import (
	"strings"

	"github.com/friendly-fhir/fhenix/internal/fhirig"
	"github.com/friendly-fhir/fhenix/internal/model/raw"
)

type TypeSource struct {
	Package             *fhirig.Package
	File                string
	StructureDefinition *raw.StructureDefinition
}

type TypeKind string

const (
	TypeKindResource    TypeKind = "resource"
	TypeKindPrimitive   TypeKind = "primitive-type"
	TypeKindComplexType TypeKind = "complex-type"
	TypeKindBackbone    TypeKind = "backbone"
)

type Type struct {
	Source *TypeSource

	Short       string
	Comment     string
	Description string

	URL  string
	Name string
	Kind TypeKind

	Base    *Type
	Derived []*Type

	IsAbstract bool

	Fields   []*Field
	SubTypes []*Type
}

func CommonBase(types []*Type) *Type {
	if len(types) == 0 {
		return nil
	}
	base := types[0]
	for _, t := range types[1:] {
		base = commonBase(base, t)
	}
	return base
}

func commonBase(t1, t2 *Type) *Type {
	if t1 == t2 {
		return t1
	}
	if t1.Base == nil && t2.Base == nil {
		return nil
	}
	if t1.Base != nil {
		t1 = t1.Base
	}
	if t2.Base != nil {
		t2 = t2.Base
	}
	return commonBase(t1, t2)
}

func (t *Type) InPackage(pkg string) bool {
	parts := strings.Split(pkg, "@")
	if len(parts) == 1 {
		return t.Package() == pkg
	}
	if len(parts) == 2 {
		return t.Package() == parts[0] && t.Version() == parts[1]
	}
	return false
}

func (t *Type) Package() string {
	return t.Source.Package.Name()
}

func (t *Type) Version() string {
	return t.Source.Package.Version()
}

func (t *Type) HasDerived() bool {
	return len(t.Derived) > 0
}

func (t *Type) IsConstraint() bool {
	return t.Source.StructureDefinition.Derivation == "constraint"
}

func (t *Type) IsSpecialization() bool {
	return t.Source.StructureDefinition.Derivation == "specialization"
}

func (t *Type) IsResource() bool {
	return t.Kind == TypeKindResource
}

func (t *Type) IsPrimitive() bool {
	return t.Kind == TypeKindPrimitive
}

func (t *Type) IsComplex() bool {
	return t.Kind == TypeKindComplexType
}

func (t *Type) IsExtension() bool {
	if t.Kind != TypeKindComplexType || !t.IsConstraint() {
		return false
	}
	if t.URL == "http://hl7.org/fhir/StructureDefinition/Extension" || t.URL == "Extension" {
		return true
	}
	if t.Base != nil {
		return t.Base.IsExtension()
	}
	return false
}
