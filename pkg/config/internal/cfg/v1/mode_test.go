package cfg_test

import (
	"testing"

	rootcfg "github.com/friendly-fhir/fhenix/pkg/config/internal/cfg"
	"github.com/friendly-fhir/fhenix/pkg/config/internal/cfg/v1"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"gopkg.in/yaml.v3"
)

func TestMode(t *testing.T) {
	testCases := []struct {
		name    string
		input   string
		want    cfg.Mode
		wantErr error
	}{
		{
			name:    "'text' Mode",
			input:   "text",
			want:    cfg.Mode("text"),
			wantErr: nil,
		}, {
			name:    "'text' Mode",
			input:   "html",
			want:    cfg.Mode("html"),
			wantErr: nil,
		}, {
			name:    "invalid Mode",
			input:   "invalid",
			want:    "",
			wantErr: rootcfg.ErrInvalidField,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var got cfg.Mode
			err := yaml.Unmarshal([]byte(tc.input), &got)

			if want := tc.want; !cmp.Equal(got, want) {
				t.Errorf("Mode.UnmarshalYAML() = %v, want %v", got, want)
			}
			if got, want := err, tc.wantErr; !cmp.Equal(got, want, cmpopts.EquateErrors()) {
				t.Errorf("Mode.UnmarshalYAML(...) = %v, want %v", got, want)
			}
		})
	}
}
