package npm

import (
	"fmt"
	"strings"
)

// PackageRef is a string that is always in the form of <package>@<version>.
type PackageRef string

// NewPackageRef creates a new package reference.
func NewPackageRef(name, version string) PackageRef {
	return PackageRef(fmt.Sprintf("%s@%s", name, version))
}

// Name returns the name of the package reference.
func (p PackageRef) Name() string {
	name, _ := p.Parts()
	return name
}

// Version returns the version of the package reference.
func (p PackageRef) Version() string {
	_, version := p.Parts()
	return version
}

// Parts returns the registry, name, and version of the package reference.
func (p PackageRef) Parts() (name, version string) {
	parts := strings.SplitN(string(p), "@", 2)
	if len(parts) == 1 {
		name = parts[0]
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
