package registry

import (
	"fmt"
	"strings"
)

// PackageRef is a string that is always in the form of <registry>::<package>@<version>.
// This is intended for easily passing along package references to other
// utilities that need this data in some structured manner.
type PackageRef string

// Registry returns the registry of the package reference.
func (p PackageRef) Registry() string {
	registry, _, _ := p.Parts()
	return registry
}

// Name returns the name of the package reference.
func (p PackageRef) Name() string {
	_, name, _ := p.Parts()
	return name
}

// Version returns the version of the package reference.
func (p PackageRef) Version() string {
	_, _, version := p.Parts()
	return version
}

// Parts returns the registry, name, and version of the package reference.
func (p PackageRef) Parts() (registry, name, version string) {
	parts := strings.Split(string(p), "::")
	rest := parts[0]
	if len(parts) == 1 {
		registry = "default"
		name = parts[0]
	} else if len(parts) == 2 {
		registry = parts[0]
		name = parts[1]
		rest = parts[1]
	}
	parts = strings.Split(rest, "@")
	if len(parts) == 1 {
		version = ""
	} else if len(parts) == 2 {
		name = parts[0]
		version = parts[1]
	}
	return
}

// String returns the string representation of the package reference.
func (p PackageRef) String() string {
	return string(p)
}

var _ fmt.Stringer = (*PackageRef)(nil)

func NewPackageRef(registry, name, version string) PackageRef {
	return PackageRef(fmt.Sprintf("%s::%s@%s", registry, name, version))
}
