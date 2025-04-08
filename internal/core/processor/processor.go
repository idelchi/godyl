// Package processor handles the processing of tool installations and management.
package processor

import (
	"errors"
	"fmt"
	"path/filepath"
	"sync"

	"golang.org/x/sync/errgroup"

	cachehandler "github.com/idelchi/godyl/internal/cache"
	"github.com/idelchi/godyl/internal/cache/cache"
	"github.com/idelchi/godyl/internal/config"
	"github.com/idelchi/godyl/internal/tools"
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

	// Setup concurrency
	resultCh := make(chan result)

	// Create error group for concurrent processing
	g := &errgroup.Group{}
	if p.config.Tool.Parallel > 0 {
		g.SetLimit(p.config.Tool.Parallel)
		p.log.Info("running with %d parallel downloads", p.config.Tool.Parallel)
	}

	// Start collecting results
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for r := range resultCh {
			p.handleResult(r)
		}
	}()

	// Process each tool concurrently
	for i := range p.tools {
		tool := p.tools[i]
		g.Go(func() error {
			p.processTool(tool, tags, withoutTags, resultCh)
			return nil
		})
	}

	// Wait for all tool processing to complete
	g.Wait()

	// Close the result channel and wait for result collection to complete
	close(resultCh)
	wg.Wait()

	// Return error if any tools failed to install
	if p.hasErrors {
		return fmt.Errorf("one or more %w", ErrToolsFailedToInstall)
	}

	return nil
}

// processTool processes an individual tool and sends the result to the result channel.
func (p *Processor) processTool(tool *tools.Tool, tags, withoutTags []string, resultCh chan<- result) {
	// Apply defaults and resolve tool configuration
	tool.ApplyDefaults(p.defaults, p.cache)

	if err := tool.Resolve(tags, withoutTags); err != nil {
		resultCh <- result{Tool: tool, Err: err}
		return
	}

	// Handle dry run
	if p.config.Root.Dry {
		resultCh <- result{Tool: tool}
		return
	}

	// Apply SSL verification setting
	if p.config.Tool.NoVerifySSL {
		tool.NoVerifySSL = true
	}

	// Download the tool
	msg, found, err := tool.Download()
	resultCh <- result{Tool: tool, Found: found, Err: err, Msg: msg}

	if err != nil {
		return
	}
}

// handleResult processes the result of a tool operation.
func (p *Processor) handleResult(r result) {
	tool := r.Tool

	p.log.Info("")
	p.log.Info("%s", tool.Name)
	p.log.Debug("configuration:")
	p.log.Debug("-------")
	p.log.Debug("%s", pretty.YAMLMasked(tool))
	p.log.Debug("-------")

	if r.Err != nil {
		p.handleToolError(tool, r.Err, r.Msg)
	}

	if r.Err == nil || errors.Is(r.Err, tools.ErrRequiresUpdate) {
		p.logToolSuccess(tool, r.Found)
	}
}

// handleToolError logs errors encountered during tool processing.
func (p *Processor) handleToolError(tool *tools.Tool, err error, msg string) {
	// These errors are expected and don't indicate a failure
	isExpectedError := errors.Is(err, tools.ErrAlreadyExists) ||
		errors.Is(err, tools.ErrUpToDate) ||
		errors.Is(err, tools.ErrRequiresUpdate) ||
		errors.Is(err, tools.ErrDoesNotHaveTags) ||
		errors.Is(err, tools.ErrDoesHaveTags) ||
		errors.Is(err, tools.ErrSkipped)

	if isExpectedError {
		p.log.Warn("  %v", err)
		if err := p.cache.Save(file.New(tool.Output, tool.Exe.Name).Path(), tool.Version.Version); err != nil {
			p.log.Error("  failed to save cache: %v", err)
		}
		return
	}

	// Unexpected error - mark as installation failure
	p.hasErrors = true
	p.log.Error("  failed to install")
	p.log.Debug("configuration:")
	p.log.Debug("-------")
	p.log.Debug("%s", pretty.JSONMasked(tool))
	p.log.Debug("-------")
	p.log.Error("  %v", err)
	p.log.Error("  %s", msg)
}

// logToolSuccess logs information about successfully processed tools.
func (p *Processor) logToolSuccess(tool *tools.Tool, found file.File) {
	if tool.Version.Version != "" {
		p.log.Info("  version: %s", tool.Version.Version)
	}

	p.log.Info("  picked download %q", filepath.Base(tool.Path))

	if tool.Mode == "find" {
		p.log.Info("  picked file %q", found)
		p.log.Info("  installed successfully at %q", filepath.Join(tool.Output, tool.Exe.Name))

		// Log aliases if any
		if tool.Aliases != nil {
			p.log.Info("  symlinks:")
			for _, alias := range tool.Aliases {
				p.log.Info("    - %q", filepath.Join(tool.Output, alias))
			}
		}

		if err := p.cache.Save(file.New(tool.Output, tool.Exe.Name).Path(), tool.Version.Version); err != nil {
			p.log.Error("  failed to save cache: %v", err)
		}
	} else {
		p.log.Info("  extracted to %q", tool.Output)
	}
}
