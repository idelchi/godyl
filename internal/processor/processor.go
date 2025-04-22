// Package processor handles the processing of tool installations and management.
package processor

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"

	"github.com/jedib0t/go-pretty/v6/text"

	"github.com/idelchi/godyl/internal/cache"
	"github.com/idelchi/godyl/internal/config"
	"github.com/idelchi/godyl/internal/tools"
	"github.com/idelchi/godyl/internal/tools/tags"
	"github.com/idelchi/godyl/internal/tools/tool"
	"github.com/idelchi/godyl/internal/ui/progress"
	"github.com/idelchi/godyl/pkg/logger"
)

// ErrToolsFailedToInstall is returned when one or more tools failed to install.
var ErrToolsFailedToInstall = errors.New("tools failed to install")

// Processor handles tool installation and management.
type Processor struct {
	config      *config.Config
	log         *logger.Logger
	cache       *cache.Cache
	progressMgr *progress.ProgressManager
	tools       tools.Tools
	results     tools.Tools
	mu          sync.Mutex
	hasErrors   bool
}

// New creates a new Processor.
func New(toolsList tools.Tools, cfg *config.Config, log *logger.Logger) *Processor {
	return &Processor{
		tools:   toolsList,
		config:  cfg,
		log:     log,
		results: make(tools.Tools, 0, len(toolsList)),
	}
}

// Process installs and manages tools with the given tags.
func (p *Processor) Process(tags tags.IncludeTags, dry bool) error {
	if !p.config.Root.Cache.Disabled {
		if err := p.initializeCache(); err != nil {
			return err
		}
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
	p.cache = cache.New(p.config.Root.Cache.Dir)

	if !p.config.Root.Cache.Disabled {
		return p.cache.Load()
	}

	return nil
}

// processTools handles the concurrent processing of all tools.
func (p *Processor) processTools(tags tags.IncludeTags, dry bool) error {
	resultCh := make(chan *tool.Tool)

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
func (p *Processor) collectResult(r *tool.Tool) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.results = append(p.results, r)

	// Set global error flag for non-expected errors
	if r.Result().AsError() != nil {
		p.hasErrors = true
	}
}

// logFinalResults iterates over collected results and logs them as a table.
func (p *Processor) logFinalResults() {
	// Initialize result table
	table := NewResultTable(
		HeaderConfig{Name: "Tool", WidthMax: 100},
		HeaderConfig{Name: "Version", WidthMax: 100},
		HeaderConfig{Name: "Output Path", WidthMax: 100},
		HeaderConfig{Name: "Aliases", WidthMax: 100},
		HeaderConfig{Name: "Status", WidthMax: 100, Bold: true},
	)

	// Add all result rows to the table
	for _, r := range p.results {

		message := r.Result().Message

		// Determine the appropriate color provider based on the status
		var color text.Colors

		switch {
		case r.Result().IsFailed():
			color = ErrorColors
			message = "failed, see below for details"
		case r.Result().IsSkipped():
			color = InfoColors
		case r.Result().IsOK():
			p.UpdateCache(r)

			color = SuccessColors
		default:
			color = DefaultColors
		}

		table.AddResult(r, color, message)
	}

	// Render the table
	// if table.writer.Length() > 0 {
	p.log.Info("") // Add a blank line before the summary
	p.log.Info("Installation Summary:")
	p.log.Info("%s", table.Render())
	// }

	// Initialize result table
	// TODO(Idelchi): Adjust to fit the new table structure
	table = NewResultTable(
		HeaderConfig{Name: "Tool", WidthMax: 100},
		HeaderConfig{Name: "Error", WidthMax: 100},
	)

	// Add all result rows to the table
	for _, r := range p.results {
		message := r.Result().Error()

		if r.Result().IsFailed() {
			table.AddResult(r, ErrorColors, message)
		}
	}

	if table.writer.Length() > 0 {
		// Render the table
		p.log.Info("") // Add a blank line before the summary
		p.log.Error("Tool Error Summary:")
		p.log.Info("%s", table.Render())
	}
}
