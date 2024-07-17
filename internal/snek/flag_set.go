package snek

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/friendly-fhir/fhenix/internal/ansi"
	"github.com/spf13/pflag"
)

// FlagCompleters is a map of completion functions for flags.
type FlagCompleters map[string]Completer

// FlagSet is a wrapper around pflag.FlagSet that provides a more fluent API for
// defining flags.
type FlagSet struct {
	completers FlagCompleters

	name string

	fs *pflag.FlagSet
}

func FlagSets(fs ...*FlagSet) []*FlagSet {
	return fs
}

// FlagSetOptions is an interface for setting options on a flag.
type FlagSetOptions interface {
	// WithCompleter sets a completion function for the flag.
	WithCompleter(completer Completer) FlagSetOptions

	// MarkHidden marks the flag as hidden.
	MarkHidden()
}

type options struct {
	fs         *pflag.FlagSet
	completers *FlagCompleters
	flag       string
}

func (o *options) WithCompleter(completer Completer) FlagSetOptions {
	(*o.completers)[o.flag] = completer
	return o
}

func (o *options) MarkHidden() {
	// This function can only error if a flag doesn't exist -- but this option is
	// only presented to flags that have been created as part of the flagset.
	_ = o.fs.MarkHidden(o.flag)
}

func NewFlagSet(name string) *FlagSet {
	return &FlagSet{
		name: name,
		fs:   pflag.NewFlagSet(name, pflag.ContinueOnError),

		completers: make(FlagCompleters),
	}
}

func (fs *FlagSet) String(out *string, name, value, usage string) FlagSetOptions {
	fs.fs.StringVar(out, name, value, usage)
	return fs.options(name)
}

func (fs *FlagSet) StringP(out *string, name, shorthand, value, usage string) FlagSetOptions {
	fs.fs.StringVarP(out, name, shorthand, value, usage)
	return fs.options(name)
}

func (fs *FlagSet) Int(out *int, name string, value int, usage string) FlagSetOptions {
	fs.fs.IntVar(out, name, value, usage)
	return fs.options(name)
}

func (fs *FlagSet) IntP(out *int, name, shorthand string, value int, usage string) FlagSetOptions {
	fs.fs.IntVarP(out, name, shorthand, value, usage)
	return fs.options(name)
}

func (fs *FlagSet) Int8(out *int8, name string, value int8, usage string) FlagSetOptions {
	fs.fs.Int8Var(out, name, value, usage)
	return fs.options(name)
}

func (fs *FlagSet) Int8P(out *int8, name, shorthand string, value int8, usage string) FlagSetOptions {
	fs.fs.Int8VarP(out, name, shorthand, value, usage)
	return fs.options(name)
}

func (fs *FlagSet) Int16(out *int16, name string, value int16, usage string) FlagSetOptions {
	fs.fs.Int16Var(out, name, value, usage)
	return fs.options(name)
}

func (fs *FlagSet) Int16P(out *int16, name, shorthand string, value int16, usage string) FlagSetOptions {
	fs.fs.Int16VarP(out, name, shorthand, value, usage)
	return fs.options(name)
}

func (fs *FlagSet) Int32(out *int32, name string, value int32, usage string) FlagSetOptions {
	fs.fs.Int32Var(out, name, value, usage)
	return fs.options(name)
}

func (fs *FlagSet) Int32P(out *int32, name, shorthand string, value int32, usage string) FlagSetOptions {
	fs.fs.Int32VarP(out, name, shorthand, value, usage)
	return fs.options(name)
}

func (fs *FlagSet) Int64(out *int64, name string, value int64, usage string) FlagSetOptions {
	fs.fs.Int64Var(out, name, value, usage)
	return fs.options(name)
}

func (fs *FlagSet) Int64P(out *int64, name, shorthand string, value int64, usage string) FlagSetOptions {
	fs.fs.Int64VarP(out, name, shorthand, value, usage)
	return fs.options(name)
}

func (fs *FlagSet) Uint(out *uint, name string, value uint, usage string) FlagSetOptions {
	fs.fs.UintVar(out, name, value, usage)
	return fs.options(name)
}

func (fs *FlagSet) UintP(out *uint, name, shorthand string, value uint, usage string) FlagSetOptions {
	fs.fs.UintVarP(out, name, shorthand, value, usage)
	return fs.options(name)
}

func (fs *FlagSet) Uint8(out *uint8, name string, value uint8, usage string) FlagSetOptions {
	fs.fs.Uint8Var(out, name, value, usage)
	return fs.options(name)
}

func (fs *FlagSet) Uint8P(out *uint8, name, shorthand string, value uint8, usage string) FlagSetOptions {
	fs.fs.Uint8VarP(out, name, shorthand, value, usage)
	return fs.options(name)
}

func (fs *FlagSet) Uint16(out *uint16, name string, value uint16, usage string) FlagSetOptions {
	fs.fs.Uint16Var(out, name, value, usage)
	return fs.options(name)
}

func (fs *FlagSet) Uint16P(out *uint16, name, shorthand string, value uint16, usage string) FlagSetOptions {
	fs.fs.Uint16VarP(out, name, shorthand, value, usage)
	return fs.options(name)
}

func (fs *FlagSet) Uint32(out *uint32, name string, value uint32, usage string) FlagSetOptions {
	fs.fs.Uint32Var(out, name, value, usage)
	return fs.options(name)
}

func (fs *FlagSet) Uint32P(out *uint32, name, shorthand string, value uint32, usage string) FlagSetOptions {
	fs.fs.Uint32VarP(out, name, shorthand, value, usage)
	return fs.options(name)
}

func (fs *FlagSet) Uint64(out *uint64, name string, value uint64, usage string) FlagSetOptions {
	fs.fs.Uint64Var(out, name, value, usage)
	return fs.options(name)
}

func (fs *FlagSet) Uint64P(out *uint64, name, shorthand string, value uint64, usage string) FlagSetOptions {
	fs.fs.Uint64VarP(out, name, shorthand, value, usage)
	return fs.options(name)
}

func (fs *FlagSet) Bool(out *bool, name string, value bool, usage string) FlagSetOptions {
	fs.fs.BoolVar(out, name, value, usage)
	return fs.options(name)
}

func (fs *FlagSet) BoolP(out *bool, name, shorthand string, value bool, usage string) FlagSetOptions {
	fs.fs.BoolVarP(out, name, shorthand, value, usage)
	return fs.options(name)
}

func (fs *FlagSet) Float64(out *float64, name string, value float64, usage string) FlagSetOptions {
	fs.fs.Float64Var(out, name, value, usage)
	return fs.options(name)
}

func (fs *FlagSet) Float64P(out *float64, name, shorthand string, value float64, usage string) FlagSetOptions {
	fs.fs.Float64VarP(out, name, shorthand, value, usage)
	return fs.options(name)
}

func (fs *FlagSet) Duration(out *time.Duration, name string, value time.Duration, usage string) FlagSetOptions {
	fs.fs.DurationVar(out, name, value, usage)
	return fs.options(name)
}

func (fs *FlagSet) DurationP(out *time.Duration, name, shorthand string, value time.Duration, usage string) FlagSetOptions {
	fs.fs.DurationVarP(out, name, shorthand, value, usage)
	return fs.options(name)
}

func (fs *FlagSet) StringSlice(out *[]string, name string, value []string, usage string) FlagSetOptions {
	fs.fs.StringSliceVar(out, name, value, usage)
	return fs.options(name)
}

func (fs *FlagSet) StringSliceP(out *[]string, name, shorthand string, value []string, usage string) FlagSetOptions {
	fs.fs.StringSliceVarP(out, name, shorthand, value, usage)
	return fs.options(name)
}

func (fs *FlagSet) IntSlice(out *[]int, name string, value []int, usage string) FlagSetOptions {
	fs.fs.IntSliceVar(out, name, value, usage)
	return fs.options(name)
}

func (fs *FlagSet) IntSliceP(out *[]int, name, shorthand string, value []int, usage string) FlagSetOptions {
	fs.fs.IntSliceVarP(out, name, shorthand, value, usage)
	return fs.options(name)
}

func (fs *FlagSet) UintSlice(out *[]uint, name string, value []uint, usage string) FlagSetOptions {
	fs.fs.UintSliceVar(out, name, value, usage)
	return fs.options(name)
}

func (fs *FlagSet) UintSliceP(out *[]uint, name, shorthand string, value []uint, usage string) FlagSetOptions {
	fs.fs.UintSliceVarP(out, name, shorthand, value, usage)
	return fs.options(name)
}

func (fs *FlagSet) Var(name string, value pflag.Value, usage string) FlagSetOptions {
	fs.fs.Var(value, name, usage)
	return fs.options(name)
}

func (fs *FlagSet) VarP(value pflag.Value, name, shorthand, usage string) FlagSetOptions {
	fs.fs.VarP(value, name, shorthand, usage)
	return fs.options(name)
}

type funcFlag func(string) error

func (fn funcFlag) Set(s string) error {
	return fn(s)
}

func (fn funcFlag) Type() string {
	return "string"
}

func (fn funcFlag) String() string {
	return ""
}

var _ pflag.Value = (*funcFlag)(nil)

func (fs *FlagSet) Func(name, usage string, fn func(string) error) FlagSetOptions {
	fs.fs.Var(funcFlag(fn), name, usage)
	return fs.options(name)
}

func (fs *FlagSet) FuncP(name, shorthand, usage string, fn func(string) error) FlagSetOptions {
	fs.fs.VarP(funcFlag(fn), name, shorthand, usage)
	return fs.options(name)
}

func (fs *FlagSet) options(name string) *options {
	return &options{
		completers: &fs.completers,
		flag:       name,
		fs:         fs.fs,
	}
}

func (fs *FlagSet) CompletionFuncs() FlagCompleters {
	return nil
}

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
