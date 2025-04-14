package tools

import (
	"fmt"

	"github.com/fatih/structs"
	"github.com/idelchi/godyl/internal/cache/cache"
	"github.com/idelchi/godyl/internal/detect"
	"github.com/idelchi/godyl/internal/match"
	"github.com/idelchi/godyl/internal/tools/sources"
	"github.com/idelchi/godyl/internal/tools/sources/command"
	"github.com/idelchi/godyl/pkg/env"
	"github.com/idelchi/godyl/pkg/unmarshal"
	"github.com/idelchi/godyl/pkg/utils"
	"gopkg.in/yaml.v3"
)

// Tool represents a single tool configuration.
// It contains various fields that specify details such as the tool's name, version, path, execution settings,
// platform-specific settings, environment variables, and custom strategies for downloading, testing, or deploying.
type Tool struct {
	// Name of the tool, usually a short identifier or title.
	Name string `validate:"required"`
	// Description of the tool, giving more context about its purpose.
	Description string
	// Version specifies the version of the tool.
	Version Version
	// Path represents the URL or file path where the tool can be fetched or downloaded from.
	Path string
	// Output defines the output path where the tool will be installed or extracted.
	Output string
	// Exe specifies the executable details for the tool, such as patterns or names for locating the binary.
	Exe Exe
	// Platform defines the platform-specific details for the tool, including OS and architecture constraints.
	Platform detect.Platform
	// Aliases represent alternative names or shortcuts for the tool.
	Aliases Aliases
	// Values contains custom values or variables used in the tool's configuration.
	Values Values
	// Fallbacks defines fallback configurations in case the primary configuration fails.
	Fallbacks Fallbacks
	// Hints provide additional matching patterns or heuristics for the tool.
	Hints match.Hints
	// Source defines the source configuration, which determines how the tool is fetched (e.g., GitHub, local files).
	Source sources.Source
	// Tags are labels or markers that can be used to categorize or filter the tool.
	Tags Tags
	// Strategy defines how the tool is deployed, fetched, or managed (e.g., download strategies, handling retries).
	Strategy Strategy
	// Extensions lists additional files or behaviors that are tied to the tool.
	Extensions Extensions
	// Skip defines conditions under which certain steps (e.g., downloading, testing) are skipped.
	Skip Skip

	Commands command.Commands
	// Mode defines the operating mode for the tool, potentially controlling behavior such as silent mode or verbose mode.
	Mode Mode
	// Settings contains custom settings or options that modify the behavior of the tool.
	Settings Settings
	// Env defines the environment variables that are applied when running the tool.
	Env env.Env
	// Check defines a set of instructions for verifying the tool's integrity or functionality.
	Check Checker
	// NoVerifySSL specifies whether SSL verification should be disabled when fetching the tool.
	NoVerifySSL bool `yaml:"no_verify_ssl"`

	// NoCache disables cache interaction
	NoCache bool `mapstructure:"no-cache"`

	// Cache can be carried around for various checks
	cache *cache.Cache
}

// UnmarshalYAML implements custom YAML unmarshaling for Tool configuration.
// Supports both scalar values (treated as tool name) and map values with field validation.
// Ensures only known fields are present in the YAML configuration.
func (t *Tool) UnmarshalYAML(value *yaml.Node) error {
	// If it's a scalar (e.g., just the name), handle it directly by assigning it to the Name field.
	if value.Kind == yaml.ScalarNode {
		t.Name = value.Value

		return nil
	}

	type rawTool Tool

	// Use custom unmarshal logic with KnownFields check to ensure field validation.

	return unmarshal.DecodeWithOptionalKnownFields(value, (*rawTool)(t), true, structs.New(t).Name())
}

// ApplyDefaults merges default configuration values into the Tool.
// Applies defaults for values, environment variables, platform settings,
// and hints without overwriting existing non-zero values.
func (t *Tool) ApplyDefaults(d Defaults) {
	utils.DeepMergeMapsWithoutOverwrite(t.Values, d.Values)
	t.Env.Merge(d.Env)

	// Apply platform-specific defaults and hints.
	t.Platform.Merge(d.Platform)
	t.Hints.Append(d.Hints)
}

// Cache sets the cache for the Tool instance.
func (t *Tool) Cache(cache *cache.Cache) {
	t.cache = cache
}

// NewTool creates a new Tool instance with the provided default configuration.
// Deep copies mutable fields to prevent sharing between instances.
func NewTool(d Defaults) (*Tool, error) {
	d, err := utils.DeepCopy(d)
	if err != nil {
		return nil, fmt.Errorf("copying defaults: %w", err)
	}

	tool := &Tool{
		Output:     d.Output,
		Source:     d.Source,
		Strategy:   d.Strategy,
		Mode:       d.Mode,
		Exe:        d.Exe,
		Extensions: d.Extensions,
		Version:    d.Version,
		Values:     d.Values,
	}

	return tool, nil
}
