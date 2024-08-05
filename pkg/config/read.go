package config

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/friendly-fhir/fhenix/pkg/config/internal/opts"
	"gopkg.in/yaml.v3"
)

// FromFile reads a configuration file from the specified path and returns a
// configuration object. This will implicitly set the "root" of the path to be
// the directory containing the file.
func FromFile(file string, options ...Option) (*Config, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	var opts opts.Options
	opts.Apply(options...)

	root := opts.RootDir
	if root == "" {
		root = filepath.Dir(filepath.FromSlash(file))
	}

	opts.RootDir, err = filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	return from(data, &opts)
}

// FromReader reads a configuration file from the specified reader and returns a
// configuration object.
func FromReader(r io.Reader, options ...Option) (*Config, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return FromBytes(data, options...)
}

// FromBytes reads a configuration file from the specified byte slice and returns
// a configuration object.
func FromBytes(bytes []byte, options ...Option) (*Config, error) {
	var opts opts.Options
	opts.Apply(options...)
	return from(bytes, &opts)
}

func from(data []byte, conf *opts.Options) (*Config, error) {
	const (
		supportedVersion = 1
	)
	type version struct {
		Version int `yaml:"version"`
	}
	var v version
	if err := yaml.Unmarshal(data, &v); err != nil {
		return nil, err
	}
	switch v.Version {
	case 1:
		return fromV1(data, conf)
	}
	if v.Version < supportedVersion {
		return nil, fmt.Errorf("%w: version '%d' is obsolete, and no longer suoported", ErrInvalidVersion, v.Version)
	}
	return nil, fmt.Errorf("%w: version '%d' is not supported", ErrInvalidVersion, v.Version)
}
