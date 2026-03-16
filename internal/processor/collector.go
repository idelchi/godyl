package processor

import "sync"

// collector safely collects results from concurrent operations.
type collector struct {
	results []Result
	mu      sync.Mutex
}

// newCollector creates a new collector.
func newCollector() *collector {
	return &collector{
		results: make([]Result, 0),
	}
}

// Add adds a result to the collection in a thread-safe manner.
func (c *collector) Add(result Result) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.results = append(c.results, result)
}

// Summary returns an aggregated summary of all results.
func (c *collector) Summary() Summary {
	c.mu.Lock()
	defer c.mu.Unlock()

	summary := Summary{
		Total:      len(c.results),
		Successful: 0,
		Failed:     0,
		Skipped:    0,
		Results:    make([]Result, len(c.results)),
	}

	copy(summary.Results, c.results)

	for _, r := range c.results {
		switch r.Status {
		case StatusOK:
			summary.Successful++
		case StatusFailed:
			summary.Failed++

			summary.Errors = append(summary.Errors, ErrorDetail{
				Tool:    r.Tool.Name,
				Message: r.Message,
				Error:   r.Error,
			})
		case StatusSkipped:
			summary.Skipped++
		}
	}

	return summary
}
