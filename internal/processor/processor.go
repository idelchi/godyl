// Package processor handles the processing of tool installations and management.
package processor

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/idelchi/godyl/internal/cache"
	"github.com/idelchi/godyl/internal/config"
	"github.com/idelchi/godyl/internal/presentation"
	"github.com/idelchi/godyl/internal/progress"
	"github.com/idelchi/godyl/internal/results"
	"github.com/idelchi/godyl/internal/runner"
	"github.com/idelchi/godyl/internal/tools"
	"github.com/idelchi/godyl/internal/tools/tags"
	"github.com/idelchi/godyl/internal/tools/tool"
	"github.com/idelchi/godyl/pkg/logger"
)

// Processor is a thin orchestrator that coordinates tool processing.
type Processor struct {
	runner     runner.Runner
	results    results.Collector
	cache      cache.Manager
	progress   progress.Manager
	config     config.Config
	log        *logger.Logger
	tools      tools.Tools
	Options    []tool.ResolveOption
	NoDownload bool
}

// New creates a new Processor.
func New(toolsList tools.Tools, cfg config.Config, log *logger.Logger) *Processor {
	// Initialize cache
	var cacheManager cache.Manager

	var cacheImpl *cache.Cache

	if !cfg.Cache.Disabled {
		file, _ := cache.File(cfg.Cache.Dir)
		cacheImpl = cache.New(file)
		cacheManager = cacheImpl
	}

	// Initialize progress manager
	progressMgr := progress.NewDefaultManager(cfg.NoProgress)

	// Initialize runner
	runnerImpl := runner.NewDefaultRunner(cacheImpl, log)

	// Initialize results collector
	collector := results.NewCollector()

	return &Processor{
		tools:    toolsList,
		config:   cfg,
		log:      log,
		runner:   runnerImpl,
		results:  collector,
		cache:    cacheManager,
		progress: progressMgr,
	}
}

// Process installs and manages tools with the given tags.
func (p *Processor) Process(tags tags.IncludeTags) error {
	// 1. Setup
	if p.cache != nil {
		if err := p.cache.Load(); err != nil {
			return fmt.Errorf("loading cache: %w", err)
		}
	}

	// 2. Process tools concurrently
	ctx := context.Background()
	g, ctx := errgroup.WithContext(ctx)

	if par := p.config.Parallel; par > 0 {
		g.SetLimit(par)
		p.log.Debugf("running with %d parallel downloads", par)
	}

	// Start progress tracking
	p.progress.Start()

	for _, t := range p.tools {
		// capture
		g.Go(func() error {
			// Build run options
			var runOpts []runner.RunOption
			runOpts = append(runOpts, runner.WithProgressTracker(p.progress.Tracker()))

			if p.NoDownload {
				runOpts = append(runOpts, runner.WithNoDownload())
			}

			if p.config.NoVerifySSL {
				runOpts = append(runOpts, runner.WithNoVerifySSL())
			}

			if p.Options != nil {
				runOpts = append(runOpts, runner.WithResolveOptions(p.Options...))
			}

			// Run the tool operation
			result := p.runner.Run(ctx, t, tags, runOpts...)

			// Collect the result
			p.results.Add(result)

			// Update cache if successful
			if result.Status == runner.StatusOK && p.cache != nil {
				p.updateCache(result)
			}

			return nil
		})
	}

	// 3. Wait for completion
	if err := g.Wait(); err != nil {
		return fmt.Errorf("processing tools: %w", err)
	}

	p.progress.Wait()

	// 4. Present results
	p.presentResults()

	// 5. Return summary
	summary := p.results.Summary()
	if summary.HasErrors() {
		return summary.Error()
	}

	return nil
}

// updateCache updates the cache with a successful result.
func (p *Processor) updateCache(result runner.Result) {
	item := &cache.Item{
		ID:         result.Tool.ID(),
		Name:       result.Tool.Name,
		Version:    result.Tool.Version,
		Path:       result.Tool.Output,
		Type:       string(result.Tool.Source.Type),
		Downloaded: time.Now(),
		Updated:    time.Now(),
	}

	if err := p.cache.Set(item.ID, item); err != nil {
		p.log.Errorf("failed to update cache for %s: %v", result.Tool.Name, err)
	}
}

// presentResults formats and displays the results.
func (p *Processor) presentResults() {
	summary := p.results.Summary()

	// Create table formatter
	tableFormatter := presentation.NewTableFormatter(presentation.TableConfig{
		MaxWidth: 100,
		Verbose:  p.config.Verbose > 0,
	})

	// Render table
	tableOutput := tableFormatter.RenderResults(summary.Results)

	if tableOutput == "" {
		p.log.Info("Nothing of interest to show")

		return
	}

	// Display results
	if p.config.Verbose > 0 {
		p.log.Info("")
		p.log.Info("Installation Summary:")
		p.log.Info(tableOutput)
	} else {
		p.log.Info("Done!")
	}

	// Handle errors
	if summary.HasErrors() {
		p.presentErrors(summary)
	}
}

// presentErrors formats and displays error messages.
func (p *Processor) presentErrors(summary results.Summary) {
	// Determine error format
	format := presentation.ErrorFormatText
	if p.config.ErrorFile.Path() != "" {
		format = presentation.ErrorFormatJSON
	}

	// Create error formatter
	errorFormatter := presentation.NewErrorFormatter(presentation.ErrorConfig{
		WrapWidth: 120,
		Format:    format,
	})

	// Format errors
	errorOutput, err := errorFormatter.FormatErrors(summary.Errors)
	if err != nil {
		p.log.Errorf("failed to format errors: %v", err)

		return
	}

	// Output errors
	if p.config.ErrorFile.Path() == "" {
		p.log.Error(errorOutput)
	} else {
		if err := p.config.ErrorFile.Write([]byte(errorOutput)); err != nil {
			p.log.Errorf("failed to write error output to %q: %v", p.config.ErrorFile.Path(), err)
		} else {
			p.log.Errorf("See error file %q for details", p.config.ErrorFile.Path())
		}
	}
}

// Cache initializes the cache for the processor.
// Kept for backward compatibility.
func (p *Processor) Cache() error {
	if p.cache != nil {
		return p.cache.Load()
	}

	return nil
}
