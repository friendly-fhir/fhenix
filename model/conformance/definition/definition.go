/*
Package definition provides the raw definitions for FHIR data types.
*/
package definition

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	fhir "github.com/friendly-fhir/go-fhir/r4/core"
	"github.com/friendly-fhir/go-fhir/r4/core/resources/codesystem"
	"github.com/friendly-fhir/go-fhir/r4/core/resources/conceptmap"
	"github.com/friendly-fhir/go-fhir/r4/core/resources/structuredefinition"
	"github.com/friendly-fhir/go-fhir/r4/core/resources/valueset"
)

type StructureDefinition = structuredefinition.StructureDefinition
type ValueSets = valueset.ValueSet
type CodeSystem = codesystem.CodeSystem
type ConceptMap = conceptmap.ConceptMap

// Canonical is an interface that represents a FHIR definition that has a
// canonical URL.
type Canonical interface {
	GetURL() *fhir.URI
}

var (
	_ Canonical = (*StructureDefinition)(nil)
	_ Canonical = (*ValueSets)(nil)
	_ Canonical = (*CodeSystem)(nil)
	_ Canonical = (*ConceptMap)(nil)
)

// FromFile reads a canonical definition from a file path.
func FromFile(path string) (Canonical, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return FromReader(file)
}

// FromReader returns a canonical definition from a reader.
func FromReader(reader io.Reader) (Canonical, error) {
	bytes, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return FromJSON(bytes)
}

// FromJSON returns a canonical definition from a JSON byte definition.
func FromJSON(data []byte) (Canonical, error) {
	var resourceType struct {
		ResourceType string `json:"resourceType"`
	}
	if err := json.Unmarshal(data, &resourceType); err != nil {
		return nil, err
	}
	var result Canonical
	switch resourceType.ResourceType {
	case "StructureDefinition":
		result = &StructureDefinition{}
	case "ValueSet":
		result = &ValueSets{}
	case "CodeSystem":
		result = &CodeSystem{}
	case "ConceptMap":
		result = &ConceptMap{}
	default:
		return nil, fmt.Errorf("unknown resource type %q", resourceType.ResourceType)
	}
	err := json.Unmarshal(data, result)
	return result, err
}
