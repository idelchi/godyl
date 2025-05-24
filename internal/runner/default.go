// Package runner handles the execution of tool operations.
package runner

import (
	"context"

	"github.com/idelchi/godyl/internal/cache"
	"github.com/idelchi/godyl/internal/tools/result"
	"github.com/idelchi/godyl/internal/tools/tags"
	"github.com/idelchi/godyl/internal/tools/tool"
	"github.com/idelchi/godyl/pkg/logger"
	"github.com/idelchi/godyl/pkg/pretty"
)

// DefaultRunner is the default implementation of the Runner interface.
type DefaultRunner struct {
	cache *cache.Cache
	log   *logger.Logger
}

// NewDefaultRunner creates a new DefaultRunner instance.
func NewDefaultRunner(cache *cache.Cache, log *logger.Logger) *DefaultRunner {
	return &DefaultRunner{
		cache: cache,
		log:   log,
	}
}

// Run executes a tool operation and returns the result.
func (r *DefaultRunner) Run(ctx context.Context, t *tool.Tool, tags tags.IncludeTags, options ...RunOption) Result {
	// Apply options
	opts := &runOptions{}
	for _, opt := range options {
		opt(opts)
	}

	// Enable cache if available
	if r.cache != nil {
		t.EnableCache(r.cache)
	}

	// Log tool configuration
	r.logToolConfiguration(t)

	// Resolve the tool
	resolveResult := t.Resolve(tags, opts.resolveOptions...)

	// Convert internal result to runner.Result
	if !resolveResult.IsOK() {
		return r.convertResult(t, resolveResult)
	}

	// Check if we should skip download
	if opts.noDownload || opts.resolveOptions != nil {
		t.DisableCache()
		return r.convertResult(t, resolveResult)
	}

	// Apply SSL verification setting
	if opts.noVerifySSL {
		t.NoVerifySSL = true
	}

	// Download the tool
	downloadResult := t.Download(opts.progressTracker)

	return r.convertResult(t, downloadResult)
}

// convertResult converts an internal result.Result to a runner.Result.
func (r *DefaultRunner) convertResult(t *tool.Tool, res result.Result) Result {
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

// logToolConfiguration logs the complete tool configuration at debug level.
func (r *DefaultRunner) logToolConfiguration(tool *tool.Tool) {
	r.log.Debug("Tool:")
	r.log.Debug("-------")
	r.log.Debugf("%s", pretty.YAML(tool))
	r.log.Debug("-------")
}
