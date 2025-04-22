package defaults

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/idelchi/godyl/internal/config"
	"github.com/idelchi/godyl/internal/detect"
	"github.com/idelchi/godyl/internal/match"
	"github.com/idelchi/godyl/internal/tools/tool"

	"gopkg.in/yaml.v3"
)

// Default manages the default configuration settings for a single tool.
// It is a type alias for tool.Tool, providing direct access to tool configuration fields.
type Default tool.Tool

// Unmarshal parses YAML configuration data into the Default.
// Returns an error if the YAML data is invalid or cannot be parsed.
func (d *Default) Unmarshal(data []byte) error {
	if err := yaml.Unmarshal(data, d); err != nil {
		return fmt.Errorf("unmarshalling defaults: %w", err)
	}

	return nil
}

// FromFile loads and parses a YAML configuration file into Default.
// Returns an error if the file cannot be read or contains invalid YAML.
func (d *Default) FromFile(path string) error {
	data, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return fmt.Errorf("reading file %q: %w", path, err)
	}

	return d.Unmarshal(data)
}

// Validate performs structural validation of the Default configuration.
// Ensures all required fields are properly set and contain valid values.
func (d *Default) Validate() error {
	if d.Inherit != "" {
		return fmt.Errorf("inheritance is not supported in defaults")
	}

	return nil
}

// MergeWithConfig applies configuration overrides from flags and environment variables.
// Updates default values for output paths, source types, tokens, platform settings,
// and other configurable options. Returns an error if any values are invalid.
func (d *Default) MergeWithConfig(cfg config.Config) error {
	if cfg.Tool.IsSet("output") {
		d.Output = cfg.Tool.Output
	}

	if cfg.Tool.IsSet("source") {
		d.Source.Type = cfg.Tool.Source
	}

	if cfg.Tool.IsSet("strategy") {
		d.Strategy = cfg.Tool.Strategy
	}

	if cfg.Root.IsSet("github-token") {
		d.Source.GitHub.Token = cfg.Root.Tokens.GitHub
	}

	if cfg.Root.IsSet("gitlab-token") {
		d.Source.GitLab.Token = cfg.Root.Tokens.GitLab
	}

	if cfg.Root.IsSet("url-token") {
		d.Source.URL.Token.Token = cfg.Root.Tokens.URL.Token
	}

	if cfg.Root.IsSet("url-token-header") {
		d.Source.URL.Token.Header = cfg.Root.Tokens.URL.Header
	}

	if cfg.Tool.IsSet("os") {
		if err := d.Platform.OS.Parse(cfg.Tool.OS); err != nil {
			return fmt.Errorf("parsing OS: %w", err)
		}

		d.Platform.Extension = d.Platform.Extension.Default(d.Platform.OS)
		d.Platform.Library = d.Platform.Library.Default(
			d.Platform.OS,
			d.Platform.Distribution,
		)
	}

	if cfg.Tool.IsSet("arch") {
		if err := d.Platform.Architecture.Parse(cfg.Tool.Arch); err != nil {
			return fmt.Errorf("parsing architecture: %w", err)
		}
	}

	for _, hint := range cfg.Tool.Hints {
		d.Hints.Add(match.Hint{
			Pattern: hint,
			Weight:  "1",
		})
	}

	return nil
}

// Initialize detects the current platform and applies platform-specific defaults to the Default struct.
// It also sets up default extensions based on the detected platform.
func (d *Default) Initialize(platform detect.Platform) error {
	// Merge the detected platform details with the default platform settings.
	d.Platform.Merge(platform)

	return nil
}

// ToToolg converts the Default back to a tool.Tool configuration.
func (d *Default) ToTool() tool.Tool {
	return tool.Tool(*d)
}
