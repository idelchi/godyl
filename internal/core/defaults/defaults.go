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
	"github.com/idelchi/godyl/pkg/cobraext"
	"github.com/idelchi/godyl/pkg/utils"
)

// Defaults holds all the configuration options for godyl, including tool-specific defaults.
type Defaults struct {
	// Inline tool-specific defaults.
	tools.Defaults `yaml:",inline"`
}

// NewDefaults creates a new Defaults instance.
func NewDefaults() *Defaults {
	return &Defaults{}
}

// Unmarshal parses the provided YAML data into the Defaults struct.
func (d *Defaults) Unmarshal(data []byte) error {
	if err := yaml.Unmarshal(data, d); err != nil {
		return fmt.Errorf("unmarshalling defaults: %w", err)
	}

	return nil
}

// FromFile reads and parses a YAML file from the given path into the Defaults struct.
func (d *Defaults) FromFile(path string) error {
	data, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return fmt.Errorf("reading file %q: %w", path, err)
	}

	return d.Unmarshal(data)
}

// Default loads the embedded default YAML configuration.
func (d *Defaults) Default(defaults []byte) error {
	return d.Unmarshal(defaults)
}

// Validate checks the Defaults struct to ensure all required fields are properly set.
func (d *Defaults) Validate() error {
	validate := validator.New()
	if err := validate.Struct(d); err != nil {
		return fmt.Errorf("validating defaults: %w", err)
	}

	return nil
}

// Merge applies configuration overrides to the defaults.
func (d *Defaults) Merge(cfg config.Config) error {
	if cobraext.IsSet("hints") {
		for _, hint := range cfg.Tool.Hints {
			d.Hints.Add(match.Hint{
				Pattern: hint,
				Weight:  "1",
			})
		}
	}

	if cobraext.IsSet("output") || utils.IsZeroValue(d.Output) {
		d.Output = cfg.Tool.Output
	}

	if cobraext.IsSet("source") || utils.IsZeroValue(d.Source.Type) {
		d.Source.Type = cfg.Tool.Source
	}

	if cobraext.IsSet("strategy") || utils.IsZeroValue(d.Strategy) {
		d.Strategy = cfg.Tool.Strategy
	}

	if cobraext.IsSet("github-token") || utils.IsZeroValue(d.Source.Github.Token) {
		d.Source.Github.Token = cfg.Tool.Tokens.GitHub
	}

	if cobraext.IsSet("os") || utils.IsZeroValue(d.Platform.OS) {
		if err := d.Platform.OS.Parse(cfg.Tool.OS); err != nil {
			return fmt.Errorf("parsing OS: %w", err)
		}

		d.Platform.Extension = d.Platform.Extension.Default(d.Platform.OS)
		d.Platform.Library = d.Platform.Library.Default(
			d.Platform.OS,
			d.Platform.Distribution,
		)
	}

	if cobraext.IsSet("arch") || utils.IsZeroValue(d.Platform.Architecture) {
		if err := d.Platform.Architecture.Parse(cfg.Tool.Arch); err != nil {
			return fmt.Errorf("parsing architecture: %w", err)
		}
	}

	return nil
}

// Load loads configuration defaults from a file or uses embedded defaults if not specified.
func (d *Defaults) Load(path string, defaults []byte) error {
	if cobraext.IsSet("defaults") {
		if err := d.FromFile(path); err != nil {
			return fmt.Errorf("loading defaults from %q: %w", path, err)
		}
	} else {
		if err := d.Default(defaults); err != nil {
			return fmt.Errorf("setting defaults: %w", err)
		}
	}

	if err := d.Initialize(); err != nil {
		return fmt.Errorf("initializing defaults: %w", err)
	}

	return nil
}

// LoadDefaults loads the default configuration.
// This function is kept for backward compatibility.
func LoadDefaults(defaults *tools.Defaults, path string, embeds config.Embedded, cfg config.Config) error {
	// Create a new Defaults instance
	d := NewDefaults()

	// Load defaults from file or embedded data
	if err := d.Load(path, embeds.Defaults); err != nil {
		return err
	}

	// Apply configuration overrides
	if err := d.Merge(cfg); err != nil {
		return err
	}

	// Copy the loaded defaults to the provided defaults struct
	*defaults = d.Defaults

	return nil
}
