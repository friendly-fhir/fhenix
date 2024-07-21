package templatefuncs

import "fmt"

var (
	ErrInvalidType     = fmt.Errorf("invalid type")
	ErrIndexOutOfRange = fmt.Errorf("index out of range")
)

const (
	// StringOnError is the default string returned when an error occurs in a
	// template function.
	StringOnError = "{{ an error occurred }}"
)
