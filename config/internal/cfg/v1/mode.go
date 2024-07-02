package cfg

import (
	"fmt"

	"github.com/friendly-fhir/fhenix/config/internal/cfg"
	"gopkg.in/yaml.v3"
)

type Mode string

const (
	ModeText Mode = "text"
	ModeHTML Mode = "html"
)

// UnmarshalYAML unmarshals a YAML node into a Mode.
func (m *Mode) UnmarshalYAML(node *yaml.Node) error {
	var s string
	if err := node.Decode(&s); err != nil {
		return err
	}

	switch s {
	case "text", "html":
		*m = Mode(s)
	default:
		return &cfg.FieldError{
			Field: "mode",
			Err:   fmt.Errorf("%w: '%v', expected 'text' or 'html'", cfg.ErrInvalidField, s),
		}
	}

	return nil
}
