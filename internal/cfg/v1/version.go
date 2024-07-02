package cfg

import (
	"fmt"

	"github.com/friendly-fhir/fhenix/internal/cfg"
	"gopkg.in/yaml.v3"
)

type Version int

func (v *Version) UnmarshalYAML(node *yaml.Node) error {
	var s string
	if err := node.Decode(&s); err != nil {
		return err
	}

	switch s {
	case "1":
		*v = 1
	default:
		return &cfg.FieldError{
			Field: "version",
			Err:   fmt.Errorf("%w: '%v', expected '1'", cfg.ErrInvalidField, s),
		}
	}

	return nil
}
