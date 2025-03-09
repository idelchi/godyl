// Package processor handles the processing of tool installations and management.
package processor

import (
	"errors"
	"fmt"
	"path/filepath"
	"sync"

	"golang.org/x/sync/errgroup"

	"github.com/idelchi/godyl/internal/config"
	"github.com/idelchi/godyl/internal/tools"
	"github.com/idelchi/godyl/pkg/file"
	"github.com/idelchi/godyl/pkg/logger"
	"github.com/idelchi/godyl/pkg/pretty"
)

// ErrToolsFailedToInstall is returned when one or more tools failed to install.
var ErrToolsFailedToInstall = errors.New("tools failed to install")

// ToolResult holds the result of processing a tool.
type ToolResult struct {
	Tool  *tools.Tool
	Found file.File
	Err   error
	Msg   string
}

// ToolProcessor processes tools with concurrency support.
type ToolProcessor struct {
	toolsList tools.Tools
	defaults  tools.Defaults
	cfg       config.Config
	log       *logger.Logger

	// Components
	concurrencyManager *ConcurrencyManager
	resultHandler      ResultHandler
}

// ResultHandler handles the results of tool processing.
type ResultHandler interface {
	HandleResult(result ToolResult)
	HasErrors() bool
}

// DefaultResultHandler is the default implementation of ResultHandler.
type DefaultResultHandler struct {
	log             *logger.Logger
	hasInstallError bool
}

// ConcurrencyManager manages concurrent execution of tool processing.
type ConcurrencyManager struct {
	errGroup  *errgroup.Group
	waitGroup *sync.WaitGroup
	resultCh  chan ToolResult
	parallel  int
}

// New creates a new ToolProcessor.
func New(toolsList tools.Tools, defaults tools.Defaults, cfg config.Config, log *logger.Logger) *ToolProcessor {
	resultHandler := &DefaultResultHandler{
		log: log,
	}

	concurrencyManager := &ConcurrencyManager{
		errGroup: &errgroup.Group{},
		resultCh: make(chan ToolResult),
		parallel: cfg.Tool.Parallel,
	}

	return &ToolProcessor{
		toolsList:          toolsList,
		defaults:           defaults,
		cfg:                cfg,
		log:                log,
		concurrencyManager: concurrencyManager,
		resultHandler:      resultHandler,
	}
}

// Process starts processing tools with the given tags.
func (tp *ToolProcessor) Process(tags, withoutTags []string) error {
	// Setup concurrency
	tp.setupConcurrency()

	// Start collecting results
	go tp.collectResults()

	// Process each tool concurrently
	for i := range tp.toolsList {
		tool := &tp.toolsList[i]
		tp.concurrencyManager.errGroup.Go(func() error {
			return tp.processTool(tool, tags, withoutTags)
		})
	}

	// Wait for all tool processing to complete
	// if err := tp.concurrencyManager.errGroup.Wait(); err != nil {
	// 	return fmt.Errorf("processing tools: %w", err)
	// }
	tp.concurrencyManager.errGroup.Wait() //nolint:errcheck	// Error is checked further down instead

	// Close the result channel and wait for result collection to complete
	close(tp.concurrencyManager.resultCh)
	tp.concurrencyManager.waitGroup.Wait()

	// Check if any tools failed to install
	if tp.resultHandler.HasErrors() {
		return fmt.Errorf("one or more %w", ErrToolsFailedToInstall)
	}

	return nil
}

// setupConcurrency sets up the concurrency management.
func (tp *ToolProcessor) setupConcurrency() {
	// Set concurrency limit if specified
	if tp.cfg.Tool.Parallel > 0 {
		tp.concurrencyManager.errGroup.SetLimit(tp.cfg.Tool.Parallel)
		tp.log.Info("running with %d parallel downloads", tp.cfg.Tool.Parallel)
	}

	// Initialize wait group for result collection
	tp.concurrencyManager.waitGroup = &sync.WaitGroup{}
	tp.concurrencyManager.waitGroup.Add(1)
}

// collectResults collects and processes results from the result channel.
func (tp *ToolProcessor) collectResults() {
	defer tp.concurrencyManager.waitGroup.Done()

	for result := range tp.concurrencyManager.resultCh {
		tp.resultHandler.HandleResult(result)
	}
}

// processTool processes an individual tool.
func (tp *ToolProcessor) processTool(tool *tools.Tool, tags, withoutTags []string) error {
	// Apply defaults and resolve tool configuration
	tool.ApplyDefaults(tp.defaults)

	if err := tool.Resolve(tags, withoutTags); err != nil {
		tp.concurrencyManager.resultCh <- ToolResult{Tool: tool, Err: err}

		return fmt.Errorf("resolving tool %q: %w", tool.Name, err)
	}

	// Handle dry run
	if tp.cfg.Root.Dry {
		tp.concurrencyManager.resultCh <- ToolResult{Tool: tool}

		return nil
	}

	// Apply SSL verification setting
	if tp.cfg.Tool.NoVerifySSL {
		tool.NoVerifySSL = true
	}

	// Download the tool
	msg, found, err := tool.Download()
	tp.concurrencyManager.resultCh <- ToolResult{Tool: tool, Found: found, Err: err, Msg: msg}

	if err != nil {
		return fmt.Errorf("downloading tool %q: %w", tool.Name, err)
	}

	// Execute post-installation commands if any exist
	if len(tool.Post) > 0 {
		if err := tool.Post.Exe(); err != nil {
			tp.concurrencyManager.resultCh <- ToolResult{
				Tool: tool,
				Err:  fmt.Errorf("executing post-installation commands: %w", err),
			}
		}
	}

	return nil
}

// HandleResult processes the result of a tool operation.
func (h *DefaultResultHandler) HandleResult(result ToolResult) {
	tool := result.Tool
	err := result.Err
	msg := result.Msg
	found := result.Found

	h.log.Info("")
	h.log.Info("%s", tool.Name)
	h.log.Debug("configuration:")
	h.log.Debug("-------")
	h.log.Debug("%s", pretty.YAMLMasked(tool))
	h.log.Debug("-------")

	if err != nil {
		h.handleToolError(tool, err, msg)
	} else {
		h.logToolSuccess(tool, found)
	}
}

// HasErrors returns true if any tools failed to install.
func (h *DefaultResultHandler) HasErrors() bool {
	return h.hasInstallError
}

// handleToolError logs errors encountered during tool processing.
func (h *DefaultResultHandler) handleToolError(tool *tools.Tool, err error, msg string) {
	if errors.Is(err, tools.ErrAlreadyExists) ||
		errors.Is(err, tools.ErrUpToDate) ||
		errors.Is(err, tools.ErrDoesNotHaveTags) ||
		errors.Is(err, tools.ErrDoesHaveTags) ||
		errors.Is(err, tools.ErrSkipped) {
		h.log.Warn("  %v", err)
	} else {
		h.hasInstallError = true // Set the flag if a tool fails to install
		h.log.Error("  failed to install")
		h.log.Debug("configuration:")
		h.log.Debug("-------")
		h.log.Debug("%s", pretty.JSONMasked(tool))
		h.log.Debug("-------")
		h.log.Error("  %v", err)
		h.log.Error("  %s", msg)
	}
}

// logToolSuccess logs information about successfully processed tools.
func (h *DefaultResultHandler) logToolSuccess(tool *tools.Tool, found file.File) {
	if tool.Version.Version != "" {
		h.log.Info("  version: %s", tool.Version.Version)
	}

	h.log.Info("  picked download %q", filepath.Base(tool.Path))

	if tool.Mode == "find" {
		h.log.Info("  picked file %q", found)
		h.log.Info("  installed successfully at %q", filepath.Join(tool.Output, tool.Exe.Name))
		h.logToolAliases(tool)
	} else {
		h.log.Info("  extracted to %q", tool.Output)
	}
}

// logToolAliases logs any aliases for the tool.
func (h *DefaultResultHandler) logToolAliases(tool *tools.Tool) {
	if tool.Aliases != nil {
		h.log.Info("  symlinks:")

		for _, alias := range tool.Aliases {
			h.log.Info("    - %q", filepath.Join(tool.Output, alias))
		}
	}
}
