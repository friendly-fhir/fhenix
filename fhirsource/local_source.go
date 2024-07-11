package fhirsource

import (
	"context"
	"os"
	"path/filepath"

	"github.com/friendly-fhir/fhenix/internal/fhirig"
)

// NewLocalSource constructs a new local Source from the specified path.
func NewLocalSource(pkg *fhirig.Package, path string) Source {
	return &localSource{
		path: path,
		pkg:  pkg,
	}
}

type localSource struct {
	path string
	pkg  *fhirig.Package
}

func (ls *localSource) Bundles(ctx context.Context) ([]*Bundle, error) {
	path := filepath.FromSlash(string(ls.path))
	entries, err := ls.walkEntries(nil, path)
	if err != nil {
		return nil, err
	}
	return []*Bundle{
		{
			Package: ls.pkg,
			Files:   entries,
		},
	}, nil
}

func (ls *localSource) walkEntries(result []string, path string) ([]string, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			result, err = ls.walkEntries(result, filepath.Join(path, entry.Name()))
			if err != nil {
				return nil, err
			}
			continue
		}
		result = append(result, filepath.Join(path, entry.Name()))
	}
	return result, nil
}

var _ Source = (*localSource)(nil)
