package config

import (
	"errors"
	"fmt"

	"github.com/friendly-fhir/fhenix/internal/template"
	"gopkg.in/yaml.v3"
)

type Version int

// Transformation outlines a transformation process that takes an input and
// produces an output.
type Transformation struct {
	Version  Version  `yaml:"version"`
	Input    Input    `yaml:"input"`
	Output   string   `yaml:"output"`
	Template Template `yaml:"template"`
}

func (v *Version) UnmarshalYAML(node *yaml.Node) error {
	var s string
	if err := node.Decode(&s); err != nil {
		return err
	}

	switch s {
	case "1":
		*v = 1
	default:
		return fmt.Errorf("invalid version: %v", s)
	}

	return nil
}

// Type is a configuration node that specifies the type of input to process.
type Type string

const (
	TypeStructureDefinition Type = "StructureDefinition"
	TypeValueSet            Type = "ValueSet"
	TypeCodeSystem          Type = "CodeSystem"
)

func (t *Type) UnmarshalYAML(node *yaml.Node) error {
	var s string
	if err := node.Decode(&s); err != nil {
		return err
	}

	switch s {
	case "StructureDefinition":
		*t = TypeStructureDefinition
	case "ValueSet":
		*t = TypeValueSet
	case "CodeSystem":
		*t = TypeCodeSystem
	default:
		return errors.New("invalid type")
	}

	return nil
}

var _ yaml.Unmarshaler = (*Type)(nil)

// Conditions is a configuration node that specifies conditions that must be met for a
// process to run.
type Condition struct {
	tmpl *template.Template
}

// NewCondition creates a new condition with the given template.
func NewCondition(tmpl *template.Template) *Condition {
	return &Condition{tmpl: tmpl}
}

// Evaluate evaluates the conditions against the given data and returns true if
// all of the conditions are met.
func (c *Condition) Evaluate(data any) bool {
	if c.tmpl == nil {
		return true
	}
	truth, _ := c.tmpl.ExecuteBool(data)
	return truth
}

func (c *Condition) UnmarshalYAML(node *yaml.Node) error {
	switch node.Kind {
	case yaml.ScalarNode:
		var tmpl template.Template
		err := node.Decode(&tmpl)
		if err != nil {
			return err
		}
		c.tmpl = &tmpl
	default:
		return errors.New("invalid if")
	}
	return nil
}

var _ yaml.Unmarshaler = (*Condition)(nil)

type Input struct {
	// Type is the type of input to process.
	Type Type `yaml:"type"`

	// If is a condition that must be met for the process to run.
	If Condition `yaml:"if,omitempty"`
}

// Template is a struct that holds the header, content, and footer paths for a
// template.
type Template struct {
	// Header is the path to the header template. May be omitted.
	// Headers are used only on the first invocation of the template for any given
	// file path.
	Header string `yaml:"header,omitempty"`

	// Content is the path to the content template.
	// Content is used on every invocation of the template for any given file path.
	Content string `yaml:"content,omitempty"`

	// Footer is the path to the footer template. May be omitted.
	// Footers are used only on the last invocation of the template for any given
	// file path.
	Footer string `yaml:"footer,omitempty"`
}

type Package struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
}

type Config struct {
	Package         Package          `yaml:"package"`
	Transformations []Transformation `yaml:"transformation"`
}
