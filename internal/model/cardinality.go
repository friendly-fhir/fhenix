package model

import (
	"fmt"
	"strings"
)

type MaxCardinality int

const (
	Unbound MaxCardinality = -1
)

type Cardinality struct {
	Min int
	Max MaxCardinality
}

func (c Cardinality) IsScalar() bool {
	return c.Min == 1 && c.Max == 1
}

func (c Cardinality) IsOptional() bool {
	return c.Min == 0 && c.Max == 1
}

func (c Cardinality) IsList() bool {
	return c.Min == 0 && c.Max > 1
}

func (c Cardinality) IsUnboundedList() bool {
	return c.Min == 0 && c.Max == Unbound
}

func (c *Cardinality) String() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d..", c.Min)
	if c.Max == Unbound {
		fmt.Fprintf(&sb, "*")
	} else {
		fmt.Fprintf(&sb, "%d", int(c.Max))
	}
	return sb.String()
}
