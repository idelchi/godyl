// Package results manages result collection and aggregation.
package results

import (
	"errors"
	"fmt"

	"github.com/idelchi/godyl/internal/runner"
)

// Summary provides an aggregated view of all results.
type Summary struct {
	Total      int
	Successful int
	Failed     int
	Skipped    int
	Results    []runner.Result
	Errors     []ErrorDetail
}

// ErrorDetail contains detailed error information for a failed tool.
type ErrorDetail struct {
	Tool    string
	Message string
	Error   error
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
		return fmt.Errorf("1 tool failed to install")
	}

	return fmt.Errorf("%d tools failed to install", s.Failed)
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
func (s Summary) ByStatus(status runner.Status) []runner.Result {
	var results []runner.Result
	for _, r := range s.Results {
		if r.Status == status {
			results = append(results, r)
		}
	}
	return results
}
