package npm

import (
	"archive/tar"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"time"
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

// ReadManifest reads a package manifest from the given reader.
func ReadManifest(r io.Reader) (*PackageManifest, error) {
	var manifest PackageManifest
	if err := json.NewDecoder(r).Decode(&manifest); err != nil {
		return nil, err
	}
	return &manifest, nil
}

// PackageVisitor is an interface for visiting the contents of a package
// from a npm tarball.
type PackageVisitor interface {
	// VisitManifest is called when the package manifest is read.
	VisitManifest(manifest *PackageManifest) error

	// VisitFile is called for each file in the package, including the manifest.
	VisitFile(path string, info fs.FileInfo, r io.Reader) error

	// visitor is an unexported func to ensure that only things embedding
	// [BasePackageVisitor] can be valid package visitors.
	visitor()
}

// BasePackageVisitor is a base implementation of PackageVisitor that does
// nothing.
type BasePackageVisitor struct{}

// VisitManifest is called when the package manifest is read.
func (*BasePackageVisitor) VisitManifest(_ *PackageManifest) error {
	return nil
}

// VisitFile is called for each file in the package, including the manifest.
func (*BasePackageVisitor) VisitFile(_ string, _ fs.FileInfo, _ io.Reader) error {
	return nil
}

func (*BasePackageVisitor) visitor() {}

// VisitPackageStream reads an npm package from the given reader and calls the
// appropriate methods on the given visitor. This is assumed to be in 'tar'
// format without any compression.
//
// If visitor is nil, the reader is still consumed but no side effects are
// produced.
//
// This utility works directly with readers returned from [Client.Fetch]
func VisitPackageStream(r io.Reader, visitor PackageVisitor) error {
	if r == nil {
		return nil
	}
	if visitor == nil {
		visitor = &BasePackageVisitor{}
	}

	tarReader := tar.NewReader(r)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		// Ignore directories since the only data it provides will be directory
		// permissions, and FHIR IGs are typically flat and have no meaningful
		// directory permissions; we can use 755.
		if header.FileInfo().IsDir() {
			continue
		}
		path := filepath.Clean(header.Name)

		var reader io.Reader = tarReader
		if filepath.Base(path) == "package.json" {
			var buffer bytes.Buffer
			manifest, err := ReadManifest(io.TeeReader(tarReader, &buffer))
			if err != nil {
				return err
			}
			if err := visitor.VisitManifest(manifest); err != nil {
				return err
			}
			reader = &buffer
		}

		if err := visitor.VisitFile(path, header.FileInfo(), reader); err != nil {
			return err
		}
	}
	return nil
}

// UnpackPackageVisitor is a [PackageVisitor] that unpacks the package to a
// directory on disk. If the destination is empty, the package is not unpacked.
type UnpackPackageVisitor struct {
	// Destination is the directory to unpack the package to.
	Destination string

	BasePackageVisitor
}

// VisitFile is called for each file in the package, including the manifest.
func (v *UnpackPackageVisitor) VisitFile(path string, info fs.FileInfo, r io.Reader) error {
	if v.Destination == "" {
		return nil
	}

	dest := filepath.Join(v.Destination, path)
	dir := filepath.Dir(dest)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("unpacking %q: %w", dir, err)
	}

	file, err := os.OpenFile(dest, os.O_CREATE|os.O_WRONLY, info.Mode())
	if err != nil {
		return fmt.Errorf("unpacking %q: %w", dest, err)
	}
	if _, err := io.Copy(file, r); err != nil {
		return fmt.Errorf("unpacking %q: %w", dest, err)
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf("unpacking %q: %w", dest, err)
	}

	// It's not critical if this fails, so don't report the error
	_ = os.Chtimes(dest, time.Time{}, info.ModTime())

	return nil
}
