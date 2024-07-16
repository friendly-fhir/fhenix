package snek

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/friendly-fhir/fhenix/internal/ansi"
)

type ctxKey string

const (
	appNameKey ctxKey = "app-name"
	stdoutKey  ctxKey = "stdout"
	stderrKey  ctxKey = "stderr"
)

func withAppName(ctx context.Context, appName string) context.Context {
	return context.WithValue(ctx, appNameKey, appName)
}

func defaultAppName() string {
	if len(os.Args[0]) == 0 {
		return "snek-cli"
	}
	return filepath.Base(os.Args[0])
}

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

func withStdout(ctx context.Context, stdout io.Writer) context.Context {
	return context.WithValue(ctx, stdoutKey, stdout)
}

func getStdout(ctx context.Context) io.Writer {
	value := ctx.Value(stdoutKey)
	if value != nil {
		if stdout, ok := value.(io.Writer); ok {
			return stdout
		}
	}
	return os.Stdout
}

func withStderr(ctx context.Context, stdout io.Writer) context.Context {
	return context.WithValue(ctx, stdoutKey, stdout)
}

func getStderr(ctx context.Context) io.Writer {
	value := ctx.Value(stderrKey)
	if value != nil {
		if stderr, ok := value.(io.Writer); ok {
			return stderr
		}
	}
	return os.Stderr
}

func Printf(ctx context.Context, format string, args ...any) {
	out := getStdout(ctx)
	ansi.Fprintf(out, format, args...)
}

func Noticef(ctx context.Context, format string, args ...any) {
	prefix := NoticePrefix(ctx)
	out := getStdout(ctx)
	ansi.Fprintf(out, "%s "+format+"\n", append([]any{prefix}, args...)...)
}

func Errorf(ctx context.Context, format string, args ...any) {
	prefix := ErrorPrefix(ctx)
	out := getStderr(ctx)
	ansi.Fprintf(out, "%s "+format+"\n", append([]any{prefix}, args...)...)
}

func Warningf(ctx context.Context, format string, args ...any) {
	prefix := WarningPrefix(ctx)
	out := getStderr(ctx)
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

func ErrorPrefix(ctx context.Context) string {
	return errPrefix(getAppName(ctx))
}

func WarningPrefix(ctx context.Context) string {
	return ansi.Sprintf("%swarning:%s",
		ansi.FGYellow,
		ansi.Reset,
	)
}

func NoticePrefix(ctx context.Context) string {
	return ansi.Sprintf("%snotice:%s",
		ansi.FGCyan,
		ansi.Reset,
	)
}
