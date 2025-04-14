//go:generate go tool enumer -type=Status -output status_enumer___generated.go -transform=lower
package tools

import "fmt"

// Result represents the outcome of a tool installation operation.
// It combines a status code with a descriptive message.
type Result struct {
	Status  Status
	Message string
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

// Wrapped creates a new Result by appending the given message to the existing one.
// Preserves the original status while extending the message context.
func (r Result) Wrapped(message string) Result {
	return Result{
		Status:  r.Status,
		Message: fmt.Sprintf("%s: %s", r.Message, message),
	}
}

// Error returns an error if the Result represents a failure.
// Returns nil for successful or skipped results.
func (r Result) Error() error {
	if r.Unsuccessful() {
		return fmt.Errorf("tool failed with status: %s: %s", r.Status, r.Message)
	}

	return nil
}

// Unsuccessful returns true if the Result status is Failed.
func (r Result) Unsuccessful() bool {
	return r.Status == Failed
}

// Successful returns true if the Result status is OK.
func (r Result) Successful() bool {
	return r.Status == OK
}

// Skipped returns true if the Result status is Skipped.
func (r Result) Skipped() bool {
	return r.Status == Skipped
}
