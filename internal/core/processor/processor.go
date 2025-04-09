// Package processor handles the processing of tool installations and management.
package processor

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/hashicorp/go-getter/v2"
	"github.com/vbauerster/mpb/v8"
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
	prog := mpb.New(mpb.WithWidth(60))

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
		// p.log.Info("")
		// p.log.Info("%s", tool.Name) // Log using the pointer
		g.Go(func() error {
			// Create the progress tracker here, passing the specific tool for this goroutine
			progressTracker := progress.NewMpbProgressTracker(prog, tool)
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
	prog.Wait()

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

// logFinalResults iterates over collected results and logs them.
func (p *Processor) logFinalResults() {
	// p.log.Info("") // Add a blank line before the summary
	for _, r := range p.results {
		p.logSingleResult(r)
	}
}

// logSingleResult logs the outcome of a single tool operation.
func (p *Processor) logSingleResult(r result) {
	tool := r.Tool

	// Log the tool name again before its status details
	// p.log.Info("")
	// p.log.Info("%s", tool.Name)
	p.log.Debug("configuration:")
	p.log.Debug("-------")
	p.log.Debug("%s", pretty.YAMLMasked(tool))
	p.log.Debug("-------")

	if r.Err != nil {
		p.handleToolError(tool, r.Err, r.Msg) // Use existing error handler logic
	} else {
		// Log success only if there was no error (handleResult used to check r.Err == nil || IsRequiresUpdate)
		p.logToolSuccess(tool, r.Found)
	}
}

// handleToolError logs errors encountered during tool processing.
// Note: This function now only handles the logging part, the hasErrors flag is set in collectResult.
func (p *Processor) handleToolError(tool *tools.Tool, err error, msg string) {
	if isExpectedError(err) {
		p.log.Warn("  %v", err)
		// Attempt to save cache even for expected non-failure errors like UpToDate
		if cacheErr := p.cache.Save(file.New(tool.Output, tool.Exe.Name).Path(), tool.Version.Version); cacheErr != nil {
			p.log.Error("  failed to save cache: %v", cacheErr)
		}
		return
	}

	// Log unexpected error details
	p.log.Error("  failed to install")
	p.log.Debug("configuration:")
	p.log.Debug("-------")
	p.log.Debug("%s", pretty.JSONMasked(tool))
	p.log.Debug("-------")
	p.log.Error("  %v", err)
	if msg != "" { // Only log message if it's not empty
		p.log.Error("  %s", msg)
	}
}

// logToolSuccess logs information about successfully processed tools.
func (p *Processor) logToolSuccess(tool *tools.Tool, found file.File) {
	var message strings.Builder
	message.WriteString(tool.Exe.Name)

	if tool.Version.Version != "" {
		// p.log.Info("  version: %s", tool.Version.Version)
		message.WriteString(fmt.Sprintf(" %s", tool.Version.Version))
	}

	// Only log download path if it's relevant (not dry run, etc.)
	// if tool.Path != "" && !p.config.Root.Dry {
	// 	p.log.Info("  picked download %q", file.File(tool.Path).Unescape().Path())
	// }

	message.WriteString(fmt.Sprintf(" installed to %q", tool.Output))

	if tool.Mode == "find" {
		// Check if the found file path is not empty
		// if found.Path() != "" {
		// 	p.log.Info("  picked file %q", found)
		// }
		// p.log.Info("  installed successfully at %q", filepath.Join(tool.Output, tool.Exe.Name))

		// Log aliases if any
		if tool.Aliases != nil {
			message.WriteString(fmt.Sprintf(" with aliases: %v", tool.Aliases))
			// p.log.Info("  symlinks:")
			// for _, alias := range tool.Aliases {
			// 	p.log.Info("    - %q", filepath.Join(tool.Output, alias))
			// }
		}

		// Save cache on success
		if err := p.cache.Save(file.New(tool.Output, tool.Exe.Name).Path(), tool.Version.Version); err != nil {
			p.log.Error("  failed to save cache: %v", err)
		}
	} else { // Assuming "extract" mode
		// p.log.Info("  extracted to %q", tool.Output)
		// Cache saving might not be applicable or needed for extract mode? Depends on requirements.
	}

	p.log.Info(message.String())
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
