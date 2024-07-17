package snek

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/friendly-fhir/fhenix/internal/ansi"
)

// CommandOut returns the stdout writer from the context, if available, or
// os.Stdout.
func CommandOut(ctx context.Context) io.Writer {
	value := ctx.Value(stdoutKey)
	if value != nil {
		if stdout, ok := value.(io.Writer); ok {
			return stdout
		}
	}
	return os.Stdout
}

// CommandErr returns the stderr writer from the context, if available, or
// os.Stderr.
func CommandErr(ctx context.Context) io.Writer {
	value := ctx.Value(stderrKey)
	if value != nil {
		if stderr, ok := value.(io.Writer); ok {
			return stderr
		}
	}
	return os.Stderr
}

// DisableColor disables color output for the given context.
func DisableColor(ctx context.Context) context.Context {
	ctx = withStdout(ctx, ansi.NoColorWriter(CommandOut(ctx)))
	ctx = withStderr(ctx, ansi.NoColorWriter(CommandErr(ctx)))
	return ctx
}

// ForceColor forces color output for the given context.
func ForceColor(ctx context.Context) context.Context {
	ctx = withStdout(ctx, ansi.ColorWriter(CommandOut(ctx)))
	ctx = withStderr(ctx, ansi.ColorWriter(CommandErr(ctx)))
	return ctx
}

// Printf writes a formatted string to the output stored within the context,
// or standard output if unavailable.
func Printf(ctx context.Context, format string, args ...any) {
	out := CommandOut(ctx)
	ansi.Fprintf(out, format, args...)
}

// Noticef writes a formatted string to the output stored within the context,
// or standard output if unavailable, with a notice prefix.
func Noticef(ctx context.Context, format string, args ...any) {
	prefix := NoticePrefix(ctx)
	out := CommandOut(ctx)
	ansi.Fprintf(out, "%s "+format+"\n", append([]any{prefix}, args...)...)
}

// Warningf writes a formatted string to the output stored within the context,
// or standard error if unavailable, with a warning prefix.
func Warningf(ctx context.Context, format string, args ...any) {
	prefix := WarningPrefix(ctx)
	out := CommandErr(ctx)
	ansi.Fprintf(out, "%s "+format+"\n", append([]any{prefix}, args...)...)
}

// Errorf writes a formatted string to the output stored within the context,
// or standard error if unavailable, with an error prefix.
func Errorf(ctx context.Context, format string, args ...any) {
	prefix := ErrorPrefix(ctx)
	out := CommandErr(ctx)
	ansi.Fprintf(out, "%s "+format+"\n", append([]any{prefix}, args...)...)
}

// Panicf writes a formatted string to the output stored within the context,
// or standard error if unavailable, with a panic prefix.
func Panicf(ctx context.Context, format string, args ...any) {
	prefix := PanicPrefix(ctx)
	out := CommandErr(ctx)
	ansi.Fprintf(out, "%s "+format+"\n", append([]any{prefix}, args...)...)
}

func errPrefix(appName string) string {
	return ansi.Sprintf("%serror: %s%s:%s",
		ansi.FGRed,
		ansi.FGWhite,
		appName,
		ansi.Reset,
	)
}

func PanicPrefix(ctx context.Context) string {
	return ansi.Sprintf("%spanic: %s%s:%s",
		ansi.FGRed,
		ansi.FGWhite,
		getAppName(ctx),
		ansi.Reset,
	)
}

// ErrorPrefix returns the error prefix for the given context.
// If no application name can be extracted from the context, a default one will
// be provided.
func ErrorPrefix(ctx context.Context) string {
	return ansi.Sprintf("%serror:%s",
		ansi.FGRed,
		ansi.Reset,
	)

}

// WarningPrefix returns the warning prefix for the given context.
func WarningPrefix(ctx context.Context) string {
	return ansi.Sprintf("%swarning:%s",
		ansi.FGYellow,
		ansi.Reset,
	)
}

// NoticePrefix returns the notice prefix for the given context.
func NoticePrefix(ctx context.Context) string {
	return ansi.Sprintf("%snotice:%s",
		ansi.FGCyan,
		ansi.Reset,
	)
}

////////////////////////////////////////////////////////////////////////////////

// ctxKey is a type used for context keys in this package.
type ctxKey string

const (
	appNameKey ctxKey = "app-name"
	stdoutKey  ctxKey = "stdout"
	stderrKey  ctxKey = "stderr"
)

// withAppName sets the application name in the context.
func withAppName(ctx context.Context, appName string) context.Context {
	return context.WithValue(ctx, appNameKey, appName)
}

// defaultAppName returns the default application name.
func defaultAppName() string {
	if len(os.Args) == 0 {
		return "cli"
	}
	return filepath.Base(os.Args[0])
}

// getAppName returns the application name from the context, if evailable, or
// the default application name.
func getAppName(ctx context.Context) string {
	defaultName := defaultAppName()
	if ctx == nil {
		return defaultName
	}

	value := ctx.Value(appNameKey)
	if value != nil {
		if appName, ok := value.(string); ok {
			return appName
		}
	}
	return defaultName
}

// withStdout sets the stdout writer in the context.
func withStdout(ctx context.Context, stdout io.Writer) context.Context {
	return context.WithValue(ctx, stdoutKey, stdout)
}

// withStderr sets the stderr writer in the context.
func withStderr(ctx context.Context, stdout io.Writer) context.Context {
	return context.WithValue(ctx, stderrKey, stdout)
}
