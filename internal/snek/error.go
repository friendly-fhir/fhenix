package snek

import (
	"errors"
	"fmt"
	"os"
)

// StatusCode represents the response from the [Application].
//
// This provides the exact underlying error that the system encountered, and
// provides a suggestion on the exit status-code, with an optional helper
// func to exit the program with.
type StatusCode struct {
	Result error
	Code   int
}

func (s *StatusCode) Error() string {
	return s.Result.Error()
}

// Exit exits the program with the status code.
func (s *StatusCode) Exit() {
	os.Exit(s.Code)
}

// PanicError is an error that is used to represent a panic in the application.
type PanicError string

// Error returns the string representation of the panic error.
func (e PanicError) Error() string {
	return string(e)
}

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
