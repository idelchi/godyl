package processor

import (
	"errors"
	"fmt"

	"github.com/idelchi/godyl/internal/aimatch"
	"github.com/idelchi/godyl/internal/tools/tool"
)

// Result represents the outcome of a tool operation.
type Result struct {
	// Error contains the underlying operation failure.
	Error error
	// Tool is the tool that was processed.
	Tool *tool.Tool
	// Metadata contains structured result details for consumers.
	Metadata map[string]any
	// Message summarizes the operation outcome.
	Message string
	// Status classifies the operation outcome.
	Status Status
	// Suggestion contains a verified AI recommendation for a matching failure.
	Suggestion *aimatch.Suggestion
	// SuggestionError explains why a matching failure could not be diagnosed.
	SuggestionError error
}

// Status represents the possible states of a tool operation.
type Status int

const (
	// StatusOK indicates a successful operation.
	StatusOK Status = iota
	// StatusSkipped indicates the operation was skipped.
	StatusSkipped
	// StatusFailed indicates the operation failed.
	StatusFailed
)

// Summary provides an aggregated view of all results.
type Summary struct {
	// Results contains every processed tool result.
	Results []Result
	// Errors contains details for failed results.
	Errors []ErrorDetail
	// Total is the number of processed tools.
	Total int
	// Successful is the number of successful tools.
	Successful int
	// Failed is the number of failed tools.
	Failed int
	// Skipped is the number of skipped tools.
	Skipped int
}

// ErrorDetail contains detailed error information for a failed tool.
type ErrorDetail struct {
	// Error is the underlying failure.
	Error error
	// Tool names the tool that failed.
	Tool string
	// Message describes the failed operation.
	Message string
}

// HasErrors returns true if there are any failed results.
func (s Summary) HasErrors() bool {
	return s.Failed > 0
}

// Error returns an aggregated error if there are any failures.
func (s Summary) Error() error {
	if !s.HasErrors() {
		return nil
	}

	if s.Failed == 1 {
		return errors.New("1 tool failed to install")
	}

	return fmt.Errorf("%d tools failed to install", s.Failed)
}

// SuggestionError returns an aggregated error for a diagnostic suggestion run.
func (s Summary) SuggestionError() error {
	if !s.HasErrors() {
		return nil
	}

	if s.Failed == 1 {
		return errors.New("1 tool remains unresolved")
	}

	return fmt.Errorf("%d tools remain unresolved", s.Failed)
}

// DetailedError returns a detailed error with all failure messages.
func (s Summary) DetailedError() error {
	if !s.HasErrors() {
		return nil
	}

	errs := make([]error, 0, len(s.Errors))

	for _, e := range s.Errors {
		if e.Error != nil {
			errs = append(errs, fmt.Errorf("%s: %s: %w", e.Tool, e.Message, e.Error))
		} else {
			errs = append(errs, fmt.Errorf("%s: %s", e.Tool, e.Message))
		}
	}

	return errors.Join(errs...)
}

// ByStatus returns all results with the given status.
func (s Summary) ByStatus(status Status) []Result {
	var results []Result

	for _, r := range s.Results {
		if r.Status == status {
			results = append(results, r)
		}
	}

	return results
}
