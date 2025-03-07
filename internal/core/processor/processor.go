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

// ToolProcessor handles concurrent processing of tools.
type ToolProcessor struct {
	toolsList       tools.Tools
	defaults        tools.Defaults
	cfg             config.Config
	log             *logger.Logger
	resultCh        chan toolResult
	errGroup        *errgroup.Group
	waitGroup       *sync.WaitGroup
	hasInstallError bool
}

// toolResult struct to hold the tool and any associated error.
type toolResult struct {
	tool  *tools.Tool
	found file.File
	err   error
	msg   string
}

// New creates a new ToolProcessor.
func New(toolsList tools.Tools, defaults tools.Defaults, cfg config.Config, log *logger.Logger) *ToolProcessor {
	return &ToolProcessor{
		toolsList: toolsList,
		defaults:  defaults,
		cfg:       cfg,
		log:       log,
		resultCh:  make(chan toolResult),
		errGroup:  &errgroup.Group{},
	}
}

// Process starts processing tools with the given tags.
func (tp *ToolProcessor) Process(tags, withoutTags []string) error {
	tp.setupConcurrencyLimit()

	tp.waitGroup = &sync.WaitGroup{}
	tp.waitGroup.Add(1)

	go tp.collectResults()

	for i := range tp.toolsList {
		tool := &tp.toolsList[i]
		tp.errGroup.Go(func() error {
			return tp.processTool(tool, tags, withoutTags)
		})
	}

	if err := tp.errGroup.Wait(); err != nil {
		return fmt.Errorf("processing tools: %w", err)
	}

	close(tp.resultCh)
	tp.waitGroup.Wait()

	if tp.hasInstallError {
		return errors.New("one or more tools failed to install")
	}

	return nil
}

// setupConcurrencyLimit sets the concurrency limit if specified in the config.
func (tp *ToolProcessor) setupConcurrencyLimit() {
	if tp.cfg.Parallel > 0 {
		tp.errGroup.SetLimit(tp.cfg.Parallel)
		tp.log.Info("running with %d parallel downloads", tp.cfg.Parallel)
	}
}

// collectResults reads results from the result channel and processes them.
func (tp *ToolProcessor) collectResults() {
	defer tp.waitGroup.Done()

	for res := range tp.resultCh {
		tp.processResult(res)
	}
}

// processTool processes an individual tool.
func (tp *ToolProcessor) processTool(tool *tools.Tool, tags, withoutTags []string) error {
	tool.ApplyDefaults(tp.defaults)

	if err := tool.Resolve(tags, withoutTags); err != nil {
		tp.resultCh <- toolResult{tool: tool, err: err}
		return nil
	}

	if tp.cfg.Dry {
		tp.resultCh <- toolResult{tool: tool, err: nil}
		return nil
	}

	if tp.cfg.NoVerifySSL {
		tool.NoVerifySSL = true
	}

	msg, found, err := tool.Download()
	tp.resultCh <- toolResult{tool: tool, found: found, err: err, msg: msg}

	if err != nil {
		return nil
	}

	// Execute post-installation commands if any exist
	if len(tool.Post) > 0 {
		if err := tool.Post.Exe(); err != nil {
			tp.resultCh <- toolResult{tool: tool, err: fmt.Errorf("executing post-installation commands: %w", err), msg: ""}
			return nil
		}
	}

	return nil
}

// processResult processes the result from a tool operation.
func (tp *ToolProcessor) processResult(res toolResult) {
	tool := res.tool
	err := res.err
	msg := res.msg
	found := res.found

	tp.log.Info("")
	tp.log.Info(tool.Name)
	tp.log.Debug("configuration:")
	tp.log.Debug("-------")
	tp.log.Debug(pretty.YAMLMasked(tool))
	tp.log.Debug("-------")

	if err != nil {
		tp.handleToolError(tool, err, msg)
	} else {
		tp.logToolSuccess(tool, found)
	}
}

// handleToolError logs errors encountered during tool processing.
func (tp *ToolProcessor) handleToolError(tool *tools.Tool, err error, msg string) {
	if errors.Is(err, tools.ErrAlreadyExists) ||
		errors.Is(err, tools.ErrUpToDate) ||
		errors.Is(err, tools.ErrDoesNotHaveTags) ||
		errors.Is(err, tools.ErrDoesHaveTags) ||
		errors.Is(err, tools.ErrSkipped) {
		tp.log.Warn("  %v", err)
	} else {
		tp.hasInstallError = true // Set the flag if a tool fails to install
		tp.log.Error("  failed to install")
		tp.log.Debug("configuration:")
		tp.log.Debug("-------")
		tp.log.Debug(pretty.JSONMasked(tool))
		tp.log.Debug("-------")
		tp.log.Error("  %v", err)
		tp.log.Error("  %s", msg)
	}
}

// logToolSuccess logs information about successfully processed tools.
func (tp *ToolProcessor) logToolSuccess(tool *tools.Tool, found file.File) {
	if tool.Version.Version != "" {
		tp.log.Info("  version: %s", tool.Version.Version)
	}
	tp.log.Info("  picked download %q", filepath.Base(tool.Path))

	if tool.Mode == "find" {
		tp.log.Info("  picked file %q", found)
		tp.log.Info("  installed successfully at %q", filepath.Join(tool.Output, tool.Exe.Name))
		tp.logToolAliases(tool)
	} else {
		tp.log.Info("  extracted to %q", tool.Output)
	}
}

// logToolAliases logs any aliases for the tool.
func (tp *ToolProcessor) logToolAliases(tool *tools.Tool) {
	if tool.Aliases != nil {
		tp.log.Info("  symlinks:")
		for _, alias := range tool.Aliases {
			tp.log.Info("    - %q", filepath.Join(tool.Output, alias))
		}
	}
}
