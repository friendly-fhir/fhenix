package archive

import (
	"io"
	"os"
	"path/filepath"
)

// Unpackers is a collection of [Unpacker]s that can be used to unpack
// an archive.
type Unpackers []Unpacker

// Unpacker is an interface for unpacking an archive.
func (us Unpackers) Unpack(path string, length int64, r io.Reader) error {
	for _, u := range us {
		if err := u.Unpack(path, length, r); err != nil {
			return err
		}
	}
	return nil
}

var _ Unpacker = (*Unpackers)(nil)

// UnpackFunc is a function that implements the [Unpacker] interface.
type UnpackFunc func(string, int64, io.Reader) error

func (f UnpackFunc) Unpack(path string, length int64, r io.Reader) error {
	return f(path, length, r)
}

var _ Unpacker = (*UnpackFunc)(nil)

// DiskUnpacker is an [Unpacker] that writes files to a directory while
// teeing the contents so that it can be recorded or logged elsewhere.
type DiskUnpacker struct {
	Root string
	Tee  func(name string, r io.Reader) io.Reader
}

// Unpack writes the contents of the reader to a file at the given path.
func (du *DiskUnpacker) Unpack(path string, length int64, r io.Reader) error {
	root := filepath.FromSlash(du.Root)
	path = filepath.FromSlash(path)
	file := filepath.Clean(filepath.Join(root, path))
	if err := os.MkdirAll(filepath.Dir(file), 0755); err != nil {
		return err
	}

	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()

	if du.Tee != nil {
		r = du.Tee(path, r)
	}

	_, err = io.Copy(f, r)
	return err
}

var _ Unpacker = (*DiskUnpacker)(nil)
