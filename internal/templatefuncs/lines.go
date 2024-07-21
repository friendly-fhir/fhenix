package templatefuncs

import (
	"strings"

	"github.com/friendly-fhir/fhenix/internal/dedent"
)

type LineModule struct{}

func (m *LineModule) Prefix(prefix, text string) string {
	return prefix + strings.ReplaceAll(text, "\n", "\n"+prefix)
}

func (m *LineModule) Suffix(suffix, text string) string {
	return strings.ReplaceAll(text, "\n", suffix+"\n") + suffix
}

func (m *LineModule) TrimSpace(s string) string {
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimSpace(line)
	}
	return strings.Join(lines, "\n")
}

func (m *LineModule) Trim(cutset, s string) string {
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		lines[i] = strings.Trim(line, cutset)
	}
	return strings.Join(lines, "\n")
}

func (m *LineModule) TrimLeft(cutset, s string) string {
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimLeft(line, cutset)
	}
	return strings.Join(lines, "\n")
}

func (m *LineModule) TrimRight(cutset, s string) string {
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimRight(line, cutset)
	}
	return strings.Join(lines, "\n")
}

func (m *LineModule) CutPrefix(prefix, s string) string {
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimPrefix(line, prefix)
	}
	return strings.Join(lines, "\n")
}

func (m *LineModule) CutSuffix(suffix, s string) string {
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimSuffix(line, suffix)
	}
	return strings.Join(lines, "\n")
}

func (m *LineModule) Dedent(lines string) string {
	return dedent.String(lines)
}

func (m *LineModule) Split(lines string) []string {
	return strings.Split(lines, "\n")
}

func (m *LineModule) Resize(columns int, text string) string {
	var sb strings.Builder
	lines := strings.Split(text, "\n")
	length := 0
	for _, line := range lines {
		// Preserve empty line breaks
		if strings.TrimSpace(line) == "" {
			sb.WriteString("\n")
			length = 0
			continue
		}
		tokens := strings.Fields(line)
		for i, token := range tokens {
			if i > 0 && length+len(token) > columns {
				sb.WriteString("\n")
				length = 0
			}
			sb.WriteString(token)
			sb.WriteString(" ")
			length += len(token) + 1
		}
	}
	return strings.TrimSpace(sb.String())
}

func (m *LineModule) Indent(indent int, text string) string {
	space := strings.Repeat(" ", indent)
	return space + strings.ReplaceAll(text, "\n", "\n"+space)
}

func (m *LineModule) TabIndent(indent int, text string) string {
	tab := strings.Repeat("\t", indent)
	return tab + strings.ReplaceAll(text, "\n", "\n"+tab)
}
