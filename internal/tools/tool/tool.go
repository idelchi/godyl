// Package tool provides core functionality for managing tool configurations.
package tool

import (
	"crypto/sha256"
	"fmt"
	"strings"

	"github.com/goccy/go-yaml/ast"

	"github.com/idelchi/godyl/internal/cache"
	"github.com/idelchi/godyl/internal/detect"
	"github.com/idelchi/godyl/internal/tools/aliases"
	"github.com/idelchi/godyl/internal/tools/command"
	"github.com/idelchi/godyl/internal/tools/exe"
	"github.com/idelchi/godyl/internal/tools/fallbacks"
	"github.com/idelchi/godyl/internal/tools/hints"
	"github.com/idelchi/godyl/internal/tools/inherit"
	"github.com/idelchi/godyl/internal/tools/mode"
	"github.com/idelchi/godyl/internal/tools/result"
	"github.com/idelchi/godyl/internal/tools/skip"
	"github.com/idelchi/godyl/internal/tools/sources"
	"github.com/idelchi/godyl/internal/tools/strategy"
	"github.com/idelchi/godyl/internal/tools/tags"
	"github.com/idelchi/godyl/internal/tools/values"
	"github.com/idelchi/godyl/internal/tools/version"
	"github.com/idelchi/godyl/pkg/env"
	"github.com/idelchi/godyl/pkg/executable"
	"github.com/idelchi/godyl/pkg/path/file"
	"github.com/idelchi/godyl/pkg/unmarshal"
)

// Tool represents a single tool configuration.
// It contains various fields that specify details such as the tool's name, version, path, execution settings,
// platform-specific settings, environment variables, and custom strategies for downloading, testing, or deploying.
//
//nolint:lll 	// Demanded by the formatter
type Tool struct {
	// Name of the tool, usually a short identifier or title.
	Name string `json:"name"          mapstructure:"name"          single:"true" validate:"required" yaml:"name"`
	// Description of the tool, giving more context about its purpose.
	Description string `json:"description"   mapstructure:"description"                                     yaml:"description"`
	// Version specifies the version of the tool.
	Version version.Version `json:"version"       mapstructure:"version"                                         yaml:"version"`
	// Path represents the URL where the tool can be downloaded from.
	URL string `json:"url"           mapstructure:"url"                                             yaml:"url"`
	// Output defines the output path where the tool will be installed or extracted.
	Output string `json:"output"        mapstructure:"output"                                          yaml:"output"`
	// Exe specifies the executable details for the tool, such as patterns or names for locating the binary.
	Exe exe.Exe `json:"exe"           mapstructure:"exe"                                             yaml:"exe"`
	// Platform defines the platform-specific details for the tool, including OS and architecture constraints.
	Platform detect.Platform `json:"platform"      mapstructure:"platform"                                        yaml:"platform"`
	// Aliases represent alternative names or shortcuts for the tool.
	Aliases aliases.Aliases `json:"aliases"       mapstructure:"aliases"                                         yaml:"aliases"`
	// Values contains custom values or variables used in the tool's configuration.
	Values values.Values `json:"values"        mapstructure:"values"                                          yaml:"values"`
	// Fallbacks defines fallback configurations in case the primary configuration fails.
	Fallbacks fallbacks.Fallbacks `json:"fallbacks"     mapstructure:"fallbacks"                                       yaml:"fallbacks"`
	// Hints provide additional matching patterns or heuristics for the tool.
	Hints *hints.Hints `json:"hints"         mapstructure:"hints"                                           yaml:"hints"`
	// Source defines the source configuration, which determines how the tool is fetched (e.g., GitHub, local files).
	Source sources.Source `json:"source"        mapstructure:"source"                                          yaml:"source"`
	// Commands contains a set of commands that can be executed in the context of the tool.
	Commands command.Commands `json:"commands"      mapstructure:"commands"                                        yaml:"commands"`
	// Tags are labels or markers that can be used to categorize or filter the tool.
	Tags tags.Tags `json:"tags"          mapstructure:"tags"                                            yaml:"tags"`
	// Strategy defines how the tool is deployed, fetched, or managed (e.g., download strategies, handling retries).
	Strategy strategy.Strategy `json:"strategy"      mapstructure:"strategy"                                        yaml:"strategy"`
	// Skip defines conditions under which certain steps (e.g., downloading, testing) are skipped.
	Skip skip.Skip `json:"skip"          mapstructure:"skip"                                            yaml:"skip"`
	Mode mode.Mode `json:"mode"          mapstructure:"mode"                                            yaml:"mode"`
	// Env defines the environment variables that are applied when running the tool.
	Env env.Env `json:"env"           mapstructure:"env"                                             yaml:"env"`
	// NoVerifySSL specifies whether SSL verification should be disabled when fetching the tool.
	NoVerifySSL bool `json:"no-verify-ssl" mapstructure:"no-verify-ssl"                                   yaml:"no-verify-ssl"`
	// NoCache disables cache interaction
	NoCache bool `json:"no-cache"      mapstructure:"no-cache"                                        yaml:"no-cache"`
	// Inherit is used to determine which default configurations the tool should inherit from.
	Inherit *inherit.Inherit `json:"inherit"       mapstructure:"inherit"                                         yaml:"inherit"`
	// Cache can be carried around for various checks
	cache *cache.Cache `json:"-"`
	// currentResult stores the latest operation result
	currentResult result.Result `json:"-"`
	// populator stores the last successful populator
	populator sources.Populator `json:"-"`
}

// NewEmptyTool returns an empty tool to make sure that no pointers are nil.
func NewEmptyTool() *Tool {
	return &Tool{
		Hints:   &hints.Hints{},
		Inherit: &inherit.Inherit{},
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

// DisableCache removes the cache from the Tool instance.
func (t *Tool) DisableCache() {
	t.cache = nil
}

// SetResult sets the current result of the Tool instance.
//
// TODO(Idelchi): Get rid of currentResult. //nolint:godox // TODO comment provides valuable context for future
// development.
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
	name := t.Exe.Name
	// Append platform-specific file extension to the executable name.
	if !strings.HasSuffix(t.Exe.Name, t.Platform.Extension.String()) && !file.File(t.Exe.Name).HasExtension() {
		name += t.Platform.Extension.String()
	}

	f := file.New(t.Output, name)

	return f.Exists() && f.IsFile()
}

// Debug prints debug information for the tool if the tool name matches the specified tool.
func (t Tool) Debug(tool, s string) {
	if t.Name == tool {
		fmt.Println(tool + ": " + s) //nolint:forbidigo // Tool output intentionally uses fmt.Println for user feedback
	}
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
			return item[0].Version.Version
		}
	}

	// No cache hit, check if we have commands to determine version
	if t.Version.Commands == nil || len(*t.Version.Commands) == 0 {
		return ""
	}

	// Parse version using available commands
	parser := &executable.Parser{
		Patterns: *t.Version.Patterns,
		Commands: *t.Version.Commands,
	}

	parsed, err := exe.Parse(parser)
	if err != nil {
		return ""
	}

	return parsed
}

// GetStrategy returns the tool's strategy.
func (t Tool) GetStrategy() strategy.Strategy {
	return t.Strategy
}

// GetTargetVersion returns the tool's target version.
func (t Tool) GetTargetVersion() string {
	return t.Version.Version
}

// GetPopulator returns the last successful populator used by the tool.
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
