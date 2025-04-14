// Package processor handles the processing of tool installations and management.
package processor

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"

	cachehandler "github.com/idelchi/godyl/internal/cache"
	"github.com/idelchi/godyl/internal/cache/cache"
	"github.com/idelchi/godyl/internal/config"
	"github.com/idelchi/godyl/internal/tools"
	"github.com/idelchi/godyl/internal/ui/progress"
	"github.com/idelchi/godyl/pkg/logger"
	"github.com/jedib0t/go-pretty/v6/text"
)

// ErrToolsFailedToInstall is returned when one or more tools failed to install.
var ErrToolsFailedToInstall = errors.New("tools failed to install")

// Processor handles tool installation and management.
type Processor struct {
	// Tools is the list of tools to be processed.
	tools tools.Tools

	// Defaults contains default settings to apply to tools.
	defaults tools.Defaults

	// Config holds the global configuration settings.
	config config.Config

	// Log is the logger instance for output.
	log *logger.Logger

	// Cache manages tool version caching.
	cache *cache.Cache

	// HasErrors indicates if any tool processing resulted in errors.
	hasErrors bool

	// Results stores the processing results for all tools.
	results []result

	// Mu protects concurrent access to the results slice.
	mu sync.Mutex

	// ProgressMgr manages progress tracking for tool downloads.
	progressMgr *progress.ProgressManager
}

// New creates a new Processor.
func New(toolsList tools.Tools, defaults tools.Defaults, cfg config.Config, log *logger.Logger) *Processor {
	return &Processor{
		tools:    toolsList,
		defaults: defaults,
		config:   cfg,
		log:      log,
		results:  make([]result, 0),
	}
}

// Process installs and manages tools with the given tags.
func (p *Processor) Process(tags tools.IncludeTags, dry bool) error {
	if err := p.initializeCache(); err != nil {
		return err
	}

	// Initialize the progress manager
	p.progressMgr = progress.NewProgressManager(progress.DefaultOptions())

	if err := p.processTools(tags, dry); err != nil {
		return err
	}

	p.logFinalResults()

	if p.hasErrors {
		return fmt.Errorf("one or more %w", ErrToolsFailedToInstall)
	}

	return nil
}

// initializeCache initializes the cache for the processor.
func (p *Processor) initializeCache() error {
	cache, err := cachehandler.New(p.config.Root.Cache.Dir)
	if err != nil {
		return fmt.Errorf("creating cache: %w", err)
	}
	p.cache = cache
	return nil
}

// processTools handles the concurrent processing of all tools.
func (p *Processor) processTools(tags tools.IncludeTags, dry bool) error {
	resultCh := make(chan result)
	var progressTrackers []*progress.PrettyProgressTracker
	var progressMu sync.Mutex

	// Create error group for concurrent processing
	g, _ := errgroup.WithContext(context.Background())
	if p.config.Tool.Parallel > 0 {
		g.SetLimit(p.config.Tool.Parallel)
		p.log.Debug("running with %d parallel downloads", p.config.Tool.Parallel)
	}

	// Start collector goroutine
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for r := range resultCh {
			p.collectResult(r)
		}
	}()

	// Launch tool processing goroutines
	for i := range p.tools {
		tool := p.tools[i]
		if tool == nil {
			continue
		}

		g.Go(func() error {
			// Create progress tracker using the manager
			progressTracker := p.progressMgr.NewTracker(tool)

			progressMu.Lock()
			progressTrackers = append(progressTrackers, progressTracker)
			progressMu.Unlock()

			p.processOneTool(tool, tags, resultCh, progressTracker, dry)
			return nil
		})
	}

	// Wait for all processing to complete
	err := g.Wait()
	close(resultCh)
	wg.Wait()

	// Wait for progress bars to finish rendering
	for _, tracker := range progressTrackers {
		tracker.Wait()
	}

	// Stop the progress manager when done
	p.progressMgr.Stop()

	return err
}

// collectResult stores the result from a tool processing goroutine.
func (p *Processor) collectResult(r result) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.results = append(p.results, r)

	// Set global error flag for non-expected errors
	if r.Result.Error() != nil {
		p.hasErrors = true
	}
}

// logFinalResults iterates over collected results and logs them as a table.
func (p *Processor) logFinalResults() {
	// Initialize result table
	// TODO(Idelchi): Adjust to fit the new table structure
	table := NewResultTable(
		HeaderConfig{Name: "Tool", WidthMax: 100},
		HeaderConfig{Name: "Version", WidthMax: 100},
		HeaderConfig{Name: "Output Path", WidthMax: 100},
		HeaderConfig{Name: "Aliases", WidthMax: 100},
		HeaderConfig{Name: "Status", WidthMax: 100, Bold: true},
	)

	// Add all result rows to the table
	for _, r := range p.results {
		p.UpdateCache(r.Tool)

		message := r.Result.Message

		// Determine the appropriate color provider based on the status
		var color text.Colors
		switch {
		case r.Result.Unsuccessful():
			color = ErrorColors
			message = "failed, see below for details"
		case r.Result.Skipped():
			color = InfoColors
		case r.Result.Successful():
			color = SuccessColors
		default:
			color = DefaultColors
		}

		table.AddResult(r.Tool, color, message)
	}

	// Render the table
	p.log.Info("") // Add a blank line before the summary
	p.log.Info("Installation Summary:")
	p.log.Info(table.Render())

	// Initialize result table
	// TODO(Idelchi): Adjust to fit the new table structure
	table = NewResultTable(
		HeaderConfig{Name: "Tool", WidthMax: 100},
		HeaderConfig{Name: "Error", WidthMax: 100},
	)

	// Add all result rows to the table
	for _, r := range p.results {
		message := r.Result.Message

		if r.Result.Unsuccessful() {
			table.AddResult(r.Tool, ErrorColors, message)
		}
	}

	if table.writer.Length() > 0 {
		// Render the table
		p.log.Info("") // Add a blank line before the summary
		p.log.Error("Tool Error Summary:")
		p.log.Info(table.Render())
	}
}
