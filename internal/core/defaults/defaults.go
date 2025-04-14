// Package defaults provides functionality for managing default values and configurations.
package defaults

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"

	"github.com/idelchi/godyl/internal/config"
	"github.com/idelchi/godyl/internal/match"
	"github.com/idelchi/godyl/internal/tools"
	"github.com/idelchi/godyl/internal/tools/sources"
	"github.com/idelchi/godyl/internal/tools/sources/github"
	"github.com/idelchi/godyl/internal/tools/sources/gitlab"
	"github.com/idelchi/godyl/internal/tools/sources/url"

	"github.com/idelchi/godyl/pkg/path/file"
)

// Defaults manages the default configuration settings for the application.
// It encapsulates tool-specific defaults and provides methods for loading,
// merging, and validating configuration values.
type Defaults struct {
	// Inline tool-specific defaults.
	defaults tools.Defaults
}

// New creates a new Defaults instance with the provided configuration.
// Initializes default values from the config struct, handling platform-specific
// settings and token configurations. Command-line flags and environment
// variables are merged into the defaults struct.
func New(cfg config.Config) *Defaults {
	defaults := &Defaults{
		defaults: tools.Defaults{
			Output: cfg.Tool.Output,
			Source: sources.Source{
				Type: cfg.Tool.Source,
				GitHub: github.GitHub{
					Token: cfg.Root.Tokens.GitHub,
				},
				URL: url.URL{
					Token: cfg.Root.Tokens.URL,
				},
				GitLab: gitlab.GitLab{
					Token: cfg.Root.Tokens.GitLab,
				},
			},
			Strategy: cfg.Tool.Strategy,
		},
	}

	defaults.defaults.Platform.OS.Parse(cfg.Tool.OS)
	defaults.defaults.Platform.Extension = defaults.defaults.Platform.Extension.Default(defaults.defaults.Platform.OS)
	defaults.defaults.Platform.Library = defaults.defaults.Platform.Library.Default(
		defaults.defaults.Platform.OS,
		defaults.defaults.Platform.Distribution,
	)

	defaults.defaults.Platform.Architecture.Parse(cfg.Tool.Arch)

	return defaults
}

// Get returns the underlying tools.Defaults configuration.
func (d *Defaults) Get() tools.Defaults {
	return d.defaults
}

// Unmarshal parses YAML configuration data into the underlying tools.Defaults.
// Returns an error if the YAML data is invalid or cannot be parsed.
func (d *Defaults) Unmarshal(data []byte) error {
	if err := yaml.Unmarshal(data, &d.defaults); err != nil {
		return fmt.Errorf("unmarshalling defaults: %w", err)
	}

	return nil
}

// FromFile loads and parses a YAML configuration file into Defaults.
// Returns an error if the file cannot be read or contains invalid YAML.
func (d *Defaults) FromFile(path string) error {
	data, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return fmt.Errorf("reading file %q: %w", path, err)
	}

	return d.Unmarshal(data)
}

// Default loads the embedded default YAML configuration into Defaults.
// Used when no external configuration file is provided.
func (d *Defaults) Default(defaults []byte) error {
	return d.Unmarshal(defaults)
}

// Validate performs structural validation of the Defaults configuration.
// Ensures all required fields are properly set and contain valid values.
func (d *Defaults) Validate() error {
	validate := validator.New()
	if err := validate.Struct(d); err != nil {
		return fmt.Errorf("validating defaults: %w", err)
	}

	return nil
}

// Merge applies configuration overrides from flags and environment variables.
// Updates default values for output paths, source types, tokens, platform settings,
// and other configurable options. Returns an error if any values are invalid.
func (d *Defaults) Merge(cfg config.Config) error {
	if cfg.Tool.IsSet("output") {
		d.defaults.Output = cfg.Tool.Output
	}

	if cfg.Tool.IsSet("source") {
		d.defaults.Source.Type = cfg.Tool.Source
	}

	if cfg.Tool.IsSet("strategy") {
		d.defaults.Strategy = cfg.Tool.Strategy
	}

	if cfg.Root.IsSet("github-token") {
		d.defaults.Source.GitHub.Token = cfg.Root.Tokens.GitHub
	}

	if cfg.Root.IsSet("gitlab-token") {
		d.defaults.Source.GitLab.Token = cfg.Root.Tokens.GitLab
	}

	if cfg.Root.IsSet("url-token") {
		d.defaults.Source.URL.Token.Token = cfg.Root.Tokens.URL.Token
	}

	if cfg.Root.IsSet("url-token-header") {
		d.defaults.Source.URL.Token.Header = cfg.Root.Tokens.URL.Header
	}

	if cfg.Tool.IsSet("os") {
		if err := d.defaults.Platform.OS.Parse(cfg.Tool.OS); err != nil {
			return fmt.Errorf("parsing OS: %w", err)
		}

		d.defaults.Platform.Extension = d.defaults.Platform.Extension.Default(d.defaults.Platform.OS)
		d.defaults.Platform.Library = d.defaults.Platform.Library.Default(
			d.defaults.Platform.OS,
			d.defaults.Platform.Distribution,
		)
	}

	if cfg.Tool.IsSet("arch") {
		if err := d.defaults.Platform.Architecture.Parse(cfg.Tool.Arch); err != nil {
			return fmt.Errorf("parsing architecture: %w", err)
		}
	}

	for _, hint := range cfg.Tool.Hints {
		d.defaults.Hints.Add(match.Hint{
			Pattern: hint,
			Weight:  "1",
		})
	}

	return nil
}

// Load loads configuration from a file or falls back to embedded defaults.
// Initializes the configuration after loading. Returns an error if loading
// or initialization fails.
func (d *Defaults) Load(path file.File, defaults []byte, isSet bool) error {
	if err := d.FromFile(path.Path()); err != nil {
		if isSet {
			return fmt.Errorf("loading defaults from %q: %w", path, err)
		} else {
			if err := d.Default(defaults); err != nil {
				return fmt.Errorf("setting defaults: %w", err)
			}
		}
	}

	if err := d.defaults.Initialize(); err != nil {
		return fmt.Errorf("initializing defaults: %w", err)
	}

	return nil
}

// Load creates and configures a new Defaults instance.
// Handles loading from files or embedded data, applies configuration
// overrides, and returns the final tools.Defaults configuration.
func Load(path file.File, embeds config.Embedded, cfg config.Config) (tools.Defaults, error) {
	// Create a new Defaults instance
	defaults := New(cfg)

	// Load defaults from file or embedded data
	if err := defaults.Load(path, embeds.Defaults, cfg.Root.IsSet("defaults")); err != nil {
		return tools.Defaults{}, err
	}

	// Apply configuration overrides
	if err := defaults.Merge(cfg); err != nil {
		return tools.Defaults{}, err
	}

	return defaults.defaults, nil
}
