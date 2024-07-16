package ansi

import (
	"io"

	"golang.org/x/term"
)

// ColorWriter creates an io.Writer that forces writing of colors, even if it's
// disabled or not currently allowed.
func ColorWriter(w io.Writer) io.Writer {
	if fd, ok := w.(fdWriter); ok && enabled && term.IsTerminal(fd.Fd()) {
		return fdColorWriter{fd}
	}
	return colorWriter{w}
}

type noColorWriter struct {
	w io.Writer
}

func (w *noColorWriter) Write(p []byte) (n int, err error) {
	return w.w.Write(p)
}

func (w *noColorWriter) Close() error {
	if c, ok := w.w.(io.Closer); ok {
		return c.Close()
	}
	return nil
}

// NoColorWriter creates an io.Writer that forces writing of colors to be disabled.
func NoColorWriter(w io.Writer) io.Writer {
	return &noColorWriter{w}
}

// IsColorable checks whether the specified Writer is a colorable output destination.
//
// This will return true either if the Writer is a TTY with colors enabled, or
// if the writer is an explicit colorable writer.
func IsColorable(w io.Writer) bool {
	if _, ok := w.(interface{ alwaysColor() }); ok {
		return true
	}
	if fd, ok := w.(interface{ Fd() uintptr }); ok && enabled && term.IsTerminal(int(fd.Fd())) {
		return true
	}
	return false
}

type colorWriter struct {
	io.Writer
}

func (colorWriter) alwaysColor() {
}

type fdWriter interface {
	Fd() int
	io.Writer
}

type fdColorWriter struct {
	fdWriter
}

func (w fdColorWriter) alwaysColor() {
}
