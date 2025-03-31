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
	"github.com/idelchi/godyl/pkg/file"
	"github.com/idelchi/godyl/pkg/utils"
)

// Defaults holds all the configuration options for godyl, including tool-specific defaults.
type Defaults struct {
	// Inline tool-specific defaults.
	defaults tools.Defaults
}

// Get returns the Defaults struct.
func (d *Defaults) Get() tools.Defaults {
	return d.defaults
}

// Unmarshal parses the provided YAML data into the Defaults struct.
func (d *Defaults) Unmarshal(data []byte) error {
	if err := yaml.Unmarshal(data, &d.defaults); err != nil {
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
//
// TODO(Idelchi): This is not subcommand-agnostic.
func (d *Defaults) Merge(cfg config.Config) error {
	if cobraext.IsSet("hints") {
		for _, hint := range cfg.Tool.Hints {
			d.defaults.Hints.Add(match.Hint{
				Pattern: hint,
				Weight:  "1",
			})
		}
	}

	if cobraext.IsSet("output") || utils.IsZeroValue(d.defaults.Output) {
		d.defaults.Output = cfg.Tool.Output
	}

	if cobraext.IsSet("source") || utils.IsZeroValue(d.defaults.Source.Type) {
		d.defaults.Source.Type = cfg.Tool.Source
	}

	if cobraext.IsSet("strategy") || utils.IsZeroValue(d.defaults.Strategy) {
		d.defaults.Strategy = cfg.Tool.Strategy
	}

	// TODO(Idelchi): This is pretty bad.
	if cobraext.IsSet("github-token") || utils.IsZeroValue(d.defaults.Source.Github.Token) {
		if utils.IsZeroValue(cfg.Tool.Tokens.GitHub) {
			d.defaults.Source.Github.Token = cfg.Update.Tokens.GitHub
		} else {
			d.defaults.Source.Github.Token = cfg.Tool.Tokens.GitHub
		}
	}

	if cobraext.IsSet("os") || utils.IsZeroValue(d.defaults.Platform.OS) {
		if err := d.defaults.Platform.OS.Parse(cfg.Tool.OS); err != nil {
			return fmt.Errorf("parsing OS: %w", err)
		}

		d.defaults.Platform.Extension = d.defaults.Platform.Extension.Default(d.defaults.Platform.OS)
		d.defaults.Platform.Library = d.defaults.Platform.Library.Default(
			d.defaults.Platform.OS,
			d.defaults.Platform.Distribution,
		)
	}

	if cobraext.IsSet("arch") || utils.IsZeroValue(d.defaults.Platform.Architecture) {
		if err := d.defaults.Platform.Architecture.Parse(cfg.Tool.Arch); err != nil {
			return fmt.Errorf("parsing architecture: %w", err)
		}
	}

	return nil
}

// Load loads configuration defaults from a file or uses embedded defaults if not specified.
func (d *Defaults) Load(path file.File, defaults []byte) error {
	if err := d.FromFile(path.Name()); err != nil {
		if cobraext.IsSet("defaults") {
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

// Load loads the default configuration.
func Load(path file.File, embeds config.Embedded, cfg config.Config) (tools.Defaults, error) {
	// Create a new Defaults instance
	defaults := &Defaults{}

	// Load defaults from file or embedded data
	if err := defaults.Load(path, embeds.Defaults); err != nil {
		return tools.Defaults{}, err
	}

	// Apply configuration overrides
	if err := defaults.Merge(cfg); err != nil {
		return tools.Defaults{}, err
	}

	return defaults.defaults, nil
}
