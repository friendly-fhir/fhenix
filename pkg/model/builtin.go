package model

import (
	"fmt"
	"regexp"
	"strings"

	fhir "github.com/friendly-fhir/go-fhir/r4/core"
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

func (b *Builtin) FromType(ty *fhir.ElementDefinitionType) error {
	const (
		typeURL    = "http://hl7.org/fhir/StructureDefinition/structuredefinition-fhir-type"
		extURL     = "http://hl7.org/fhir/StructureDefinition/regex"
		fpURPrefix = "http://hl7.org/fhirpath/System."
	)
	if ext := extension(typeURL, ty.Extension); ext != nil {
		b.Name = ext.GetValueURL().GetValue()
		if ext := extension(extURL, ty.Extension); ext != nil {
			var err error
			b.Regex, err = regexp.Compile(ext.GetValueCode().GetValue())
			if err != nil {
				return fmt.Errorf("unable to compile regular expression %s: %w", ext.GetValueCode().GetValue(), err)
			}
		}
	}
	if url, ok := strings.CutPrefix(ty.GetCode().GetValue(), fpURPrefix); ok {
		b.Name = strings.ToLower(url)
		return nil
	}
	return fmt.Errorf("unable to determine builtin type for %s", ty.GetCode().GetValue())
}

func extension(url string, exts []*fhir.Extension) *fhir.Extension {
	for _, ext := range exts {
		if ext.URL == url {
			return ext
		}
	}
	return nil
}
