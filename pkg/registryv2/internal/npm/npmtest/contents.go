package npmtest

import (
	"archive/tar"
	"embed"
	"errors"
	"io"
	"io/fs"
	"testing"
)

//go:embed testdata/package
var testdataPackage embed.FS

// SimplePackageFS returns an fs.FS with a flat package structure and no
// dependencies.
func SimplePackageFS() fs.FS {
	return testdataPackage
}

// SimplePackageTarball returns a reader to a tarball with a flat package
// structure and no dependencies.
func SimplePackageTarball() io.ReadCloser {
	return NewTarballFromFS(testdataPackage)
}

// NewTarballFromFS returns a reader to a tarball with the contents of an
// [fs.FS]. The caller is responsible for closing the returned reader.
func NewTarballFromFS(f fs.FS) io.ReadCloser {
	pr, pw := io.Pipe()
	go func() {
		writer := tar.NewWriter(pw)
		if err := writer.AddFS(f); err != nil {
			_ = pw.CloseWithError(err)
			return
		}

		err := writer.Close()
		_ = pw.CloseWithError(err)
	}()
	return pr
}

// TarContents reads the contents of a tarball and returns a map of file paths
// to their contents.
func TarContents(t *testing.T, r io.ReadCloser) map[string]string {
	t.Helper()
	if r == nil {
		return nil
	}
	defer r.Close()
	files := make(map[string]string)

	tr := tar.NewReader(r)
	for {
		hdr, err := tr.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			t.Fatalf("unable to read tar entry: %v", err)
		}
		bytes, err := io.ReadAll(tr)
		if err != nil {
			t.Fatalf("unable to read tar entry contents: %v", err)
		}
		files[hdr.Name] = string(bytes)
	}

	return files
}

// FSContents reads the contents of an fs.FS and returns a map of file paths to
// their contents.
func FSContents(t *testing.T, f fs.FS) map[string]string {
	t.Helper()
	files := make(map[string]string)

	err := fs.WalkDir(f, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		file, err := fs.ReadFile(f, path)
		if err != nil {
			return err
		}
		files[path] = string(file)
		return nil
	})
	if err != nil {
		t.Fatalf("unable to read fs contents: %v", err)
	}

	return files
}
