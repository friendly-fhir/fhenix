package config

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/friendly-fhir/fhenix/internal/template"
	"gopkg.in/yaml.v3"
)

type Version int

// Transformation outlines a transformation process that takes an input and
// produces an output.
type Transformation struct {
	Input    Input    `yaml:"input"`
	Output   Output   `yaml:"output"`
	Partials Partials `yaml:"partials"`
}

// Default is a configuration node that specifies default values to use for
// output and templates, if not specified in a transformation.
// This just helps to reduce the boilerplate when several transformations use
// the same set of templates.
type Default struct {
	Dist     string   `yaml:"dist-dir"`
	Output   Output   `yaml:"output"`
	Partials Partials `yaml:"partials"`
}

type Partials map[string]string

func (t *Partials) UnmarshalYAML(node *yaml.Node) error {
	switch node.Kind {
	case yaml.MappingNode:
		*t = make(Partials)
		for i := 0; i < len(node.Content); i += 2 {
			key := node.Content[i].Value
			var tmpl string
			if err := node.Content[i+1].Decode(&tmpl); err != nil {
				return err
			}
			(*t)[key] = tmpl
		}
	case yaml.ScalarNode:
		var tmpl string
		if err := node.Decode(&tmpl); err != nil {
			return err
		}
		*t = Partials{"default": tmpl}
	default:
		return errors.New("config: invalid 'template' definition, expected mapping or string")
	}
	return nil
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
	Default         Default          `yaml:"default"`
	Package         Package          `yaml:"package"`
	Transformations []Transformation `yaml:"transformations"`
	BasePath        string           `yaml:"-"`
}

// FromYAML parses a configuration from a YAML byte slice.
func FromYAML(data []byte) (*Config, error) {
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	for _, transform := range cfg.Transformations {
		if transform.Partials == nil {
			transform.Partials = make(Partials)
		}
		for key, value := range cfg.Default.Partials {
			if _, ok := transform.Partials[key]; !ok {
				transform.Partials[key] = value
			}
		}
		if transform.Output.tmpl == nil {
			transform.Output = cfg.Default.Output
		}
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
	cfg, err := FromReader(file)
	if err != nil {
		return nil, err
	}
	cfg.BasePath = filepath.Dir(path)
	for _, t := range cfg.Transformations {
		for k, v := range t.Partials {
			if !filepath.IsAbs(v) {
				t.Partials[k] = filepath.Join(cfg.BasePath, v)
			}
		}
	}
	if !filepath.IsAbs(cfg.Default.Dist) {
		cfg.Default.Dist = filepath.Join(cfg.BasePath, cfg.Default.Dist)
	}
	return cfg, nil
}
