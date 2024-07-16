package dedent_test

import (
	"testing"

	"github.com/friendly-fhir/fhenix/internal/dedent"
	"github.com/google/go-cmp/cmp"
)

func TestString(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name: "Words with same indentation",
			input: `
				Hello,
				World!
			`,
			want: "Hello,\nWorld!",
		},
		{
			name: "Words with different indentation",
			input: `
				Hello,
					World!
			`,
			want: "Hello,\n\tWorld!",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := dedent.String(test.input)

			if want := test.want; got != want {
				t.Errorf("String(%s): got %q, want %q", test.name, want, got)
			}
		})
	}
}

func TestLines(t *testing.T) {
	tests := []struct {
		name  string
		input []string
		want  []string
	}{
		{
			name: "Words with same indentation",
			input: []string{
				"    Hello,",
				"    World!",
			},
			want: []string{
				"Hello,",
				"World!",
			},
		},
		{
			name: "Words with different indentation",
			input: []string{
				"    Hello,",
				"        World!",
			},
			want: []string{
				"Hello,",
				"    World!",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := dedent.Lines(test.input...)

			if want := test.want; !cmp.Equal(got, want) {
				t.Errorf("Lines(%s): got %v, want %v", test.name, want, got)
			}
		})
	}
}
