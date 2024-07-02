package cfg

import "fmt"

// FieldError is an error that occurred while parsing a specific config field.
type FieldError struct {
	// Field is the name of the field that caused the error.
	Field string

	// Err is the error that occurred.
	Err error
}

func (e *FieldError) Error() string {
	return fmt.Sprintf("config field %q: %v", e.Field, e.Err)
}

func (e *FieldError) Unwrap() error {
	return e.Err
}

var _ error = (*FieldError)(nil)

var (
	ErrMissingField = fmt.Errorf("mandatory field not specified")
	ErrInvalidField = fmt.Errorf("invalid field value")
)
