package registrytest

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"io/fs"
)

// TarballBytes returns a tarball of the given filesystem as a byte slice.
func TarballBytes(fs fs.FS) []byte {
	var buf bytes.Buffer
	w := tar.NewWriter(&buf)
	if err := w.AddFS(fs); err != nil {
		panic(err)
	}
	if err := w.Close(); err != nil {
		panic(err)
	}
	return buf.Bytes()
}

// GzipTarballBytes returns a gzipped tarball of the given filesystem as a byte slice.
func GzipTarballBytes(fs fs.FS) []byte {
	var buf bytes.Buffer
	w := tar.NewWriter(gzip.NewWriter(&buf))
	if err := w.AddFS(fs); err != nil {
		panic(err)
	}
	if err := w.Close(); err != nil {
		panic(err)
	}
	return buf.Bytes()
}
