package snek

import (
	"cmp"
	"go/build"
	"runtime"
	"runtime/debug"
	"slices"
	"strings"
	"text/template"

	"github.com/friendly-fhir/fhenix/internal/ansi"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var (
	FormatLink    = ansi.Format(ansi.Underline, ansi.FGWhite)
	FormatHeading = ansi.FGYellow
	FormatFlag    = ansi.FGCyan
	FormatArg     = ansi.Format(ansi.Bold, ansi.FGWhite)
	FormatCommand = ansi.FGGreen
	FormatKeyword = ansi.FGCyan
	FormatStrong  = ansi.Bold
	FormatCall    = ansi.Format(ansi.FGWhite, ansi.Bold)
	FormatQuote   = ansi.FGGray
	FormatError   = ansi.FGRed
	FormatWarning = ansi.FGYellow
	FormatInfo    = ansi.FGBlue

	funcs = template.FuncMap{
		"AppName": appName,
		"Flags":   commandFlags.Get,

		"FormatLink":    format(FormatLink),
		"FormatHeading": format(FormatHeading),
		"FormatFlag":    format(FormatFlag),
		"FormatArg":     format(FormatArg),
		"FormatCommand": format(FormatCommand),
		"FormatKeyword": format(FormatKeyword),
		"FormatStrong":  format(FormatStrong),
		"FormatCall":    format(FormatCall),
		"FormatQuote":   format(FormatQuote),
		"FormatError":   format(FormatError),
		"FormatWarning": format(FormatWarning),
		"FormatInfo":    format(FormatInfo),

		"ToUpper": strings.ToUpper,
		"ToLower": strings.ToLower,

		"GoVersion":     runtime.Version,
		"GoArch":        variable(build.Default.GOARCH),
		"GoOS":          variable(build.Default.GOOS),
		"GoBuildTags":   variable(strings.Join(build.Default.BuildTags, ",")),
		"VCS":           vcs,
		"VCSRevision":   vcsRevision,
		"VCSTime":       vcsTime,
		"OrderedGroups": orderedGroups,

		"PrefixLines": prefix,
		"SuffixLines": suffix,
		"Indent":      indent,
		"FitColumns":  fitColumns,
		"FitTerm":     fitTerm,

		"FlagUsages": flagUsages,
	}
)

func appName(cmd *cobra.Command) string {
	return getAppName(cmd.Context())
}

func fitTerm(content string) string {
	width, _, err := term.GetSize(0)
	if err != nil {
		return content
	}
	const (
		minWidth = 60
		maxWidth = 100
	)
	width = max(minWidth, width)
	width = min(maxWidth, width)
	return fitColumns(width, content)
}

func fitColumns(columns int, content string) string {
	contentLines := strings.Split(content, "\n")
	var lines []string
	var sb strings.Builder
	for _, contentLine := range contentLines {
		if strings.TrimSpace(contentLine) == "" {
			next := sb.String()
			lines = append(lines, next)
			if next != "" {
				lines = append(lines, "")
			}
			sb.Reset()
			continue
		}
		words := strings.Fields(contentLine)
		for _, word := range words {
			if sb.Len()+len(word) > columns {
				lines = append(lines, sb.String())
				sb.Reset()
			}
			if sb.Len() > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(word)
		}
	}
	if sb.Len() > 0 {
		lines = append(lines, sb.String())
	}
	return strings.Join(lines, "\n")
}

func suffix(suffix, lines string) string {
	return strings.ReplaceAll(lines, "\n", suffix+"\n") + suffix
}

func prefix(prefix, lines string) string {
	return prefix + strings.ReplaceAll(lines, "\n", "\n"+prefix)
}

func indent(n int, s string) string {
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		lines[i] = strings.Repeat(" ", n) + line
	}
	return strings.Join(lines, "\n")
}

func format(displays ...ansi.Display) func(string, ...any) string {
	return func(format string, args ...any) string {
		return ansi.Format(displays...).Format(format, args...)
	}
}

func orderedGroups(cmd *cobra.Command) []*cobra.Group {
	groups := append([]*cobra.Group{}, cmd.Groups()...)
	slices.SortFunc(groups, func(lhs, rhs *cobra.Group) int {
		return cmp.Compare(lhs.Title, rhs.Title)
	})
	return groups
}

func settings() []debug.BuildSetting {
	if info, ok := debug.ReadBuildInfo(); ok {
		return info.Settings
	}
	return nil
}

func vcs() string {
	for _, setting := range settings() {
		if setting.Key == "vcs" {
			return setting.Value
		}
	}
	return "unknown"
}

func vcsRevision() string {
	for _, setting := range settings() {
		if setting.Key == "vcs.revision" {
			return setting.Value
		}
	}
	return "unknown"
}

func vcsTime() string {
	for _, setting := range settings() {
		if setting.Key == "vcs.time" {
			return setting.Value
		}
	}
	return "unknown"
}

func variable[T any](v T) func() T {
	return func() T { return v }
}

func flagUsages(f *FlagSet) string {
	return f.FormattedFlagUsages(&FormatOptions{
		ArgFormat:        FormatArg,
		FlagFormat:       FormatFlag,
		DeprecatedFormat: FormatWarning,
	})
}
