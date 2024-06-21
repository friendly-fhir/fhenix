package model

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/friendly-fhir/fhenix/internal/fhirig"
	"github.com/friendly-fhir/fhenix/internal/model/raw"
)

type FHIRDefinition[T any] struct {
	Definition *T
	Package    *fhirig.Package
	File       string
}

// FHIRCache is a helper type that assists in caching the raw parsed FHIR
// definition resources.
type FHIRCache struct {
	base string
	sds  map[string]*FHIRDefinition[raw.StructureDefinition]
	css  map[string]*FHIRDefinition[raw.CodeSystem]
}

// NewFHIRCache creates a new definition cache instance.
func NewFHIRCache(base string) *FHIRCache {
	if !strings.HasSuffix(base, "/") {
		base += "/"
	}
	return &FHIRCache{
		base: base,
		sds:  map[string]*FHIRDefinition[raw.StructureDefinition]{},
		css:  map[string]*FHIRDefinition[raw.CodeSystem]{},
	}
}

// DefaultFHIRCache creates a new definition cache instance with the
// default root URL.
func DefaultFHIRCache() *FHIRCache {
	return NewFHIRCache("http://hl7.org/fhir")
}

func (dc *FHIRCache) Base() string {
	return strings.TrimSuffix(dc.base, "/")
}

// StructureDefinitions returns all the structure definitions in the cache sorted
// by their URL.
func (dc *FHIRCache) StructureDefinitions() []*FHIRDefinition[raw.StructureDefinition] {
	result := make([]*FHIRDefinition[raw.StructureDefinition], 0, len(dc.sds))
	for _, sd := range dc.sds {
		result = append(result, sd)
	}
	slices.SortFunc(result, func(lhs, rhs *FHIRDefinition[raw.StructureDefinition]) int {
		return strings.Compare(dc.sdRoot(lhs.Definition.URL), dc.sdRoot(rhs.Definition.URL))
	})
	return result
}

func (dc *FHIRCache) sdRoot(url string) string {
	if strings.Contains(url, "/") {
		return url
	}
	return dc.base + "StructureDefinition/" + url
}

func (dc *FHIRCache) csRoot(url string) string {
	if strings.Contains(url, "/") {
		return url
	}
	return dc.base + "CodeSystem/" + url
}

// GetStructureDefinition returns the structure definition with the given URL.
func (dc *FHIRCache) GetStructureDefinition(url string) *FHIRDefinition[raw.StructureDefinition] {
	return dc.sds[dc.sdRoot(url)]
}

func (dc *FHIRCache) LookupStructureDefinition(url string) (sd *FHIRDefinition[raw.StructureDefinition], ok bool) {
	sd, ok = dc.sds[dc.sdRoot(url)]
	return
}

// AddStructureDefinition adds a structure definition to the cache, overwriting
// any with the same URL if it already existed.
func (dc *FHIRCache) AddStructureDefinition(pkg *fhirig.Package, file string, sd *raw.StructureDefinition) {
	dc.sds[dc.sdRoot(sd.URL)] = &FHIRDefinition[raw.StructureDefinition]{
		Definition: sd,
		Package:    pkg,
		File:       file,
	}
}

// ParseStructureDefinition parses a structure definition from the given reader
// and adds it to the cache.
func (dc *FHIRCache) ParseStructureDefinition(pkg *fhirig.Package, file string, r io.Reader) error {
	var sd raw.StructureDefinition
	if err := json.NewDecoder(r).Decode(&sd); err != nil {
		return err
	}
	dc.AddStructureDefinition(pkg, file, &sd)
	return nil
}

// ParseStructureDefinitionFromFile parses a structure definition from the given
// file and adds it to the cache.
func (dc *FHIRCache) ParseStructureDefinitionFromFile(pkg *fhirig.Package, path string) error {
	r, err := os.Open(path)
	if err != nil {
		return err
	}
	defer r.Close()
	return dc.ParseStructureDefinition(pkg, filepath.Base(path), r)
}

// GetCodeSystem returns the code system with the given URL.
func (dc *FHIRCache) GetCodeSystem(url string) *FHIRDefinition[raw.CodeSystem] {
	return dc.css[dc.csRoot(url)]
}

// LookupCodeSystem looks up a code system by URL in the cache. If the code system
// does not exist, ok will be false.
func (dc *FHIRCache) LookupCodeSystem(url string) (cs *FHIRDefinition[raw.CodeSystem], ok bool) {
	cs, ok = dc.css[dc.csRoot(url)]
	return
}

// CodeSystems returns all the code systems in the cache sorted by their URL.
func (dc *FHIRCache) CodeSystems() []*FHIRDefinition[raw.CodeSystem] {
	result := make([]*FHIRDefinition[raw.CodeSystem], 0, len(dc.css))
	for _, cs := range dc.css {
		result = append(result, cs)
	}
	slices.SortFunc(result, func(lhs, rhs *FHIRDefinition[raw.CodeSystem]) int {
		return strings.Compare(dc.csRoot(lhs.Definition.URL), dc.csRoot(rhs.Definition.URL))
	})
	return result
}

// AddCodeSystem adds a code system to the cache, overwriting any with the same
// URL if it already existed.
func (dc *FHIRCache) AddCodeSystem(pkg *fhirig.Package, file string, cs *raw.CodeSystem) {
	dc.css[cs.URL] = &FHIRDefinition[raw.CodeSystem]{
		Definition: cs,
		Package:    pkg,
		File:       file,
	}
}

// ParseCodeSystem parses a code system from the given reader and adds it to the
// cache.
func (dc *FHIRCache) ParseCodeSystem(pkg *fhirig.Package, file string, r io.Reader) error {
	var cs raw.CodeSystem
	if err := json.NewDecoder(r).Decode(&cs); err != nil {
		return err
	}
	dc.AddCodeSystem(pkg, file, &cs)
	return nil
}

// ParseCodeSystemFromFile parses a code system from the given file and adds it
// to the cache.
func (dc *FHIRCache) ParseCodeSystemFromFile(pkg *fhirig.Package, path string) error {
	r, err := os.Open(path)
	if err != nil {
		return err
	}
	defer r.Close()
	return dc.ParseCodeSystem(pkg, filepath.Base(path), r)
}
