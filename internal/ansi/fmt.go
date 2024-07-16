package ansi

import (
	"fmt"
	"io"
	"os"
	"reflect"
)

// Fsprintf behaves like fmt.Sprintf, except that it will conditionally support
// color formatting depending on whether the Writer is colorable.
func Fsprintf(w io.Writer, format string, args ...any) string {
	return fmt.Sprintf(format, formatted(w, args)...)
}

// Fprint behaves like fmt.Fprint, except that it will conditionally support
// color formatting depending on whether the Writer is colorable.
func Fprint(w io.Writer, args ...any) (int, error) {
	return fmt.Fprint(w, formatted(w, args)...)
}

// Fprintln behaves like fmt.Fprintln, except that it will conditionally support
// color formatting depending on whether the Writer is colorable.
func Fprintln(w io.Writer, args ...any) (int, error) {
	return fmt.Fprintln(w, formatted(w, args)...)
}

// Fprintf behaves like fmt.Fprintf, except that it will conditionally support
// color formatting depending on whether the Writer is colorable.
func Fprintf(w io.Writer, format string, args ...any) (int, error) {
	return fmt.Fprintf(w, format, formatted(w, args)...)
}

// print behaves like fmt.Print, except that it will conditionally support
// color formatting depending if os.Stdout is colorable.
func Print(args ...any) (int, error) {
	return Fprint(os.Stdout, args...)
}

// Println behaves like fmt.Println, except that it will conditionally support
// color formatting depending if os.Stdout is colorable.
func Println(args ...any) (int, error) {
	return Fprintln(os.Stdout, args...)
}

// Printf behaves like fmt.Printf, except that it will conditionally support
// color formatting depending if os.Stdout is colorable.
func Printf(format string, args ...any) (int, error) {
	return Fprintf(os.Stdout, format, args...)
}

// Eprint behaves like fmt.Print, except that it will conditionally support
// color formatting depending if os.Stderr is colorable.
func Eprint(args ...any) (int, error) {
	return Fprint(os.Stderr, args...)
}

// Eprintln behaves like fmt.Println, except that it will conditionally support
// color formatting depending if os.Stderr is colorable.
func Eprintln(args ...any) (int, error) {
	return Fprintln(os.Stderr, args...)
}

// Eprintf behaves like fmt.Printf, except that it will conditionally support
// color formatting depending if os.Stderr is colorable.
func Eprintf(format string, args ...any) (int, error) {
	return Fprintf(os.Stderr, format, args...)
}

// Sprintf behaves exactly like fmt.Sprintf, except it will disable color codes
// if NOCOLOR or NO_COLOR are set.
//
// Note: ansi.Sprintf cannot obey explicit colored writers since it coalesces
// down to a single string type. If NOCOLOR is defined, then Sprintf behaves
// identically to fmt.Sprintf.
func Sprintf(format string, args ...any) string {
	return formatFunc(format, args...)
}

func formatted(w io.Writer, args []any) []any {
	formatter := selectFormatter(w)
	newArgs := make([]any, 0, len(args))
	for _, arg := range args {
		newArgs = append(newArgs, formatter(arg))
	}
	return newArgs
}

func selectFormatter(w io.Writer) func(any) any {
	if _, ok := w.(interface{ alwaysColor() }); ok {
		return alwaysFormat
	}
	if IsColorable(w) {
		return defaultFormat
	}
	return neverFormat
}

func alwaysFormat(one any) any {
	if formatter, ok := one.(interface{ formatString() string }); ok {
		return formatter.formatString()
	}
	return one
}

func defaultFormat(one any) any {
	return one
}

func neverFormat(one any) any {
	if _, ok := one.(interface{ formatString() string }); ok {
		return ""
	}
	if rt := reflect.TypeOf(one); rt.Kind() == reflect.String {
		bytes := fmt.Append(nil, one)
		rv := reflect.New(rt)
		stripped := string(stripCodes.ReplaceAll(bytes, nil))
		rv.Elem().SetString(stripped)
		return rv.Elem().Interface()
	}
	return one
}

// StripFormat will remove all ANSI escape codes from the input format string.
func StripFormat(format string) string {
	return string(stripCodes.ReplaceAll([]byte(format), nil))
}
