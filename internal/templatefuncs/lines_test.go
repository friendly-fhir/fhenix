package templatefuncs_test

import (
	"testing"

	"github.com/friendly-fhir/fhenix/internal/templatefuncs"
)

func TestLines_Prefix(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		prefix string
		text   string
		want   string
	}{
		{"empty", "", "", ""},
		{"no prefix", "prefix", "text", "prefixtext"},
		{"single line", "prefix", "text", "prefixtext"},
		{"multi line", "prefix", "line1\nline2", "prefixline1\nprefixline2"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &templatefuncs.LineModule{}

			if got := m.Prefix(tt.prefix, tt.text); got != tt.want {
				t.Errorf("LineModule.Prefix() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLines_Suffix(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		suffix string
		text   string
		want   string
	}{
		{"empty", "", "", ""},
		{"no suffix", "suffix", "text", "textsuffix"},
		{"single line", "suffix", "text", "textsuffix"},
		{"multi line", "suffix", "line1\nline2", "line1suffix\nline2suffix"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &templatefuncs.LineModule{}

			if got := m.Suffix(tt.suffix, tt.text); got != tt.want {
				t.Errorf("LineModule.Suffix() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLines_TrimSpace(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		text string
		want string
	}{
		{"empty", "", ""},
		{"no space", "text", "text"},
		{"single line", " text ", "text"},
		{"multi line", " line1 \n line2 ", "line1\nline2"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &templatefuncs.LineModule{}

			if got := m.TrimSpace(tt.text); got != tt.want {
				t.Errorf("LineModule.TrimSpace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLine_Indent(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		indent int
		text   string
		want   string
	}{
		{"empty", 0, "", ""},
		{"no indent", 2, "text", "  text"},
		{"single line", 2, "text", "  text"},
		{"multi line", 2, "line1\nline2", "  line1\n  line2"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &templatefuncs.LineModule{}

			if got := m.Indent(tt.indent, tt.text); got != tt.want {
				t.Errorf("LineModule.Indent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLine_TabIndent(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		indent int
		text   string
		want   string
	}{
		{"empty", 0, "", ""},
		{"no indent", 2, "text", "\t\ttext"},
		{"single line", 2, "text", "\t\ttext"},
		{"multi line", 2, "line1\nline2", "\t\tline1\n\t\tline2"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &templatefuncs.LineModule{}

			if got := m.TabIndent(tt.indent, tt.text); got != tt.want {
				t.Errorf("LineModule.TabIndent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLine_Trim(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		cutset string
		text   string
		want   string
	}{
		{"empty", "", "", ""},
		{"no trim", " ", "text", "text"},
		{"single line", " ", " text ", "text"},
		{"multi line", " ", " line1 \n line2 ", "line1\nline2"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &templatefuncs.LineModule{}

			if got := m.Trim(tt.cutset, tt.text); got != tt.want {
				t.Errorf("LineModule.Trim() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLine_TrimLeft(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		cutset string
		text   string
		want   string
	}{
		{"empty", "", "", ""},
		{"no trim", " ", "text", "text"},
		{"single line", " ", " text ", "text "},
		{"multi line", " ", " line1 \n line2 ", "line1 \nline2 "},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &templatefuncs.LineModule{}

			if got := m.TrimLeft(tt.cutset, tt.text); got != tt.want {
				t.Errorf("LineModule.TrimLeft() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLine_TrimRight(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		cutset string
		text   string
		want   string
	}{
		{"empty", "", "", ""},
		{"no trim", " ", "text", "text"},
		{"single line", " ", " text ", " text"},
		{"multi line", " ", " line1 \n line2 ", " line1\n line2"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &templatefuncs.LineModule{}

			if got := m.TrimRight(tt.cutset, tt.text); got != tt.want {
				t.Errorf("LineModule.TrimRight() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLine_CutPrefix(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		prefix string
		text   string
		want   string
	}{
		{"empty", "", "", ""},
		{"no prefix", "prefix", "text", "text"},
		{"single line", "prefix", "text", "text"},
		{"multi line", "prefix", "prefixline1\nprefixline2", "line1\nline2"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &templatefuncs.LineModule{}

			if got := m.CutPrefix(tt.prefix, tt.text); got != tt.want {
				t.Errorf("LineModule.CutPrefix() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLine_CutSuffix(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		suffix string
		text   string
		want   string
	}{
		{"empty", "", "", ""},
		{"no suffix", "suffix", "text", "text"},
		{"single line", "suffix", "text", "text"},
		{"multi line", "suffix", "line1suffix\nline2suffix", "line1\nline2"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &templatefuncs.LineModule{}

			if got := m.CutSuffix(tt.suffix, tt.text); got != tt.want {
				t.Errorf("LineModule.CutSuffix() = %v, want %v", got, tt.want)
			}
		})
	}
}
