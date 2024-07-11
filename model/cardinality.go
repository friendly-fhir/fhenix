package model

import (
	"strconv"

	fhir "github.com/friendly-fhir/go-fhir/r4/core"
)

const Unbound = -1

// Cardinality represents the cardinality of an element.
type Cardinality struct {
	Min int
	Max int
}

func (c *Cardinality) String() string {
	if c.Max == Unbound {
		return strconv.Itoa(c.Min) + "..*"
	}
	return strconv.Itoa(c.Min) + ".." + strconv.Itoa(c.Max)
}

func (c *Cardinality) IsRequired() bool {
	return c.Min >= 1
}

func (c *Cardinality) IsDisabled() bool {
	return c.Max == 0
}

func (c *Cardinality) IsScalar() bool {
	return c.Max == 1
}

func (c *Cardinality) IsOptional() bool {
	return c.Min == 0 && c.IsScalar()
}

func (c *Cardinality) IsList() bool {
	return c.Max > 1 || c.IsUnboundedList()
}

func (c *Cardinality) IsUnboundedList() bool {
	return c.Max == -1
}

func (c *Cardinality) FromElementDefinition(ed *fhir.ElementDefinition) (err error) {
	c.Min = int(ed.GetMin().GetValue())
	c.Max = 1
	if ed.GetMax().GetValue() == "*" {
		c.Max = -1
	} else {
		c.Max, err = strconv.Atoi(ed.GetMax().GetValue())
	}
	return
}

func (c *Cardinality) FromBaseElementDefinition(ed *fhir.ElementDefinition) (err error) {
	if ed.Base == nil {
		return c.FromElementDefinition(ed)
	}
	c.Min = int(ed.GetBase().GetMin().GetValue())
	c.Max = 1
	if ed.GetBase().GetMax().GetValue() == "*" {
		c.Max = -1
	} else {
		c.Max, err = strconv.Atoi(ed.GetBase().GetMax().GetValue())
	}
	return
}
