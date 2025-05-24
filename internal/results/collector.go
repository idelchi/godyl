// Package results manages result collection and aggregation.
package results

import (
	"sync"

	"github.com/idelchi/godyl/internal/runner"
)

// Collector safely collects results from concurrent operations.
type Collector interface {
	// Add adds a result to the collection.
	Add(result runner.Result)
	// All returns all collected results.
	All() []runner.Result
	// Summary returns an aggregated summary of all results.
	Summary() Summary
}

// DefaultCollector is the default implementation of Collector.
type DefaultCollector struct {
	mu      sync.Mutex
	results []runner.Result
}

// NewCollector creates a new DefaultCollector.
func NewCollector() *DefaultCollector {
	return &DefaultCollector{
		results: make([]runner.Result, 0),
	}
}

// Add adds a result to the collection in a thread-safe manner.
func (c *DefaultCollector) Add(result runner.Result) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.results = append(c.results, result)
}

// All returns a copy of all collected results.
func (c *DefaultCollector) All() []runner.Result {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Return a copy to prevent external modification
	results := make([]runner.Result, len(c.results))
	copy(results, c.results)
	return results
}

// Summary returns an aggregated summary of all results.
func (c *DefaultCollector) Summary() Summary {
	c.mu.Lock()
	defer c.mu.Unlock()

	summary := Summary{
		Total:      len(c.results),
		Successful: 0,
		Failed:     0,
		Skipped:    0,
		Results:    make([]runner.Result, len(c.results)),
	}

	copy(summary.Results, c.results)

	for _, r := range c.results {
		switch r.Status {
		case runner.StatusOK:
			summary.Successful++
		case runner.StatusFailed:
			summary.Failed++
			summary.Errors = append(summary.Errors, ErrorDetail{
				Tool:    r.Tool.Name,
				Message: r.Message,
				Error:   r.Error,
			})
		case runner.StatusSkipped:
			summary.Skipped++
		}
	}

	return summary
}
