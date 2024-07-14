package registry_test

import (
	"testing"

	"github.com/friendly-fhir/fhenix/registry"
)

func TestPackageRef_Registry(t *testing.T) {
	testCases := []struct {
		name  string
		input registry.PackageRef
		want  string
	}{
		{
			name:  "good reference definition",
			input: "default::test.package@1.0.0",
			want:  "default",
		}, {
			name:  "no registry name",
			input: "test.package@1.0.0",
			want:  "default",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got, want := tc.input.Registry(), tc.want; got != want {
				t.Errorf("PackageRef.Registry() = %q, want %q", got, want)
			}
		})
	}
}

func TestPackageRef_Name(t *testing.T) {
	testCases := []struct {
		name  string
		input registry.PackageRef
		want  string
	}{
		{
			name:  "good reference definition",
			input: "default::test.package@1.0.0",
			want:  "test.package",
		}, {
			name:  "no package name",
			input: "default::@1.0.0",
			want:  "",
		}, {
			name:  "no version",
			input: "default::test.package",
			want:  "test.package",
		}, {
			name:  "bad package definition",
			input: "package",
			want:  "package",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got, want := tc.input.Name(), tc.want; got != want {
				t.Errorf("PackageRef.Name() = %q, want %q", got, want)
			}
		})
	}
}

func TestPackageRef_Version(t *testing.T) {
	testCases := []struct {
		name  string
		input registry.PackageRef
		want  string
	}{
		{
			name:  "good reference definition",
			input: "default::test.package@1.0.0",
			want:  "1.0.0",
		}, {
			name:  "no version",
			input: "default::test.package",
			want:  "",
		}, {
			name:  "bad package definition",
			input: "package",
			want:  "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got, want := tc.input.Version(), tc.want; got != want {
				t.Errorf("PackageRef.Version() = %q, want %q", got, want)
			}
		})
	}
}
