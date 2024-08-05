/*
Package opts provides default option configuration used internally for
configuration dispatching.
*/
package opts

import (
	"os"
	"path/filepath"
)

// Option is an interface for setting options on configuration objects.
type Option interface {
	set(*Options)
}

type OptionFunc func(*Options)

func (o OptionFunc) set(opts *Options) {
	o(opts)
}

var _ Option = (*OptionFunc)(nil)

// Options is the configuration for options provided to the reader for
// configurations.
type Options struct {
	// OutputDir is the directory where generated output will be written.
	OutputDir string

	// RootDir is the root directory which all configuration paths will be
	// considered relative to.
	RootDir string
}

func (o *Options) Apply(opts ...Option) {
	for _, opt := range opts {
		opt.set(o)
	}
}

// RootPath returns the absolute path of the specified path relative to the
// root directory. If the path is already absolute, it will be returned as is.
// If the root directory is not set, the current working directory will be used.
func (o *Options) RootPath(path string) (string, error) {
	return o.path(o.RootDir, path)
}

// OutputPath returns the absolute path of the specified path relative to the
// output directory. If the path is already absolute, it will be returned as is.
// If the output directory is not set, the current working directory will be used.
func (o *Options) OutputPath(path string) (string, error) {
	return o.path(o.OutputDir, path)
}

func (c *Options) path(root, path string) (string, error) {
	if filepath.IsAbs(path) {
		return path, nil
	}
	if root == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return "", err
		}
		return filepath.Join(cwd, filepath.FromSlash(path)), nil
	}
	return filepath.Clean(filepath.Join(filepath.FromSlash(root), filepath.FromSlash(path))), nil
}
