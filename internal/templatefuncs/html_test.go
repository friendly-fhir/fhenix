package templatefuncs_test

import (
	"testing"

	"github.com/friendly-fhir/fhenix/internal/templatefuncs"
)

func TestHTML_Escape(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		s    string
		want string
	}{
		{"empty", "", ""},
		{"simple", "hello", "hello"},
		{"escape", "hello\nworld", "hello\nworld"},
		{"quote", `"hello"`, "&#34;hello&#34;"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &templatefuncs.HTMLModule{}
			if got := m.Escape(tt.s); got != tt.want {
				t.Errorf("HTMLModule.Escape() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHTML_Unescape(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		s    string
		want string
	}{
		{"empty", "", ""},
		{"simple", "hello", "hello"},
		{"escape", "hello\nworld", "hello\nworld"},
		{"quote", "&#34;hello&#34;", `"hello"`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &templatefuncs.HTMLModule{}
			if got := m.Unescape(tt.s); got != tt.want {
				t.Errorf("HTMLModule.Unescape() = %v, want %v", got, tt.want)
			}
		})
	}
}
