// Package processor handles the processing of tool installations and management.
package processor

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/hashicorp/go-getter/v2"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"golang.org/x/sync/errgroup"

	cachehandler "github.com/idelchi/godyl/internal/cache"
	"github.com/idelchi/godyl/internal/cache/cache"
	"github.com/idelchi/godyl/internal/config"
	"github.com/idelchi/godyl/internal/tools"
	"github.com/idelchi/godyl/internal/ui/progress"
	"github.com/idelchi/godyl/pkg/logger"
	"github.com/idelchi/godyl/pkg/path/file"
	"github.com/idelchi/godyl/pkg/pretty"
)

// ErrToolsFailedToInstall is returned when one or more tools failed to install.
var ErrToolsFailedToInstall = errors.New("tools failed to install")

// result holds the result of processing a tool.
// It's used internally to pass information between goroutines.
type result struct {
	Tool  *tools.Tool
	Found file.File
	Err   error
	Msg   string
}

// Processor handles tool installation and management.
type Processor struct {
	tools     tools.Tools
	defaults  tools.Defaults
	config    config.Config
	log       *logger.Logger
	cache     *cache.Cache
	hasErrors bool
	results   []result   // Slice to store results for deferred logging
	mu        sync.Mutex // Mutex to protect results slice
}

// New creates a new Processor.
func New(toolsList tools.Tools, defaults tools.Defaults, cfg config.Config, log *logger.Logger) *Processor {
	return &Processor{
		tools:    toolsList,
		defaults: defaults,
		config:   cfg,
		log:      log,
	}
}

// Process installs and manages tools with the given tags.
func (p *Processor) Process(tags, withoutTags []string) error {
	cache, err := cachehandler.New(p.config.Root.Cache.Dir, p.config.Root.Cache.Type)
	if err != nil {
		return fmt.Errorf("creating cache: %w", err)
	}
	p.cache = cache
	p.results = make([]result, 0, len(p.tools)) // Initialize results slice

	// Setup concurrency and progress bar container
	resultCh := make(chan result)
	var progressTrackers []*progress.PrettyProgressTracker
	var progressMu sync.Mutex

	// Create error group for concurrent processing
	g, _ := errgroup.WithContext(context.Background())
	if p.config.Tool.Parallel > 0 {
		g.SetLimit(p.config.Tool.Parallel)
		p.log.Info("running with %d parallel downloads", p.config.Tool.Parallel)
	}

	// Goroutine to collect results from the channel and store them
	var wg sync.WaitGroup // Use WaitGroup to ensure collector finishes
	wg.Add(1)
	go func() {
		defer wg.Done()
		for r := range resultCh {
			p.collectResult(r)
		}
	}()

	// Process each tool concurrently
	for i := range p.tools {
		// Capture the pointer for this iteration directly.
		tool := p.tools[i]
		if tool == nil { // Skip nil tools if any
			continue
		}
		// Log tool name before starting processing/download
		g.Go(func() error {
			// Create the progress tracker here, passing the specific tool for this goroutine
			progressTracker := progress.NewPrettyProgressTracker(tool)

			// Store the tracker for later cleanup
			progressMu.Lock()
			progressTrackers = append(progressTrackers, progressTracker)
			progressMu.Unlock()

			// Pass the tracker to processTool
			p.processTool(tool, tags, withoutTags, resultCh, progressTracker)
			return nil
		})
	}

	// Wait for all tool processing goroutines to complete
	err = g.Wait()

	// Close the result channel
	close(resultCh)

	// Wait for the collector goroutine to finish processing all results
	wg.Wait()

	// Wait for progress bars to finish rendering
	for _, tracker := range progressTrackers {
		tracker.Wait()
	}

	// Now, log all collected results sequentially
	p.logFinalResults()

	// Handle potential errors from errgroup (if any weren't handled via channel)
	if err != nil && !p.hasErrors {
		p.log.Error("unexpected error during processing: %v", err)
		return fmt.Errorf("unexpected error during processing: %w", err)
	}

	// Return final status based on handled results
	if p.hasErrors {
		return fmt.Errorf("one or more %w", ErrToolsFailedToInstall)
	}

	return nil
}

// processTool processes an individual tool and sends the result to the result channel.
func (p *Processor) processTool(
	tool *tools.Tool,
	tags, withoutTags []string,
	resultCh chan<- result,
	progressTracker getter.ProgressTracker, // Receive the tracker instance
) {
	// Apply defaults and resolve tool configuration
	tool.ApplyDefaults(p.defaults, p.cache)

	if err := tool.Resolve(tags, withoutTags); err != nil {
		resultCh <- result{Tool: tool, Err: err}
		return
	}

	// Handle dry run
	if p.config.Root.Dry {
		resultCh <- result{Tool: tool} // Send result even for dry run for consistent logging order
		return
	}

	// Apply SSL verification setting
	if p.config.Tool.NoVerifySSL {
		tool.NoVerifySSL = true
	}

	// Download the tool, passing the progress tracker
	msg, found, err := tool.Download(progressTracker)
	resultCh <- result{Tool: tool, Found: found, Err: err, Msg: msg}

	// Error handling is done via the result channel now
}

// collectResult stores the result from a tool processing goroutine.
func (p *Processor) collectResult(r result) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.results = append(p.results, r)
	// Check if the result contains an error that should set the global error flag
	if r.Err != nil && !isExpectedError(r.Err) {
		p.hasErrors = true
	}
}

// logFinalResults iterates over collected results and logs them as a table.
func (p *Processor) logFinalResults() {
	p.log.Info("") // Add a blank line before the summary

	// Import the table package at the top of your file:
	// import (
	//   "github.com/jedib0t/go-pretty/v6/table"
	//   "github.com/jedib0t/go-pretty/v6/text"
	//   "os"
	// )

	t := table.NewWriter()
	// We'll use our own logging rather than direct output

	t.AppendHeader(table.Row{"Tool", "Version", "Output Path", "Aliases", "Status"})

	// Set up a row painter to colorize rows based on their status
	t.SetRowPainter(func(row table.Row) text.Colors {
		status, ok := row[4].(string)
		if !ok {
			return nil // default color
		}

		// Color based on status message
		if strings.HasPrefix(status, "Error:") {
			return text.Colors{text.FgRed}
		} else if strings.HasPrefix(status, "Info:") {
			return text.Colors{text.FgYellow}
		} else if status == "Successfully installed" || status == "Success" {
			return text.Colors{text.FgGreen}
		}

		return nil // default color
	})

	// Now add all the rows
	for _, r := range p.results {
		p.appendResultToTable(t, r)
	}

	// Set some style options
	t.SetStyle(table.StyleRounded)

	// Configure colors for different parts if your logger supports colors
	t.Style().Color.Header = text.Colors{text.FgBlue, text.Bold}

	// Set column configurations for better formatting
	t.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, WidthMax: 20},                                 // Tool name column
		{Number: 2, WidthMax: 15},                                 // Version column
		{Number: 3, WidthMax: 40},                                 // Output Path column
		{Number: 4, WidthMax: 30},                                 // Aliases column
		{Number: 5, WidthMax: 50, Colors: text.Colors{text.Bold}}, // Status column always bold
	})

	// Render to string and log through the logger
	p.log.Info("Tool Installation Summary:")
	p.log.Info(t.Render())
}

// appendResultToTable adds a single result as a row to the table
func (p *Processor) appendResultToTable(t table.Writer, r result) {
	tool := r.Tool

	// Format aliases as a comma-separated string if present
	aliases := ""
	if tool.Aliases != nil && len(tool.Aliases) > 0 {
		aliases = strings.Join(tool.Aliases, ", ")
	}

	// Determine status based on error type
	status := "Success"
	if r.Err != nil {
		if isExpectedError(r.Err) {
			// For expected errors like "Already up to date"
			status = fmt.Sprintf("Info: %v", r.Err)

			// Attempt to save cache even for expected non-failure errors
			if cacheErr := p.cache.Save(file.New(tool.Output, tool.Exe.Name).Path(), tool.Version.Version); cacheErr != nil {
				p.log.Error("  failed to save cache: %v", cacheErr)
			}
		} else {
			// For unexpected errors
			status = fmt.Sprintf("Error: %v", r.Err)
			if r.Msg != "" {
				status += fmt.Sprintf(" (%s)", r.Msg)
			}
		}
	} else if tool.Mode == "find" {
		// Success case with more details for find mode
		status = "Successfully installed"

		// Save cache on success
		if tool.Version.Version != "" {
			if err := p.cache.Save(file.New(tool.Output, tool.Exe.Name).Path(), tool.Version.Version); err != nil {
				p.log.Error("  failed to save cache: %v", err)
			}
		}
	}

	// Append the row to the table
	t.AppendRow(table.Row{
		tool.Exe.Name,
		tool.Version.Version,
		tool.Output,
		aliases,
		status,
	})

	// Log detailed configuration at debug level
	p.log.Debug("configuration:")
	p.log.Debug("-------")
	p.log.Debug("%s", pretty.YAMLMasked(tool))
	p.log.Debug("-------")
}

// isExpectedError checks if the error is one that doesn't indicate a complete failure.
// Moved from handleToolError to be reusable.
func isExpectedError(err error) bool {
	return errors.Is(err, tools.ErrAlreadyExists) ||
		errors.Is(err, tools.ErrUpToDate) ||
		errors.Is(err, tools.ErrRequiresUpdate) ||
		errors.Is(err, tools.ErrDoesNotHaveTags) ||
		errors.Is(err, tools.ErrDoesHaveTags) ||
		errors.Is(err, tools.ErrSkipped)
}
