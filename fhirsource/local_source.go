package fhirsource

import (
	"context"
	"os"
	"path/filepath"
)

type localSource string

func (ls localSource) Definitions(ctx context.Context) ([]string, error) {
	path := filepath.FromSlash(string(ls))
	return ls.walkEntries(nil, path)
}

func (ls localSource) walkEntries(result []string, path string) ([]string, error) {
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
