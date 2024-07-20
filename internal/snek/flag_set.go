package snek

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/friendly-fhir/fhenix/internal/ansi"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// FlagCompleters is a map of completion functions for flags.
type FlagCompleters map[string]Completer

// Flag is an interface for setting options on a flag.
type Flag struct {
	flag       *pflag.Flag
	fs         *FlagSet
	completers *FlagCompleters
}

// SetCompleter sets the completer for the flag.
func (f *Flag) SetCompleter(completer Completer) *Flag {
	(*f.completers)[f.flag.Name] = completer
	return f
}

// MarkHidden marks the flag as hidden.
func (f *Flag) MarkHidden() *Flag {
	// This function can only error if a flag doesn't exist -- but this option is
	// only presented to flags that have been created as part of the flagset.
	_ = f.fs.fs.MarkHidden(f.flag.Name)
	return f
}

// MarkRequired marks the flag as required.
func (f *Flag) MarkRequired(required bool) *Flag {
	f.fs.transforms = append(f.fs.transforms, func(cmd *cobra.Command) {
		_ = cmd.MarkFlagRequired(f.flag.Name)
	})
	return f
}

// MarkDeprecated marks the flag as deprecated.
func (f *Flag) MarkDeprecated(message string) *Flag {
	_ = f.fs.fs.MarkDeprecated(f.flag.Name, message)
	return f
}

// FlagSet is a wrapper around pflag.FlagSet that provides a more fluent API for
// defining flags.
type FlagSet struct {
	completers FlagCompleters

	name string

	fs *pflag.FlagSet

	transforms []func(cmd *cobra.Command)
}

// FlagSets returns a slice of FlagSetss.
func FlagSets(fs ...*FlagSet) []*FlagSet {
	return fs
}

// NewFlagSet creates a new FlagSet with the specified name.
func NewFlagSet(name string) *FlagSet {
	return &FlagSet{
		name: name,
		fs:   pflag.NewFlagSet(name, pflag.ContinueOnError),

		completers: make(FlagCompleters),
	}
}

// RequireTogether marks the flags as required together.
func (fs *FlagSet) RequireTogether(flags ...*Flag) {
	fs.transforms = append(fs.transforms, func(cmd *cobra.Command) {
		var labels []string
		for _, flag := range flags {
			labels = append(labels, flag.flag.Name)
		}
		cmd.MarkFlagsRequiredTogether(labels...)
	})
}

// RequireOneOf marks the flags as requiring at least one of them to be set.
func (fs *FlagSet) RequireOneOf(flags ...*Flag) {
	fs.transforms = append(fs.transforms, func(cmd *cobra.Command) {
		labels := make([]string, 0, len(flags))
		for _, flag := range flags {
			labels = append(labels, flag.flag.Name)
		}
		cmd.MarkFlagsOneRequired(labels...)
	})
}

// MarkMutuallyExclusive marks the flags as mutually exclusive.
func (fs *FlagSet) MarkMutuallyExclusive(flags ...*Flag) {
	fs.transforms = append(fs.transforms, func(cmd *cobra.Command) {
		labels := make([]string, 0, len(flags))
		for _, flag := range flags {
			labels = append(labels, flag.flag.Name)
		}
		cmd.MarkFlagsMutuallyExclusive(labels...)
	})
}

// String defines a string flag with the specified name, default value, and
// usage string.
func (fs *FlagSet) String(out *string, name, value, usage string) *Flag {
	fs.fs.StringVar(out, name, value, usage)
	return fs.flag(name)
}

// StringP defines a string flag with the specified name, shorthand, default
// value, and usage string.
func (fs *FlagSet) StringP(out *string, name, shorthand, value, usage string) *Flag {
	fs.fs.StringVarP(out, name, shorthand, value, usage)
	return fs.flag(name)
}

// Int defines an int flag with the specified name, default value, and usage
// string.
func (fs *FlagSet) Int(out *int, name string, value int, usage string) *Flag {
	fs.fs.IntVar(out, name, value, usage)
	return fs.flag(name)
}

// IntP defines an int flag with the specified name, shorthand, default value,
// and usage string.
func (fs *FlagSet) IntP(out *int, name, shorthand string, value int, usage string) *Flag {
	fs.fs.IntVarP(out, name, shorthand, value, usage)
	return fs.flag(name)
}

// Int8 defines an int8 flag with the specified name, default value, and usage
// string.
func (fs *FlagSet) Int8(out *int8, name string, value int8, usage string) *Flag {
	fs.fs.Int8Var(out, name, value, usage)
	return fs.flag(name)
}

// Int8P defines an int8 flag with the specified name, shorthand, default value,
// and usage string.
func (fs *FlagSet) Int8P(out *int8, name, shorthand string, value int8, usage string) *Flag {
	fs.fs.Int8VarP(out, name, shorthand, value, usage)
	return fs.flag(name)
}

// Int16 defines an int16 flag with the specified name, default value, and usage
// string.
func (fs *FlagSet) Int16(out *int16, name string, value int16, usage string) *Flag {
	fs.fs.Int16Var(out, name, value, usage)
	return fs.flag(name)
}

// Int16P defines an int16 flag with the specified name, shorthand, default value,
// and usage string.
func (fs *FlagSet) Int16P(out *int16, name, shorthand string, value int16, usage string) *Flag {
	fs.fs.Int16VarP(out, name, shorthand, value, usage)
	return fs.flag(name)
}

// Int32 defines an int32 flag with the specified name, default value, and usage
// string.
func (fs *FlagSet) Int32(out *int32, name string, value int32, usage string) *Flag {
	fs.fs.Int32Var(out, name, value, usage)
	return fs.flag(name)
}

// Int32P defines an int32 flag with the specified name, shorthand, default value,
// and usage string.
func (fs *FlagSet) Int32P(out *int32, name, shorthand string, value int32, usage string) *Flag {
	fs.fs.Int32VarP(out, name, shorthand, value, usage)
	return fs.flag(name)
}

// Int64 defines an int64 flag with the specified name, default value, and usage
// string.
func (fs *FlagSet) Int64(out *int64, name string, value int64, usage string) *Flag {
	fs.fs.Int64Var(out, name, value, usage)
	return fs.flag(name)
}

// Int64P defines an int64 flag with the specified name, shorthand, default value,
// and usage string.
func (fs *FlagSet) Int64P(out *int64, name, shorthand string, value int64, usage string) *Flag {
	fs.fs.Int64VarP(out, name, shorthand, value, usage)
	return fs.flag(name)
}

// Uint defines a uint flag with the specified name, default value, and usage
// string.
func (fs *FlagSet) Uint(out *uint, name string, value uint, usage string) *Flag {
	fs.fs.UintVar(out, name, value, usage)
	return fs.flag(name)
}

// UintP defines a uint flag with the specified name, shorthand, default value,
// and usage string.
func (fs *FlagSet) UintP(out *uint, name, shorthand string, value uint, usage string) *Flag {
	fs.fs.UintVarP(out, name, shorthand, value, usage)
	return fs.flag(name)
}

// Uint8 defines a uint8 flag with the specified name, default value, and usage
// string.
func (fs *FlagSet) Uint8(out *uint8, name string, value uint8, usage string) *Flag {
	fs.fs.Uint8Var(out, name, value, usage)
	return fs.flag(name)
}

// Uint8P defines a uint8 flag with the specified name, shorthand, default value,
// and usage string.
func (fs *FlagSet) Uint8P(out *uint8, name, shorthand string, value uint8, usage string) *Flag {
	fs.fs.Uint8VarP(out, name, shorthand, value, usage)
	return fs.flag(name)
}

// Uint16 defines a uint16 flag with the specified name, default value, and usage
// string.
func (fs *FlagSet) Uint16(out *uint16, name string, value uint16, usage string) *Flag {
	fs.fs.Uint16Var(out, name, value, usage)
	return fs.flag(name)
}

// Uint16P defines a uint16 flag with the specified name, shorthand, default value,
// and usage string.
func (fs *FlagSet) Uint16P(out *uint16, name, shorthand string, value uint16, usage string) *Flag {
	fs.fs.Uint16VarP(out, name, shorthand, value, usage)
	return fs.flag(name)
}

// Uint32 defines a uint32 flag with the specified name, default value, and usage
// string.
func (fs *FlagSet) Uint32(out *uint32, name string, value uint32, usage string) *Flag {
	fs.fs.Uint32Var(out, name, value, usage)
	return fs.flag(name)
}

// Uint32P defines a uint32 flag with the specified name, shorthand, default value,
// and usage string.
func (fs *FlagSet) Uint32P(out *uint32, name, shorthand string, value uint32, usage string) *Flag {
	fs.fs.Uint32VarP(out, name, shorthand, value, usage)
	return fs.flag(name)
}

// Uint64 defines a uint64 flag with the specified name, default value, and usage
// string.
func (fs *FlagSet) Uint64(out *uint64, name string, value uint64, usage string) *Flag {
	fs.fs.Uint64Var(out, name, value, usage)
	return fs.flag(name)
}

// Uint64P defines a uint64 flag with the specified name, shorthand, default value,
// and usage string.
func (fs *FlagSet) Uint64P(out *uint64, name, shorthand string, value uint64, usage string) *Flag {
	fs.fs.Uint64VarP(out, name, shorthand, value, usage)
	return fs.flag(name)
}

// Bool defines a bool flag with the specified name, default value, and usage
// string.
func (fs *FlagSet) Bool(out *bool, name string, value bool, usage string) *Flag {
	fs.fs.BoolVar(out, name, value, usage)
	return fs.flag(name)
}

// BoolP defines a bool flag with the specified name, shorthand, default value,
// and usage string.
func (fs *FlagSet) BoolP(out *bool, name, shorthand string, value bool, usage string) *Flag {
	fs.fs.BoolVarP(out, name, shorthand, value, usage)
	return fs.flag(name)
}

// Float64 defines a float64 flag with the specified name, default value, and
// usage string.
func (fs *FlagSet) Float64(out *float64, name string, value float64, usage string) *Flag {
	fs.fs.Float64Var(out, name, value, usage)
	return fs.flag(name)
}

// Float64P defines a float64 flag with the specified name, shorthand, default
// value, and usage string.
func (fs *FlagSet) Float64P(out *float64, name, shorthand string, value float64, usage string) *Flag {
	fs.fs.Float64VarP(out, name, shorthand, value, usage)
	return fs.flag(name)
}

// Duration defines a time.Duration flag with the specified name, default value,
// and usage string.
func (fs *FlagSet) Duration(out *time.Duration, name string, value time.Duration, usage string) *Flag {
	fs.fs.DurationVar(out, name, value, usage)
	return fs.flag(name)
}

// DurationP defines a time.Duration flag with the specified name, shorthand,
// default value, and usage string.
func (fs *FlagSet) DurationP(out *time.Duration, name, shorthand string, value time.Duration, usage string) *Flag {
	fs.fs.DurationVarP(out, name, shorthand, value, usage)
	return fs.flag(name)
}

// StringSlice defines a string slice flag with the specified name, default
// value, and usage string.
func (fs *FlagSet) StringSlice(out *[]string, name string, value []string, usage string) *Flag {
	fs.fs.StringSliceVar(out, name, value, usage)
	return fs.flag(name)
}

// StringSliceP defines a string slice flag with the specified name, shorthand,
// default value, and usage string.
func (fs *FlagSet) StringSliceP(out *[]string, name, shorthand string, value []string, usage string) *Flag {
	fs.fs.StringSliceVarP(out, name, shorthand, value, usage)
	return fs.flag(name)
}

// IntSlice defines an int slice flag with the specified name, default value,
// and usage string.
func (fs *FlagSet) IntSlice(out *[]int, name string, value []int, usage string) *Flag {
	fs.fs.IntSliceVar(out, name, value, usage)
	return fs.flag(name)
}

// IntSliceP defines an int slice flag with the specified name, shorthand,
// default value, and usage string.
func (fs *FlagSet) IntSliceP(out *[]int, name, shorthand string, value []int, usage string) *Flag {
	fs.fs.IntSliceVarP(out, name, shorthand, value, usage)
	return fs.flag(name)
}

func (fs *FlagSet) UintSlice(out *[]uint, name string, value []uint, usage string) *Flag {
	fs.fs.UintSliceVar(out, name, value, usage)
	return fs.flag(name)
}

func (fs *FlagSet) UintSliceP(out *[]uint, name, shorthand string, value []uint, usage string) *Flag {
	fs.fs.UintSliceVarP(out, name, shorthand, value, usage)
	return fs.flag(name)
}

func (fs *FlagSet) Var(name string, value pflag.Value, usage string) *Flag {
	fs.fs.Var(value, name, usage)
	return fs.flag(name)
}

func (fs *FlagSet) VarP(value pflag.Value, name, shorthand, usage string) *Flag {
	fs.fs.VarP(value, name, shorthand, usage)
	return fs.flag(name)
}

// CompletionFuncs returns the completion functions for the flags in the FlagSet.
func (fs *FlagSet) CompletionFuncs() FlagCompleters {
	return fs.completers
}

// FlagSet returns the underlying pflag.FlagSet.
func (fs *FlagSet) FlagSet() *pflag.FlagSet {
	return fs.fs
}

// Name returns the name of the FlagSet.
func (f *FlagSet) Name() string {
	return f.name
}

// FlagUsages returns a string containing the usage information for the flags
// in the FlagSet.
func (f *FlagSet) FlagUsages() string {
	return f.FormattedFlagUsages(nil)
}

func (fs *FlagSet) flag(name string) *Flag {
	return &Flag{
		completers: &fs.completers,
		flag:       fs.fs.Lookup(name),
		fs:         fs,
	}
}

// FormatOptions is used to configure the formatting of the flag usage.
type FormatOptions struct {
	ArgFormat        ansi.Display
	FlagFormat       ansi.Display
	DeprecatedFormat ansi.Display
}

func (fo *FormatOptions) argFormat() ansi.Display {
	if fo == nil || fo.ArgFormat == nil {
		return ansi.None
	}
	return fo.ArgFormat
}

func (fo *FormatOptions) flagFormat() ansi.Display {
	if fo == nil || fo.FlagFormat == nil {
		return ansi.None
	}
	return fo.FlagFormat
}

func (fo *FormatOptions) deprecatedFormat() ansi.Display {
	if fo == nil || fo.FlagFormat == nil {
		return ansi.None
	}
	return fo.DeprecatedFormat
}

// FormattedFlagUsages returns a string containing the usage information for
// the flags in the FlagSet, formatted according to the provided options.
func (f *FlagSet) FormattedFlagUsages(opts *FormatOptions) string {
	buf := new(bytes.Buffer)

	var lines []string

	maxlen := 0
	f.fs.VisitAll(func(flag *pflag.Flag) {
		if flag.Hidden {
			return
		}

		line := ""
		if flag.Shorthand != "" && flag.ShorthandDeprecated == "" {
			line = opts.flagFormat().Format("  -%s, --%s", flag.Shorthand, flag.Name)
		} else {
			line = opts.flagFormat().Format("      --%s", flag.Name)
		}

		varname, usage := pflag.UnquoteUsage(flag)
		if varname != "" {
			line += " " + opts.argFormat().Format(varname)
		}

		// This special character will be replaced with spacing once the
		// correct alignment is calculated
		line += "\x00"
		maxlen = max(maxlen, len(ansi.StripFormat(line)))

		line += usage
		if len(flag.Deprecated) != 0 {
			line += opts.deprecatedFormat().Format(" (DEPRECATED: %s)", flag.Deprecated)
		}

		lines = append(lines, line)
	})

	for _, line := range lines {
		sidx := strings.Index(ansi.StripFormat(line), "\x00")
		sidx2 := strings.Index(line, "\x00")

		spacing := strings.Repeat(" ", maxlen-sidx)
		// maxlen + 2 comes from + 1 for the \x00 and + 1 for the (deliberate) off-by-one in maxlen-sidx
		fmt.Fprintln(buf, line[:sidx2], spacing, wrap(maxlen+2, 0, line[sidx2+1:]))
	}

	return buf.String()
}

// Wraps the string `s` to a maximum width `w` with leading indent
// `i`. The first line is not indented (this is assumed to be done by
// caller). Pass `w` == 0 to do no wrapping
func wrap(i, w int, s string) string {
	if w == 0 {
		return strings.Replace(s, "\n", "\n"+strings.Repeat(" ", i), -1)
	}

	// space between indent i and end of line width w into which
	// we should wrap the text.
	wrap := w - i

	var r, l string

	// Not enough space for sensible wrapping. Wrap as a block on
	// the next line instead.
	if wrap < 24 {
		i = 16
		wrap = w - i
		r += "\n" + strings.Repeat(" ", i)
	}
	// If still not enough space then don't even try to wrap.
	if wrap < 24 {
		return strings.Replace(s, "\n", r, -1)
	}

	// Try to avoid short orphan words on the final line, by
	// allowing wrapN to go a bit over if that would fit in the
	// remainder of the line.
	slop := 5
	wrap = wrap - slop

	// Handle first line, which is indented by the caller (or the
	// special case above)
	l, s = wrapN(wrap, slop, s)
	r = r + strings.Replace(l, "\n", "\n"+strings.Repeat(" ", i), -1)

	// Now wrap the rest
	for s != "" {
		var t string

		t, s = wrapN(wrap, slop, s)
		r = r + "\n" + strings.Repeat(" ", i) + strings.Replace(t, "\n", "\n"+strings.Repeat(" ", i), -1)
	}
	return r
}

func wrapN(i, slop int, s string) (string, string) {
	if i+slop > len(s) {
		return s, ""
	}

	w := strings.LastIndexAny(s[:i], " \t\n")
	if w <= 0 {
		return s, ""
	}
	nlPos := strings.LastIndex(s[:i], "\n")
	if nlPos > 0 && nlPos < w {
		return s[:nlPos], s[nlPos+1:]
	}
	return s[:w], s[w+1:]
}
