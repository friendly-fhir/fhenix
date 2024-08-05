package cfg_test

import (
	"testing"

	rootcfg "github.com/friendly-fhir/fhenix/pkg/config/internal/cfg"
	"github.com/friendly-fhir/fhenix/pkg/config/internal/cfg/v1"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"gopkg.in/yaml.v3"
)

func TestVersion(t *testing.T) {
	testCases := []struct {
		name    string
		input   string
		want    cfg.Version
		wantErr error
	}{
		{
			name:    "valid version",
			input:   "1",
			want:    cfg.Version(1),
			wantErr: nil,
		}, {
			name:    "invalid version",
			input:   "2",
			want:    cfg.Version(0),
			wantErr: rootcfg.ErrInvalidField,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var got cfg.Version
			err := yaml.Unmarshal([]byte(tc.input), &got)

			if want := tc.want; !cmp.Equal(got, want) {
				t.Errorf("Version.UnmarshalYAML() = %v, want %v", got, want)
			}
			if got, want := err, tc.wantErr; !cmp.Equal(got, want, cmpopts.EquateErrors()) {
				t.Errorf("Version.UnmarshalYAML(...) = %v, want %v", got, want)
			}
		})
	}
}
