package model

type Field struct {
	// Name represents the name of the field type
	Name string

	// Path is the key path to the field in the type
	Path string

	// Type is a reference to the underlying type of this field
	Type *Type

	Definition string

	Comment string

	// Cardinality is the cardinality of the defined type.
	Cardinality Cardinality
}

func (f *Field) IsBuiltin() bool {
	return f.Type.Kind == "builtin"
}
