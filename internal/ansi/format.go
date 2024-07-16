package ansi

import (
	"fmt"
	"strings"
)

const (
	// ControlSequenceIntroducer is the prefix for ANSI escape commands.
	ControlSequenceIntroducer = "\033["

	// ControlCode is a rune representing the ANSI control code (0x1B, or 033).
	ControlCode rune = 033

	// Separator is the separator between different control codes.
	Separator = ';'

	// SGR is the suffix code for the Select Graphics Renderer control codes.
	SGRSuffix = 'm'
)

// Format is a collection of ANSI displayable formatting attributes.
func Format(displays ...Display) Display {
	return format(displays)
}

type format []Display

// Format the input format string as Sprintf would, but wrap it in an ANSI
// format and reset pair.
//
// If color is disabled, this will return only the formatted string
// without any ANSI control codes.
func (f format) Format(format string, args ...any) string {
	prefix := f.String()
	suffix := Reset.String()
	content := formatFunc(format, args...)
	if len(prefix) == 0 && len(suffix) == 0 {
		return content
	}

	sb := strings.Builder{}
	sb.Grow(len(prefix) + len(suffix) + len(content))
	sb.WriteString(prefix)
	sb.WriteString(content)
	sb.WriteString(suffix)

	return sb.String()
}

var _ Formatter = (*format)(nil)

// String implements fmt.Stringer.
//
// If color is disabled, this will return an empty string.
func (f format) String() string {
	if len(f) == 0 {
		return ""
	}
	return f.formatString()
}

var _ fmt.Stringer = (*format)(nil)

// FormatString returns this formatted string with escape codes, regardless
// of whether color is enabled.
func (f format) FormatString() string {
	return f.formatString()
}

func (f format) formatString() string {
	if len(f) == 0 {
		return ""
	}
	numBytes := 0
	for _, format := range f {
		numBytes += format.len()
	}
	if numBytes == 0 {
		return ""
	}
	codes := f.codes()
	return createFormatFunc(codes...)
}

func (f format) codes() []byte {
	if !enabled {
		return nil
	}
	length := 0
	for _, display := range f {
		length += display.len()
	}
	codes := make([]byte, 0, length)
	for _, display := range f {
		codes = append(codes, display.codes()...)
	}
	return codes
}

func (f format) len() int {
	return len(f)
}

var _ Display = (*format)(nil)
