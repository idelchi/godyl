// Package processor handles the processing of tool installations and management.
package processor

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/idelchi/godyl/internal/cache"
	"github.com/idelchi/godyl/internal/config/root"
	"github.com/idelchi/godyl/internal/data"
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
	cache      *cache.Cache
	progress   progress.Manager
	config     root.Config
	log        *logger.Logger
	tools      tools.Tools
	Options    []tool.ResolveOption
	NoDownload bool
}

// New creates a new Processor.
func New(toolsList tools.Tools, cfg root.Config, log *logger.Logger) *Processor {
	// Initialize cache
	var cacheManager *cache.Cache

	if !cfg.Cache.Disabled {
		cacheManager = cache.New(data.CacheFile(cfg.Cache.Dir))
	}

	// Initialize progress manager
	progressMgr := progress.NewDefaultManager(cfg.NoProgress)

	// Initialize runner
	runnerImpl := runner.NewDefaultRunner(cacheManager, log)

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

	if p.config.Tokens.GitHub == "" {
		p.config.Parallel = 1
	}

	if p.config.Parallel > 0 {
		g.SetLimit(p.config.Parallel)
	}

	p.log.Debugf("running with %d parallel downloads", p.config.Parallel)

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

			if p.Options != nil {
				runOpts = append(runOpts, runner.WithResolveOptions(p.Options...))
			}

			// Run the tool operation
			result := p.runner.Run(ctx, t, tags, runOpts...)

			// Collect the result
			p.results.Add(result)

			// Update cache if successful
			if result.Status == runner.StatusOK && p.cache != nil {
				p.updateCache(result) //nolint:contextcheck	// Unclear what this is about.
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
	if result.Tool.Version.Version == "" {
		result.Tool.Version.Version = result.Tool.GetCurrentVersion()
	}

	if result.Tool.Version.Version == "" {
		return // No version information available
	}

	now := time.Now()

	item := &cache.Item{
		ID: result.Tool.ID(),
		// TODO(Idelchi): Name is too ambiguous and can be used for several tools (especially repos that store multiple
		// tools), consider using something else.
		Name:       result.Tool.Name,
		Version:    result.Tool.Version,
		Path:       result.Tool.AbsPath(),
		Type:       result.Tool.Source.Type.String(),
		Downloaded: now,
		Updated:    now,
	}

	if err := p.cache.Add(item); err != nil {
		p.log.Errorf("failed to update cache for %s: %v", result.Tool.Name, err)
	}
}

// presentResults formats and displays the results.
func (p *Processor) presentResults() {
	const tableMaxWidth = 100

	summary := p.results.Summary()

	// Create table formatter
	tableFormatter := presentation.NewTableFormatter(presentation.TableConfig{
		MaxWidth: tableMaxWidth,
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

		p.log.Infof("%d tools processed", len(summary.Results))
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

	const errorWrapWidth = 120

	// Create error formatter
	errorFormatter := presentation.NewErrorFormatter(presentation.ErrorConfig{
		WrapWidth: errorWrapWidth,
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
