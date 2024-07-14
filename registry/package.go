package registry

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
)

// PackageManifest represents the manifest of a package, in NPM package.json
// format.
type PackageManifest struct {
	Name            string            `json:"name"`
	Version         string            `json:"version"`
	FHIRVersionList []string          `json:"fhir-version-list,omitempty"`
	Type            string            `json:"type"`
	Dependencies    map[string]string `json:"dependencies"`
	License         string            `json:"license"`
	Title           string            `json:"title"`
	Description     string            `json:"description"`
	Author          string            `json:"author"`
	URL             string            `json:"url"`
	ToolsVersion    int               `json:"tools-version"`
	Canonical       string            `json:"canonical"`
	Homepage        string            `json:"homepage"`
}

// Package represents a FHIR package that has been unpacked on-disk.
type Package struct {
	// Path to the location on-disk where the registry has been unpacked to.
	Path string

	// Manifest is the package manifest.
	Manifest *PackageManifest
}

// NewPackage creates a new package from the content at the specified file path.
// The path must contain a package.json content that defines information about
// the FHIR package.
func NewPackage(path string) (*Package, error) {
	var manifest PackageManifest
	metadata := filepath.Join(filepath.FromSlash(path), "package.json")
	file, err := os.Open(metadata)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	if err := json.NewDecoder(file).Decode(&manifest); err != nil {
		return nil, err
	}

	return &Package{
		Path:     path,
		Manifest: &manifest,
	}, nil
}

// Name returns the name of the package.
func (p *Package) Name() string {
	if p.Manifest == nil {
		return ""
	}
	return p.Manifest.Name
}

// Version returns the version of the package.
func (p *Package) Version() string {
	if p.Manifest == nil {
		return ""
	}
	return p.Manifest.Version
}

// Dependencies returns the dependencies for the package.
func (p *Package) Dependencies() map[string]string {
	if p.Manifest == nil {
		return nil
	}
	return p.Manifest.Dependencies
}

// FHIRVersionList returns the list of FHIR versions supported by the package.
func (p *Package) FHIRVersionList() []string {
	if p.Manifest == nil {
		return nil
	}
	return p.Manifest.FHIRVersionList
}

// Canonical returns the canonical URL for the package.
func (p *Package) Canonical() string {
	if p.Manifest == nil {
		return ""
	}
	return p.Manifest.Canonical
}

// Files returns the list of files in the package.
func (p *Package) Files() ([]string, error) {
	var files []string
	err := filepath.Walk(p.Path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		base := filepath.Base(path)
		if base == "package.json" || base == "package.tar.gz" {
			return nil
		}

		files = append(files, path)
		return nil
	})
	slices.Sort(files)
	return files, err
}
