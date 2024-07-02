package cfg

import (
	"errors"
	"regexp"
	"strings"

	"github.com/friendly-fhir/fhenix/internal/cfg"
	"gopkg.in/yaml.v3"
)

// Input is a configuration node for specifying the input package.
type Input struct {
	// Package is the name of the input package to parse (mandatory).
	Package *InputPackage `yaml:"package"`

	// IncludeDependencies is an optional field that will include dependencies
	// of the package when passing contents to the generator. The default for this
	// is false.
	IncludeDependencies bool `yaml:"include-dependencies"`
}

func (i *Input) UnmarshalYAML(node *yaml.Node) error {
	type input Input
	var out input
	if err := node.Decode(&out); err != nil {
		return err
	}

	if out.Package == nil {
		return &cfg.FieldError{Field: "input.package", Err: cfg.ErrMissingField}
	}
	*i = Input(out)

	return nil
}

type InputPackage struct {
	// Name is the name of the package (mandatory).
	Name string `yaml:"name"`

	// Version is a version string for the package version (mandatory).
	Version string `yaml:"version"`

	// Path is an optional path to specify to where the package is located.
	// If specified, this will override the package being fetched from the
	// package registry.
	Path string `yaml:"path"`
}

var (
	packageNameRegex = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_\-\.]+$`)
	versionRegex     = regexp.MustCompile(`^[0-9]+\.[0-9]+\.[0-9]+$`)
)

func (ip *InputPackage) UnmarshalYAML(node *yaml.Node) error {
	type input InputPackage
	var out input
	if err := node.Decode(&out); err != nil {
		return err
	}

	var errs []error
	if strings.TrimSpace(out.Name) == "" {
		errs = append(errs, &cfg.FieldError{Field: "package.name", Err: cfg.ErrMissingField})
	}

	if strings.TrimSpace(out.Version) == "" {
		errs = append(errs, &cfg.FieldError{Field: "package.version", Err: cfg.ErrMissingField})
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	if !packageNameRegex.MatchString(out.Name) {
		errs = append(errs, &cfg.FieldError{Field: "package.name", Err: cfg.ErrInvalidField})
	}

	if !versionRegex.MatchString(out.Version) {
		errs = append(errs, &cfg.FieldError{Field: "package.version", Err: cfg.ErrInvalidField})
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	*ip = InputPackage(out)
	return nil
}

var _ yaml.Unmarshaler = (*InputPackage)(nil)
