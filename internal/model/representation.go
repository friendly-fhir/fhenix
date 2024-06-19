package model

type Repr int

const (
	ReprStruct Repr = iota
	ReprBuiltin
	ReprNumber
	ReprBoolean
	ReprBase64
)
