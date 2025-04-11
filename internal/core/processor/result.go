package processor

import (
	"errors"
	"fmt"

	"github.com/idelchi/godyl/internal/tools"
	"github.com/idelchi/godyl/pkg/path/file"
)

// result holds the result of processing a tool.
type result struct {
	Tool    *tools.Tool
	Details details
}

// details contains information about the processing result.
type details struct {
	err      error
	messages []string
}

// determineStatus generates a status string based on the result and handles caching if needed
func (p *Processor) determineStatus(r result) string {
	tool := r.Tool

	if r.Details.err != nil {
		if isExpectedError(r.Details.err) {
			// For expected errors like "Already up to date"
			status := fmt.Sprintf("Info: %v", r.Details.err)

			// Try to save cache even for expected errors
			p.tryUpdateCache(tool)

			return status
		}

		// For unexpected errors
		status := fmt.Sprintf("Error: %v", r.Details.err)
		if len(r.Details.messages) > 0 && r.Details.messages[0] != "" {
			status += fmt.Sprintf(" (%s)", r.Details.messages[0])
		}
		return status
	}

	p.tryUpdateCache(tool)

	return "Success"
}

// tryUpdateCache attempts to update the cache with tool version information
func (p *Processor) tryUpdateCache(tool *tools.Tool) {
	if tool.Version.Version == "" {
		return
	}

	toolPath := file.New(tool.Output, tool.Exe.Name).Path()
	if err := p.cache.Save(toolPath, tool.Version.Version); err != nil {
		p.log.Error("  failed to save cache: %v", err)
	}
}

// isExpectedError checks if the error is one that doesn't indicate a complete failure.
func isExpectedError(err error) bool {
	return errors.Is(err, tools.ErrAlreadyExists) ||
		errors.Is(err, tools.ErrUpToDate) ||
		errors.Is(err, tools.ErrRequiresUpdate) ||
		errors.Is(err, tools.ErrDoesNotHaveTags) ||
		errors.Is(err, tools.ErrDoesHaveTags) ||
		errors.Is(err, tools.ErrSkipped)
}
