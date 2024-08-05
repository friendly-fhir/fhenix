package template_test

import (
	"testing"

	"github.com/friendly-fhir/fhenix/pkg/transform/internal/template"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestFromString(t *testing.T) {
	testCases := []struct {
		name    string
		in      string
		want    template.Engine
		wantErr error
	}{
		{
			name:    "invalid template",
			in:      "invalid",
			want:    nil,
			wantErr: cmpopts.AnyError,
		}, {
			name: "text template",
			in:   "text",
			want: template.Text(),
		}, {
			name: "html template",
			in:   "html",
			want: template.HTML(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := template.FromString(tc.in)

			if got, want := err, tc.wantErr; !cmp.Equal(got, want, cmpopts.EquateErrors()) {
				t.Fatalf("FromString() = %v, want %v", got, want)
			}
			if got, want := got, tc.want; !cmp.Equal(got, want) {
				t.Errorf("FromString() = %v, want %v", got, want)
			}
		})
	}
}
