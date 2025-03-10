// Package defaults provides functionality for managing default values and configurations.
package defaults

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/idelchi/godyl/internal/config"
	"github.com/idelchi/godyl/internal/match"
	"github.com/idelchi/godyl/internal/tools"
	"github.com/idelchi/godyl/pkg/utils"

	"gopkg.in/yaml.v3"
)

// Loader is responsible for loading defaults from various sources.
type Loader interface {
	LoadFromFile(path string) error
	LoadFromBytes(data []byte) error
	Initialize() error
}

// ConfigMerger is responsible for merging configuration into defaults.
type ConfigMerger interface {
	MergeConfig(cfg config.Config) error
}

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
	// Using the yaml tag to ensure proper unmarshaling
	err := yaml.Unmarshal(
		data,
		d,
	) // nolint:musttag		// TODO(Idelchi): Not sure what is expected here, check later.
	if err != nil {
		return fmt.Errorf("unmarshalling defaults: %w", err)
	}

	return nil
}

// LoadDefaults loads the default configuration.
// This function is kept for backward compatibility.
func LoadDefaults(defaults *tools.Defaults, path string, defaultEmbedded []byte, cfg config.Config) error {
	// Create a new Manager
	manager := NewDefaultsManager()

	// Load defaults from file or embedded data
	if err := manager.LoadDefaults(path, defaultEmbedded); err != nil {
		return err
	}

	// Apply configuration overrides
	if err := manager.ApplyConfig(cfg); err != nil {
		return err
	}

	// Copy the loaded defaults to the provided defaults struct
	*defaults = manager.defaults.Defaults

	return nil
}

// Manager manages the loading and merging of defaults.
type Manager struct {
	defaults *Defaults
}

// NewDefaultsManager creates a new Manager.
func NewDefaultsManager() *Manager {
	return &Manager{
		defaults: NewDefaults(),
	}
}

// LoadDefaults loads defaults from a file or embedded data.
func (m *Manager) LoadDefaults(path string, defaultEmbedded []byte) error {
	if config.IsSet("defaults") {
		data, err := os.ReadFile(filepath.Clean(path))
		if err != nil {
			return fmt.Errorf("reading defaults file %q: %w", path, err)
		}

		if err := m.defaults.Unmarshal(data); err != nil {
			return fmt.Errorf("unmarshalling defaults: %w", err)
		}
	} else {
		if err := m.defaults.Unmarshal(defaultEmbedded); err != nil {
			return fmt.Errorf("unmarshalling embedded defaults: %w", err)
		}
	}

	if err := m.defaults.Initialize(); err != nil {
		return fmt.Errorf("initializing defaults: %w", err)
	}

	return nil
}

// ApplyConfig applies configuration overrides to the defaults.
func (m *Manager) ApplyConfig(cfg config.Config) error {
	if config.IsSet("hints") {
		for _, hint := range cfg.Tool.Hints {
			m.defaults.Hints.Add(match.Hint{
				Pattern: hint,
				Weight:  "1",
			})
		}
	}

	if config.IsSet("output") || utils.IsZeroValue(m.defaults.Output) {
		m.defaults.Output = cfg.Tool.Output
	}

	if config.IsSet("source") || utils.IsZeroValue(m.defaults.Source.Type) {
		m.defaults.Source.Type = cfg.Tool.Source
	}

	if config.IsSet("strategy") || utils.IsZeroValue(m.defaults.Strategy) {
		m.defaults.Strategy = cfg.Tool.Strategy
	}

	if config.IsSet("github-token") || utils.IsZeroValue(m.defaults.Source.Github.Token) {
		m.defaults.Source.Github.Token = cfg.Tool.Tokens.GitHub
	}

	if config.IsSet("os") || utils.IsZeroValue(m.defaults.Platform.OS) {
		if err := m.defaults.Platform.OS.Parse(cfg.Tool.OS); err != nil {
			return fmt.Errorf("parsing OS: %w", err)
		}

		m.defaults.Platform.Extension = m.defaults.Platform.Extension.Default(m.defaults.Platform.OS)
		m.defaults.Platform.Library = m.defaults.Platform.Library.Default(
			m.defaults.Platform.OS,
			m.defaults.Platform.Distribution,
		)
	}

	if config.IsSet("arch") || utils.IsZeroValue(m.defaults.Platform.Architecture) {
		if err := m.defaults.Platform.Architecture.Parse(cfg.Tool.Arch); err != nil {
			return fmt.Errorf("parsing architecture: %w", err)
		}
	}

	return nil
}
