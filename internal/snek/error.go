package snek

import (
	"errors"
	"fmt"
)

// errNotImplemented is returned when a command is not implemented.
var errNotImplemented = errors.New("not implemented")

// UsageError returns an error that should be displayed to the user as a usage
// error. This will be accompanied with the usage instructions for the command.
func UsageError(message string) error {
	return UsageErrorf(message)
}

// UsageErrorf returns an error that should be displayed to the user as a usage
// error. This will be accompanied with the usage instructions for the command.
func UsageErrorf(format string, args ...interface{}) error {
	return &usageError{Message: fmt.Sprintf(format, args...)}
}

// IsUsageError returns true if the error is a usage error.
func IsUsageError(err error) bool {
	return errors.Is(err, errUsage)
}

var errUsage = errors.New("usage error")

type usageError struct {
	Message string
}

func (e *usageError) Error() string {
	return e.Message
}

func (e *usageError) Unwrap() error {
	return errUsage
}
