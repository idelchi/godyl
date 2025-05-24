// Package runner handles the execution of tool operations.
package runner

import (
	"context"

	"github.com/hashicorp/go-getter/v2"

	"github.com/idelchi/godyl/internal/tools/tags"
	"github.com/idelchi/godyl/internal/tools/tool"
)

// Result represents the outcome of a tool operation.
type Result struct {
	Tool     *tool.Tool
	Status   Status
	Message  string
	Error    error
	Metadata map[string]any
}

// Status represents the possible states of a tool operation.
type Status int

const (
	// StatusOK indicates a successful operation.
	StatusOK Status = iota
	// StatusSkipped indicates the operation was skipped.
	StatusSkipped
	// StatusFailed indicates the operation failed.
	StatusFailed
)

// Runner defines the interface for executing tool operations.
type Runner interface {
	// Run executes a tool operation and returns the result.
	Run(ctx context.Context, tool *tool.Tool, tags tags.IncludeTags, options ...RunOption) Result
}

// RunOption configures how a tool operation is executed.
type RunOption func(*runOptions)

// runOptions holds configuration for a tool run.
type runOptions struct {
	progressTracker getter.ProgressTracker
	noDownload      bool
	noVerifySSL     bool
	resolveOptions  []tool.ResolveOption
}

// WithProgressTracker sets the progress tracker for downloads.
func WithProgressTracker(tracker getter.ProgressTracker) RunOption {
	return func(opts *runOptions) {
		opts.progressTracker = tracker
	}
}

// WithNoDownload skips the download phase.
func WithNoDownload() RunOption {
	return func(opts *runOptions) {
		opts.noDownload = true
	}
}

// WithNoVerifySSL disables SSL verification.
func WithNoVerifySSL() RunOption {
	return func(opts *runOptions) {
		opts.noVerifySSL = true
	}
}

// WithResolveOptions passes options to the tool's Resolve method.
func WithResolveOptions(resolveOpts ...tool.ResolveOption) RunOption {
	return func(opts *runOptions) {
		opts.resolveOptions = resolveOpts
	}
}
