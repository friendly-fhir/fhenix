package config

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/friendly-fhir/fhenix/internal/template"
	"gopkg.in/yaml.v3"
)

type Version int

// Transformation outlines a transformation process that takes an input and
// produces an output.
type Transformation struct {
	Version  Version `yaml:"version"`
	Input    Input   `yaml:"input"`
	Output   Output  `yaml:"output"`
	Template string  `yaml:"template"`
}

type Output struct {
	tmpl *template.Template
}

func (o *Output) UnmarshalYAML(node *yaml.Node) error {
	switch node.Kind {
	case yaml.ScalarNode:
		var tmpl template.Template
		err := node.Decode(&tmpl)
		if err != nil {
			return err
		}
		o.tmpl = &tmpl
	default:
		return errors.New("invalid output")
	}
	return nil
}

func (o *Output) Evaluate(data any) (string, error) {
	if o.tmpl == nil {
		return "", nil
	}
	var sb strings.Builder
	if err := o.tmpl.Execute(&sb, data); err != nil {
		return "", err
	}
	return sb.String(), nil
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

type Package struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
}

type Config struct {
	Package         Package          `yaml:"package"`
	Transformations []Transformation `yaml:"transformations"`
}

// FromYAML parses a configuration from a YAML byte slice.
func FromYAML(data []byte) (*Config, error) {
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// FromReader parses a configuration from an io.Reader.
func FromReader(r io.Reader) (*Config, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return FromYAML(data)
}

// FromFile parses a configuration from a file.
func FromFile(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return FromReader(file)
}
