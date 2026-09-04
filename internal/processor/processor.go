// Package processor handles the processing of tool installations and management.
package processor

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/idelchi/godyl/internal/aimatch"
	"github.com/idelchi/godyl/internal/cache"
	"github.com/idelchi/godyl/internal/config/root"
	"github.com/idelchi/godyl/internal/data"
	"github.com/idelchi/godyl/internal/match"
	"github.com/idelchi/godyl/internal/tools"
	"github.com/idelchi/godyl/internal/tools/result"
	"github.com/idelchi/godyl/internal/tools/tags"
	"github.com/idelchi/godyl/internal/tools/tool"
	"github.com/idelchi/godyl/pkg/logger"
	"github.com/idelchi/godyl/pkg/pretty"
)

type assetSuggester interface {
	Suggest(ctx context.Context, request aimatch.SuggestionRequest) (aimatch.Suggestion, error)
}

// Processor is a thin orchestrator that coordinates tool processing.
type Processor struct {
	// assetSelector is consulted only for ambiguous deterministic asset matches.
	assetSelector match.AssetSelector
	// assetSuggester diagnoses deterministic asset matching failures.
	assetSuggester assetSuggester
	// results collects outcomes from concurrent tool operations.
	results *collector
	// cache persists successful tool resolutions when enabled.
	cache *cache.Cache
	// progress coordinates download progress rendering.
	progress *progressMgr
	// config contains the resolved root configuration.
	config root.Config
	// log records processing diagnostics and results.
	log *logger.Logger
	// tools contains the tools selected for this processing run.
	tools tools.Tools
	// Options customize tool resolution.
	Options []tool.ResolveOption
	// NoDownload stops after resolving tools without installing them.
	NoDownload bool
}

// New creates a new Processor.
func New(toolsList tools.Tools, cfg root.Config, log *logger.Logger) *Processor {
	// Initialize cache
	var cacheManager *cache.Cache

	if !cfg.Cache.Disabled && !cfg.Install.Suggest {
		cacheManager = cache.New(data.CacheFile(cfg.Cache.Dir))
	}

	processor := &Processor{
		tools:    toolsList,
		config:   cfg,
		log:      log,
		results:  newCollector(),
		cache:    cacheManager,
		progress: newProgressMgr(cfg.NoProgress),
	}

	if cfg.AI.Enabled {
		processor.assetSelector = aimatch.New(cfg.AI)
	}

	if cfg.Install.Suggest {
		processor.assetSuggester = aimatch.New(cfg.AI)
	}

	return processor
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
	artifacts := newArtifactStore()

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
			result := p.runTool(ctx, t, tags, artifacts)

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
		return Summary{}, errors.Join(
			fmt.Errorf("processing tools: %w", err),
			artifacts.cleanup(),
		)
	}

	p.progress.Wait()

	summary := p.results.Summary()

	if err := artifacts.cleanup(); err != nil {
		if summary.HasErrors() {
			p.log.Errorf("failed to clean artifact downloads: %v", err)

			return summary, nil
		}

		return summary, fmt.Errorf("cleaning artifact downloads: %w", err)
	}

	return summary, nil
}

// runTool executes a tool operation and returns the result.
func (p *Processor) runTool(ctx context.Context, t *tool.Tool, tags tags.IncludeTags, artifacts *artifactStore) Result {
	if p.assetSuggester != nil {
		t.NoCache = true
	}

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
	options := p.Options
	if p.assetSelector != nil {
		options = append(slices.Clone(options), tool.WithAssetSelector(p.assetSelector))
	}

	resolveResult := t.Resolve(tags, options...)

	// Convert internal result to Result
	if !resolveResult.IsOK() {
		converted := p.convertResult(t, resolveResult)

		if p.assetSuggester != nil {
			p.addSuggestion(ctx, &converted, t, resolveResult.AsError())
		}

		return converted
	}

	// Check if we should skip download
	if p.NoDownload || p.Options != nil {
		t.DisableCache()

		return p.convertResult(t, resolveResult)
	}

	downloadResult := p.downloadTool(ctx, t, artifacts)

	return p.convertResult(t, downloadResult)
}

func (p *Processor) addSuggestion(ctx context.Context, converted *Result, t *tool.Tool, resolveErr error) {
	var selectionErr *match.SelectionError

	if !errors.As(resolveErr, &selectionErr) {
		return
	}

	suggestion, err := p.assetSuggester.Suggest(ctx, aimatch.SuggestionRequest{
		Target:          t.Name,
		Description:     t.Description,
		Source:          t.Source.Type.String(),
		Version:         t.Version.Version,
		Mode:            t.Mode.String(),
		Executable:      t.Exe.Name,
		ChecksumPattern: t.Checksum.Pattern,
		Failure:         selectionErr,
	})
	if err != nil {
		converted.SuggestionError = err

		return
	}

	converted.Suggestion = &suggestion
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
