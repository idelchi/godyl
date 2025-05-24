// Package result provides a Result struct to represent the outcome of a tool installation operation.
//
//go:generate go tool enumer -type=Status -output status_enumer___generated.go -transform=lower
package result

import (
	"errors"
	"fmt"
)

// TODO(Idelchi): Make more sane functions. Wrap and Wrapped are confusing.

// Result represents the outcome of a tool installation operation.
// It combines a status code with a descriptive message and an optional error.
type Result struct {
	err     error
	Message string
	Status  Status
}

// New creates a new Result with the specified message and status.
func New(message string, status Status) Result {
	return Result{
		Message: message,
		Status:  status,
	}
}

// Status represents the possible states of a tool installation operation.
type Status int

const (
	// OK indicates a successful tool installation.
	OK Status = iota

	// Skipped indicates the tool installation was intentionally skipped.
	Skipped

	// Failed indicates the tool installation encountered an error.
	Failed
)

func (r Result) String() string {
	return fmt.Sprintf("%s: %s", r.Status, r.Message)
}

// Wrapped creates a new Result by appending the given message to the existing one.
// Preserves the original status while extending the message context.
func (r Result) Wrapped(message string) Result {
	return Result{
		Status:  r.Status,
		Message: fmt.Sprintf("%s: %s", r.Message, message),
		err:     r.err,
	}
}

// Error implements the error interface.
func (r Result) Error() string {
	if r.err != nil {
		return fmt.Sprintf("%s: %s: %v", r.Status, r.Message, r.err)
	}

	return fmt.Sprintf("%s: %s", r.Status, r.Message)
}

// AsError returns an error if the Result represents a failure.
// Returns nil for successful or skipped results.
func (r Result) AsError() error {
	if r.IsFailed() {
		if r.err != nil {
			return r.err
		}

		return fmt.Errorf("%s: %s", r.Status, r.Message)
	}

	return nil
}

// Unwrap allows for error unwrapping.
func (r Result) Unwrap() error {
	return r.err
}

// Wrap wraps an error with additional context.
func (r Result) Wrap(err error) Result {
	if err == nil {
		return r
	}

	return Result{
		Status:  r.Status,
		Message: r.Message,
		err:     fmt.Errorf("%s: %w", r.Message, err),
	}
}

// UnsIsFaileduccessful returns true if the Result status is Failed.
func (r Result) IsFailed() bool {
	return r.Status == Failed
}

// IsOK returns true if the Result status is OK.
func (r Result) IsOK() bool {
	return r.Status == OK
}

// IsSkipped returns true if the Result status is Skipped.
func (r Result) IsSkipped() bool {
	return r.Status == Skipped
}

// WithFailed creates a new Result with a Failed status, message, and optional errors.
// If multiple errors are provided, they're joined using errors.Join.
func WithFailed(message string, errs ...error) Result {
	return Result{
		Status:  Failed,
		Message: message,
		err:     errors.Join(errs...),
	}
}

// WithOK creates a new Result with an OK status and the provided message.
func WithOK(message string) Result {
	return New(message, OK)
}

// WithSkipped creates a new Result with a Skipped status and the provided message.
func WithSkipped(message string) Result {
	return New(message, Skipped)
}
