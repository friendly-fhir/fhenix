package ansi

// Display represents a displayable format element.
type Display interface {
	codes() []byte
	len() int
	Formatter
}

// Formatter is an interface for representing things capable of formatting
// strings with ANSI control codes.
type Formatter interface {
	Format(format string, args ...any) string
}

var None = Format()
