package model

type Field struct {
	Name string
	Path string

	Short      string
	Comment    string
	Definition string

	Type         *Type
	Alternatives []*Type
	Builtin      *Builtin

	Cardinality     Cardinality
	BaseCardinality Cardinality
}

func (f *Field) IsScalar() bool {
	return f.Cardinality.IsScalar()
}

func (f *Field) IsOptional() bool {
	return f.Cardinality.IsOptional()
}

func (f *Field) IsList() bool {
	return f.Cardinality.IsList()
}

func (f *Field) IsUnboundedList() bool {
	return f.Cardinality.IsUnboundedList()
}

func (f *Field) IsDisabled() bool {
	return f.Cardinality.IsDisabled()
}

func (f *Field) IsRequired() bool {
	return f.Cardinality.IsRequired()
}

func (f *Field) IsNarrowed() bool {
	return f.Cardinality.Min > f.BaseCardinality.Min || f.Cardinality.Max < f.BaseCardinality.Max
}
