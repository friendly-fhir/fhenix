/*
Package dedent provides a utility for dedenting multiline strings.

This package is really useful for providing multiline string literals as inputs
to cobra commands in Go.
*/
package dedent

import (
	"math"
	"strings"
	"unicode"
)

// String dedents a multiline string and returns the result as a string without
// any common prefix whitespace. This will also trim and remove any leading and
// trailing whitespace lines.
func String(s string) string {
	return strings.Join(Lines(strings.Split(s, "\n")...), "\n")
}

// Strings dedents a slice of strings and returns the result as a single string
// with the common prefix whitespace removed. This will also trim any leading
func Strings(lines ...string) string {
	return strings.Join(Lines(lines...), "\n")
}

// Lines dedents a slice of strings and returns the result as a slice of strings
// with the common prefix whitespace removed. This will also trim any leading
// and trailing whitespace lines.
func Lines(lines ...string) []string {
	lines = trimLeadingBlankLines(lines)
	lines = trimTrailingBlankLines(lines)

	prefix := commonWhitespacePrefix(lines)
	if prefix == 0 {
		return lines
	}

	newlines := make([]string, len(lines))
	for i, line := range lines {
		if len(line) >= prefix {
			newlines[i] = line[prefix:]
		}
	}
	return newlines
}

func trimLeadingBlankLines(lines []string) []string {
	for i, line := range lines {
		if strings.TrimSpace(line) != "" {
			return lines[i:]
		}
	}
	return lines
}

func trimTrailingBlankLines(lines []string) []string {
	for i := len(lines) - 1; i >= 0; i-- {
		if strings.TrimSpace(lines[i]) != "" {
			return lines[:i+1]
		}
	}
	return lines
}

func commonWhitespacePrefix(lines []string) int {
	if len(lines) == 0 {
		return 0
	}
	prefix := math.MaxInt
	for _, line := range lines {
		// Skip empty lines, since they should not contribute to whitespace prefix.
		if len(strings.TrimSpace(line)) == 0 {
			continue
		}
		prefix = min(prefix, countLeadingWhitespace(line))
	}
	if prefix == math.MaxInt {
		return 0
	}
	return prefix
}

func countLeadingWhitespace(s string) int {
	for i, ch := range s {
		if !unicode.IsSpace(ch) {
			return i
		}
	}
	return len(s)
}
