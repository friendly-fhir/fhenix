/*
Package structuredefinition provides raw JSON bindings for StructureDefinition
resources so that they can be easily marshalled and parsed.

This is not a complete implementation; it will deliberately ignore many fields,
since all it really cares for are the "Snapshot" and "Differential" fields along
with the defining attributes (name, comments, etc). This baseline is what enables
generating for FHIR resources.
*/
package raw

type BackboneElement struct {
	Element []*ElementDefinition `json:"element"`
}

type Extension struct {
	URL          string `json:"url"`
	ValueURL     string `json:"valueURL"`
	ValueCode    string `json:"valueCode,omitempty"`
	ValueInteger int    `json:"valueInteger,omitempty"`
}

type Extensions []Extension

func (e Extensions) GetByURL(url string) *Extension {
	for _, ext := range e {
		if ext.URL == url {
			return &ext
		}
	}
	return nil
}

type Type struct {
	Code       string      `json:"code"`
	Extensions []Extension `json:"extension"`
}

func (t *Type) Extension(url string) *Extension {
	return Extensions(t.Extensions).GetByURL(url)
}

type ElementDefinition struct {
	ID   string `json:"id"`
	Path string `json:"path"`
	Min  int    `json:"min"`
	Max  string `json:"max"`
	Base *struct {
		Min int    `json:"min"`
		Max string `json:"max"`
	} `json:"base,omitempty"`
	Short      string   `json:"short"`
	Definition string   `json:"definition"`
	Comment    string   `json:"comment"`
	Types      []Type   `json:"type"`
	Binding    *Binding `json:"binding,omitempty"`
}

type Binding struct {
	Strength string `json:"strength"`
	ValueSet string `json:"valueSet"`
}

// StructureDefinition represents a FHIR Profile StructureDefinition
type StructureDefinition struct {
	ResourceType   string          `json:"resourceType"`
	BaseDefinition string          `json:"baseDefinition"`
	Status         string          `json:"status"`
	URL            string          `json:"url"`
	Name           string          `json:"name"`
	Type           string          `json:"type"`
	Kind           string          `json:"kind"`
	Abstract       bool            `json:"abstract"`
	Short          string          `json:"short"`
	Description    string          `json:"description"`
	Comment        string          `json:"comment"`
	FHIRVersion    string          `json:"fhirVersion"`
	Derivation     string          `json:"derivation"`
	Snapshot       BackboneElement `json:"snapshot"`
	Differential   BackboneElement `json:"differential"`
}

func ReadStructureDefinition(path string) (*StructureDefinition, error) {
	return readJSON[StructureDefinition](path)
}
