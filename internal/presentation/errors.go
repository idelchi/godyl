// Package presentation handles all UI and formatting logic.
package presentation

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/idelchi/godyl/internal/results"
)

// ErrorFormatter handles error message formatting.
type ErrorFormatter struct {
	config ErrorConfig
}

// ErrorConfig configures the error formatter.
type ErrorConfig struct {
	Format    ErrorFormat
	WrapWidth int
}

// ErrorFormat represents the output format for errors.
type ErrorFormat string

const (
	// ErrorFormatText formats errors as plain text.
	ErrorFormatText ErrorFormat = "text"
	// ErrorFormatJSON formats errors as JSON.
	ErrorFormatJSON ErrorFormat = "json"
)

// NewErrorFormatter creates a new error formatter.
func NewErrorFormatter(config ErrorConfig) *ErrorFormatter {
	return &ErrorFormatter{
		config: config,
	}
}

// FormatErrors formats error details based on the configured format.
func (f *ErrorFormatter) FormatErrors(errors []results.ErrorDetail) (string, error) {
	if len(errors) == 0 {
		return "", nil
	}

	switch f.config.Format {
	case ErrorFormatJSON:
		return f.formatJSON(errors)
	case ErrorFormatText:
		fallthrough
	default:
		return f.formatText(errors), nil
	}
}

// FormatSummary formats a summary message.
func (f *ErrorFormatter) FormatSummary(summary results.Summary) string {
	var parts []string

	if summary.Successful > 0 {
		parts = append(parts, fmt.Sprintf("%d successful", summary.Successful))
	}

	if summary.Failed > 0 {
		parts = append(parts, fmt.Sprintf("%d failed", summary.Failed))
	}

	if summary.Skipped > 0 {
		parts = append(parts, fmt.Sprintf("%d skipped", summary.Skipped))
	}

	if len(parts) == 0 {
		return "No tools processed"
	}

	return fmt.Sprintf("Processed %d tools: %s", summary.Total, strings.Join(parts, ", "))
}

// formatJSON formats errors as JSON.
func (f *ErrorFormatter) formatJSON(errors []results.ErrorDetail) (string, error) {
	type jsonError struct {
		Tool    string `json:"tool"`
		Message string `json:"message"`
		Error   string `json:"error,omitempty"`
	}

	jsonErrors := make([]jsonError, 0, len(errors))

	for _, e := range errors {
		je := jsonError{
			Tool:    e.Tool,
			Message: e.Message,
		}
		if e.Error != nil {
			je.Error = e.Error.Error()
		}

		jsonErrors = append(jsonErrors, je)
	}

	bytes, err := json.MarshalIndent(jsonErrors, "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshaling errors to JSON: %w", err)
	}

	return string(bytes), nil
}

// formatText formats errors as plain text.
func (f *ErrorFormatter) formatText(errors []results.ErrorDetail) string {
	var sb strings.Builder

	for i, e := range errors {
		if i > 0 {
			sb.WriteString("\n\n")
		}

		// Tool name
		sb.WriteString(e.Tool + "\n")

		// Error details if present
		if e.Error != nil {
			sb.WriteString(fmt.Sprintf("%v", e.Error))
		}
	}

	return sb.String()
}
