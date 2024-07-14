package archive

import (
	"archive/tar"
	"compress/gzip"
	"io"
)

type Option interface {
	apply(*Archive)
}

type option func(*Archive)

func (o option) apply(a *Archive) {
	o(a)
}

// Visitor is a function that visits each file in the archive.
type Visitor func(name string, r io.Reader) error

// Transform returns an [Option] that sets the transform function for the
// archive.
func Transform(fn func(string) string) Option {
	return option(func(a *Archive) {
		a.transform = fn
	})
}

// Filter returns an [Option] that sets the filter function for the archive.
// If the filter returns false, the file will be skipped.
func Filter(fn func(string) bool) Option {
	return option(func(a *Archive) {
		a.filter = fn
	})
}

// Archive represents a FHIR IG package in an archived format.
type Archive struct {
	reader          io.Reader
	transform       func(string) string
	filter          func(string) bool
	transformReader func(string, io.Reader) io.Reader
}

// New creates a new archive from a reader.
func New(r io.Reader, opts ...Option) *Archive {
	archive := &Archive{
		reader:          r,
		transform:       func(s string) string { return s },
		filter:          func(s string) bool { return true },
		transformReader: func(_ string, r io.Reader) io.Reader { return r },
	}
	for _, opt := range opts {
		opt.apply(archive)
	}
	return archive
}

// VisitFiles visits each file in the archive, calling the provided function
// with the name and file reader.
func (a *Archive) VisitFiles(visitor Visitor) error {
	gzipReader, err := gzip.NewReader(a.reader)
	if err != nil {
		return err
	}
	tarReader := tar.NewReader(gzipReader)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if header.Typeflag != tar.TypeReg {
			continue
		}
		if !a.filter(header.Name) {
			continue
		}

		path := a.transform(header.Name)
		if err := visitor(path, a.transformReader(path, tarReader)); err != nil {
			return err
		}
	}
	return nil
}
