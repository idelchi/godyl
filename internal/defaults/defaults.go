// Package defaults provides functionality for managing default values and configurations.
package defaults

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/idelchi/godyl/internal/config"
	"github.com/idelchi/godyl/internal/detect"
	"github.com/idelchi/godyl/pkg/path/file"

	"gopkg.in/yaml.v3"
)

// Defaults manages a collection of default configurations for various tools.
type Defaults map[string]*Default

// NewDefaults creates a new Defaults instance.
func NewDefaults() *Defaults {
	return &Defaults{}
}

// Unmarshal parses YAML configuration data into the Defaults map.
// Returns an error if the YAML data is invalid or cannot be parsed.
func (d *Defaults) Unmarshal(data []byte) error {
	if err := yaml.Unmarshal(data, d); err != nil {
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

// Validate performs structural validation of the Defaults configuration.
// Ensures all required fields are properly set and contain valid values.
func (d *Defaults) Validate() error {
	var errs []error

	// Validate each individual Default configuration
	for key, def := range *d {
		if def == nil {
			continue
		}

		if err := def.Validate(); err != nil {
			errs = append(errs, fmt.Errorf("%q: %w", key, err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("\n%w", errors.Join(errs...))
	}

	return nil
}

// MergeWithConfig applies configuration overrides from flags and environment variables.
// Updates default values for output paths, source types, tokens, platform settings,
// and other configurable options. Returns an error if any values are invalid.
func (d *Defaults) MergeWithConfig(cfg config.Config) error {
	// For each default, merge the configuration
	for key, def := range *d {
		if def == nil {
			continue
		}

		if err := def.MergeWithConfig(cfg); err != nil {
			return fmt.Errorf("merging defaults for %q: %w", key, err)
		}
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
			if err := d.Unmarshal(defaults); err != nil {
				return fmt.Errorf("setting defaults: %w", err)
			}
		}
	}

	if err := d.Initialize(); err != nil {
		return fmt.Errorf("initializing defaults: %w", err)
	}

	return nil
}

// Initialize detects the current platform and applies platform-specific defaults to all Default configurations.
func (d *Defaults) Initialize() error {
	// Detect the current platform (e.g., OS, architecture).
	platform := detect.Platform{}
	if err := platform.Detect(); err != nil {
		return err
	}

	// For each default, initialize with platform configuration
	for _, def := range *d {
		if def == nil {
			continue
		}

		if err := def.Initialize(platform); err != nil {
			return err
		}
	}

	return nil
}

// GetDefault returns a specific Default configuration by name.
// If the configuration doesn't exist, it returns nil.
func (d *Defaults) GetDefault(name string) *Default {
	if def, ok := (*d)[name]; ok {
		return def
	}

	return nil
}

// Load creates and configures a new Defaults instance.
// Handles loading from files or embedded data, applies configuration
// overrides, and returns the Defaults object and an error if any.
func Load(path file.File, embeds config.Embedded, cfg config.Config) (*Defaults, error) {
	// Create a new Defaults instance
	defaults := NewDefaults()

	// Load defaults from file or embedded data
	if err := defaults.Load(path, embeds.Defaults, cfg.Root.IsSet("defaults")); err != nil {
		return nil, err
	}

	// Apply configuration overrides
	if err := defaults.MergeWithConfig(cfg); err != nil {
		return nil, err
	}

	return defaults, nil
}
