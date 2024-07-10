package model

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/friendly-fhir/fhenix/model/raw"
)

// Builtin is a helper type that is used to represent language-level builtin types.
type Builtin struct {
	// Name is the name of the builtin type.
	Name string

	// Regex is a regular expression that the value must match.
	Regex *regexp.Regexp
}

func (b *Builtin) ValidateString(s string) bool {
	if b.Regex == nil {
		return true
	}
	return b.Regex.MatchString(s)
}

func (b *Builtin) FromType(ty *raw.Type) error {
	const (
		typeURL    = "http://hl7.org/fhir/StructureDefinition/structuredefinition-fhir-type"
		extURL     = "http://hl7.org/fhir/StructureDefinition/regex"
		fpURPrefix = "http://hl7.org/fhirpath/System."
	)
	if ext := ty.Extension(typeURL); ext != nil {
		b.Name = ext.ValueURL
		if ext := ty.Extension(extURL); ext != nil {
			var err error
			b.Regex, err = regexp.Compile(ext.ValueCode)
			if err != nil {
				return fmt.Errorf("unable to compile regular expression %s: %w", ext.ValueCode, err)
			}
		}
	}
	if url, ok := strings.CutPrefix(ty.Code, fpURPrefix); ok {
		b.Name = strings.ToLower(url)
		return nil
	}
	return fmt.Errorf("unable to determine builtin type for %s", ty.Code)
}
