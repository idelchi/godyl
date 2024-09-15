// Package tool provides core functionality for managing tool configurations.
package tool

import (
	"crypto/sha256"
	"fmt"

	"github.com/goccy/go-yaml/ast"

	"github.com/idelchi/godyl/internal/cache"
	"github.com/idelchi/godyl/internal/detect"
	"github.com/idelchi/godyl/internal/match"
	"github.com/idelchi/godyl/internal/tools/aliases"
	"github.com/idelchi/godyl/internal/tools/command"
	"github.com/idelchi/godyl/internal/tools/env"
	"github.com/idelchi/godyl/internal/tools/exe"
	"github.com/idelchi/godyl/internal/tools/extensions"
	"github.com/idelchi/godyl/internal/tools/fallbacks"
	"github.com/idelchi/godyl/internal/tools/inherit"
	"github.com/idelchi/godyl/internal/tools/mode"
	"github.com/idelchi/godyl/internal/tools/result"
	"github.com/idelchi/godyl/internal/tools/skip"
	"github.com/idelchi/godyl/internal/tools/sources"
	"github.com/idelchi/godyl/internal/tools/strategy"
	"github.com/idelchi/godyl/internal/tools/tags"
	"github.com/idelchi/godyl/internal/tools/values"
	"github.com/idelchi/godyl/internal/tools/version"
	"github.com/idelchi/godyl/pkg/executable"
	"github.com/idelchi/godyl/pkg/path/file"
	"github.com/idelchi/godyl/pkg/unmarshal"
)

// Tool represents a single tool configuration.
// It contains various fields that specify details such as the tool's name, version, path, execution settings,
// platform-specific settings, environment variables, and custom strategies for downloading, testing, or deploying.
type Tool struct {
	// Name of the tool, usually a short identifier or title.
	Name string `validate:"required" single:"true"`
	// Description of the tool, giving more context about its purpose.
	Description string
	// Version specifies the version of the tool.
	Version version.Version
	// Path represents the URL where the tool can be downloaded from.
	URL string
	// Output defines the output path where the tool will be installed or extracted.
	Output string
	// Exe specifies the executable details for the tool, such as patterns or names for locating the binary.
	Exe exe.Exe
	// Platform defines the platform-specific details for the tool, including OS and architecture constraints.
	Platform detect.Platform
	// Aliases represent alternative names or shortcuts for the tool.
	Aliases aliases.Aliases
	// Values contains custom values or variables used in the tool's configuration.
	Values values.Values
	// Fallbacks defines fallback configurations in case the primary configuration fails.
	Fallbacks fallbacks.Fallbacks
	// Hints provide additional matching patterns or heuristics for the tool.
	Hints *match.Hints
	// Source defines the source configuration, which determines how the tool is fetched (e.g., GitHub, local files).
	Source sources.Source
	// Commands contains a set of commands that can be executed in the context of the tool.
	Commands command.Commands
	// Tags are labels or markers that can be used to categorize or filter the tool.
	Tags tags.Tags
	// Strategy defines how the tool is deployed, fetched, or managed (e.g., download strategies, handling retries).
	Strategy strategy.Strategy
	// Extensions lists additional files or behaviors that are tied to the tool.
	Extensions *extensions.Extensions
	// Skip defines conditions under which certain steps (e.g., downloading, testing) are skipped.
	Skip skip.Skip
	// Mode defines the operating mode for the tool, potentially controlling behavior such as silent mode or verbose mode.
	Mode mode.Mode
	// Env defines the environment variables that are applied when running the tool.
	Env env.Env
	// NoVerifySSL specifies whether SSL verification should be disabled when fetching the tool.
	NoVerifySSL bool `mapstructure:"no-verify-ssl" yaml:"no-verify-ssl"`
	// NoCache disables cache interaction
	NoCache bool `mapstructure:"no-cache" yaml:"no-cache"`
	// Inherit is used to determine which default configurations the tool should inherit from.
	Inherit *inherit.Inherit
	// Cache can be carried around for various checks
	cache *cache.Cache
	// currentResult stores the latest operation result
	currentResult result.Result
	// populator stores the last successful populator
	populator sources.Populator
}

// NewEmptyTool returns an empty tool to make sure that no pointers are nil.
func NewEmptyTool() *Tool {
	return &Tool{
		Hints:      &match.Hints{},
		Extensions: &extensions.Extensions{},
		Inherit:    &inherit.Inherit{},
	}
}

// UnmarshalYAML implements custom YAML unmarshaling for Tool configuration.
// Supports both scalar values (treated as tool name) and map values.
func (t *Tool) UnmarshalYAML(node ast.Node) error {
	type raw Tool

	return unmarshal.SingleStringOrStruct(node, (*raw)(t))
}

// EnableCache sets the cache for the Tool instance.
func (t *Tool) EnableCache(cache *cache.Cache) {
	t.cache = cache
}

// EnableCache sets the cache for the Tool instance.
func (t *Tool) DisableCache() {
	t.cache = nil
}

// SetResult sets the current result of the Tool instance.
//
// TODO(Idelchi): Get rid of currentResult.
func (t *Tool) SetResult(res result.Result) {
	t.currentResult = res
}

// Result returns the current result of the Tool instance.
func (t *Tool) Result() *result.Result {
	return &t.currentResult
}

// Exists checks if the tool's executable exists in the configured output path.
// Returns true if the file exists and is a regular file.
func (t Tool) Exists() bool {
	f := file.New(t.Output, t.Exe.Name)

	return f.Exists() && f.IsFile()
}

// GetCurrentVersion attempts to retrieve the current version of the tool.
func (t Tool) GetCurrentVersion() string {
	if !t.Exists() {
		return ""
	}

	// Parse the version of the existing tool.
	exe := executable.New(t.Output, t.Exe.Name)

	// Try to get version - first from cache, then using commands
	if !t.NoCache {
		if item, err := t.cache.Get(t.ID()); err == nil {
			return item.Version.Version
		}
	}

	// No cache hit, check if we have commands to determine version
	if t.Version.Commands == nil {
		return ""
	}

	// Parse version using available commands
	parser := &executable.Parser{
		Patterns: *t.Version.Patterns,
		Commands: *t.Version.Commands,
	}

	if parsed, err := exe.Parse(parser); err != nil {
		return ""
	} else {
		return parsed
	}
}

// GetStrategy returns the tool's strategy.
func (t Tool) GetStrategy() strategy.Strategy {
	return t.Strategy
}

// GetTargetVersion returns the tool's target version.
func (t Tool) GetTargetVersion() string {
	return t.Version.Version
}

func (t Tool) GetPopulator() sources.Populator {
	return t.populator
}

// ID generates a unique identifier for the tool based on its output path and name.
func (t Tool) ID() string {
	path := file.New(t.Output, t.Name).Path()
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(path)))

	return hash
}

// AbsPath returns the absolute path of the tool's executable.
func (t Tool) AbsPath() string {
	return file.New(t.Output, t.Exe.Name).Absolute().Path()
}
