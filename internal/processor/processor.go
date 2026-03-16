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
	"github.com/idelchi/godyl/internal/tools"
	"github.com/idelchi/godyl/internal/tools/result"
	"github.com/idelchi/godyl/internal/tools/tags"
	"github.com/idelchi/godyl/internal/tools/tool"
	"github.com/idelchi/godyl/pkg/logger"
	"github.com/idelchi/godyl/pkg/pretty"
)

// Processor is a thin orchestrator that coordinates tool processing.
type Processor struct {
	results    *collector
	cache      *cache.Cache
	progress   *progressMgr
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

	return &Processor{
		tools:    toolsList,
		config:   cfg,
		log:      log,
		results:  newCollector(),
		cache:    cacheManager,
		progress: newProgressMgr(cfg.NoProgress),
	}
}

// Process installs and manages tools with the given tags.
// Returns the aggregated summary and any infrastructure error (e.g. cache load failure).
func (p *Processor) Process(tags tags.IncludeTags) (Summary, error) {
	// 1. Setup
	if p.cache != nil {
		if err := p.cache.Load(); err != nil {
			return Summary{}, fmt.Errorf("loading cache: %w", err)
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
			// Run the tool operation
			result := p.runTool(ctx, t, tags)

			// Collect the result
			p.results.Add(result)

			// Update cache if successful
			if result.Status == StatusOK && p.cache != nil {
				p.updateCache(result) //nolint:contextcheck	// Unclear what this is about.
			}

			return nil
		})
	}

	// 3. Wait for completion
	if err := g.Wait(); err != nil {
		return Summary{}, fmt.Errorf("processing tools: %w", err)
	}

	p.progress.Wait()

	return p.results.Summary(), nil
}

// runTool executes a tool operation and returns the result.
func (p *Processor) runTool(ctx context.Context, t *tool.Tool, tags tags.IncludeTags) Result {
	// Enable cache if available
	if p.cache != nil {
		t.EnableCache(p.cache)
	}

	// Log tool configuration
	p.log.Debug("Tool:")
	p.log.Debug("-------")
	p.log.Debugf("%s", pretty.YAML(t))
	p.log.Debug("-------")

	// Resolve the tool
	resolveResult := t.Resolve(tags, p.Options...)

	// Convert internal result to Result
	if !resolveResult.IsOK() {
		return p.convertResult(t, resolveResult)
	}

	// Check if we should skip download
	if p.NoDownload || p.Options != nil {
		t.DisableCache()

		return p.convertResult(t, resolveResult)
	}

	// Download the tool
	downloadResult := t.Download(ctx, p.progress.Tracker())

	return p.convertResult(t, downloadResult)
}

// convertResult converts an internal result.Result to a processor Result.
func (p *Processor) convertResult(t *tool.Tool, res result.Result) Result {
	var status Status

	switch {
	case res.IsOK():
		status = StatusOK
	case res.IsSkipped():
		status = StatusSkipped
	case res.IsFailed():
		status = StatusFailed
	}

	return Result{
		Tool:    t,
		Status:  status,
		Message: res.Message,
		Error:   res.AsError(),
		Metadata: map[string]any{
			"url":     t.URL,
			"version": t.Version.Version,
			"output":  t.Output,
		},
	}
}

// updateCache updates the cache with a successful result.
func (p *Processor) updateCache(result Result) {
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
