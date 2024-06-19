package model

import "regexp"

type Type struct {
	// Package is the package name of the type definition
	Package string

	// Version is the version of the package
	Version string

	// Status is the publication status of the type definition
	Status string

	// Name represents the actual name of the defined type
	Name string

	// Kind represents the kind of type definition
	Kind string

	// The base type (if one exists)
	Base *Type

	// Representation is how the type is represented when serialized. By default
	// all types are serialized as structures; but primitive types have inline
	// representations instead.
	Representation Repr

	// Regexp is the regular expression that the type must match.
	Regex *regexp.Regexp

	// Short is the short description provided to the structure definition.
	Short string

	// Description is the full description provided to the structure definition.
	Description string

	// Derivation is the derivation type of the structure definition.
	Derivation string

	// Comment represents a comment that has been applied to a structure definition.
	Comment string

	// IsAbstract is set to true if a type is an abstract definition and not a
	// concrete object instance.
	IsAbstract bool

	// URL is the string URL of the type definition
	URL string

	// Fields is a mapping of all the field names to the underlying type
	// definition.
	Fields map[string]*Field

	// DerivedTypes is a list of all the types that are derived from this type.
	DerivedTypes []*Type

	// NestedTypes is a list of all the child types that are defined in this type.
	// This corresponds to nested types in the structure definition, such as
	// backbone elements.
	NestedTypes []*Type
}

func (t *Type) IsConstraint() bool {
	return t.Derivation == "constraint"
}

func (t *Type) IsSpecialization() bool {
	return t.Derivation == "specialization"
}

func (t *Type) IsResource() bool {
	return t.Kind == "resource"
}

func (t *Type) IsPrimitive() bool {
	return t.Kind == "primitive-type"
}

func (t *Type) IsComplex() bool {
	return t.Kind == "complex-type"
}

// IsExtension returns true if the type is derived from an extension.
func (t *Type) IsExtension() bool {
	if t.Kind != "complex-type" {
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

// HasDerived returns true if the type has any derived types.
func (t *Type) HasDerived() bool {
	return len(t.DerivedTypes) > 0
}

// FieldNames returns a slice containing all the field names that this Type
// defines.
func (t *Type) FieldNames() []string {
	keys := make([]string, 0, len(t.Fields))
	for key := range t.Fields {
		keys = append(keys, key)
	}
	return keys
}

// LookupField looks up a field by name in the type definition. If the field
// does not exist, ok will be false.
func (t *Type) LookupField(name string) (field *Field, ok bool) {
	field, ok = t.Fields[name]
	return
}

// Field returns the field definition for the given field name. If the field
// does not exist, this will return nil.
func (t *Type) Field(name string) *Field {
	return t.Fields[name]
}
