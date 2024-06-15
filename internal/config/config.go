package config

import (
	"errors"

	"gopkg.in/yaml.v3"
)

// Transformation outlines a transformation process that takes an input and
// produces an output.
type Transformation struct {
	Input    Input    `yaml:"input"`
	Output   string   `yaml:"output"`
	Template Template `yaml:"template"`
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
type Conditions []string

func (c *Conditions) UnmarshalYAML(node *yaml.Node) error {
	switch node.Kind {
	case yaml.ScalarNode:
		var str string
		if err := node.Decode(&str); err == nil {
			*c = []string{str}
		}
	case yaml.SequenceNode:
		var s []string
		if err := node.Decode(&s); err != nil {
			return err
		}
		*c = s
	default:
		return errors.New("invalid if")
	}
	return nil
}

var _ yaml.Unmarshaler = (*Conditions)(nil)

type Input struct {
	// Type is the type of input to process.
	Type Type `yaml:"type"`

	// If is a list of conditions that must be met for the process to run.
	If Conditions `yaml:"if,omitempty"`
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

type Config struct {
	Transformations []Transformation `yaml:"transformation"`
}
