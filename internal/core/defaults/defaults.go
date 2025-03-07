package defaults

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-playground/validator/v10"

	"github.com/idelchi/godyl/internal/config"
	"github.com/idelchi/godyl/internal/tools"

	"gopkg.in/yaml.v3"
)

// DefaultsLoader is responsible for loading defaults from various sources.
type DefaultsLoader interface {
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
	return yaml.Unmarshal(data, d)
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
		return fmt.Errorf("validating Defaults: %w", err)
	}

	return nil
}

// Merge applies values from a Config object into the Defaults struct, only if corresponding values are set.
func (d *Defaults) Merge(cfg config.Config) error {
	if config.IsSet("output") {
		d.Output = cfg.Output
	}

	if config.IsSet("source") {
		d.Source.Type = cfg.Source
	}

	if config.IsSet("strategy") {
		d.Strategy = cfg.Strategy
	}

	if config.IsSet("github-token") {
		d.Source.Github.Token = cfg.Tokens.GitHub
	}

	if config.IsSet("os") {
		if err := d.Platform.OS.Parse(cfg.OS); err != nil {
			return fmt.Errorf("parsing OS: %w", err)
		}
		d.Platform.Extension = d.Platform.Extension.Default(d.Platform.OS)
		d.Platform.Library = d.Platform.Library.Default(d.Platform.OS, d.Platform.Distribution)
	}

	if config.IsSet("arch") {
		if err := d.Platform.Architecture.Parse(cfg.Arch); err != nil {
			return fmt.Errorf("parsing architecture: %w", err)
		}
	}

	if err := d.Validate(); err != nil {
		return fmt.Errorf("merging defaults: %w", err)
	}

	return nil
}

// Load loads configuration defaults from a file or uses embedded defaults if not specified.
func (d *Defaults) Load(path string, defaults []byte) error {
	if config.IsSet("defaults") {
		if err := d.FromFile(path); err != nil {
			return fmt.Errorf("loading defaults from %q: %w", path, err)
		}
	} else {
		if err := d.Default(defaults); err != nil {
			return fmt.Errorf("setting defaults: %w", err)
		}
	}

	if err := d.Initialize(); err != nil {
		return fmt.Errorf("setting tool defaults: %w", err)
	}

	return nil
}

// LoadDefaults loads the default configuration.
// This function is kept for backward compatibility.
func LoadDefaults(defaults *tools.Defaults, path string, defaultEmbedded []byte, cfg config.Config) error {
	// Create a new DefaultsManager
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

// DefaultsManager manages the loading and merging of defaults.
type DefaultsManager struct {
	defaults *Defaults
}

// NewDefaultsManager creates a new DefaultsManager.
func NewDefaultsManager() *DefaultsManager {
	return &DefaultsManager{
		defaults: NewDefaults(),
	}
}

// LoadDefaults loads defaults from a file or embedded data.
func (m *DefaultsManager) LoadDefaults(path string, defaultEmbedded []byte) error {
	if config.IsSet("defaults") {
		data, err := os.ReadFile(path)
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
func (m *DefaultsManager) ApplyConfig(cfg config.Config) error {
	// Apply configuration overrides
	if config.IsSet("output") {
		m.defaults.Output = cfg.Output
	}

	if config.IsSet("source") {
		m.defaults.Source.Type = cfg.Source
	}

	if config.IsSet("strategy") {
		m.defaults.Strategy = cfg.Strategy
	}

	if config.IsSet("github-token") {
		m.defaults.Source.Github.Token = cfg.Tokens.GitHub
	}

	if config.IsSet("os") {
		if err := m.defaults.Platform.OS.Parse(cfg.OS); err != nil {
			return fmt.Errorf("parsing OS: %w", err)
		}
		m.defaults.Platform.Extension = m.defaults.Platform.Extension.Default(m.defaults.Platform.OS)
		m.defaults.Platform.Library = m.defaults.Platform.Library.Default(m.defaults.Platform.OS, m.defaults.Platform.Distribution)
	}

	if config.IsSet("arch") {
		if err := m.defaults.Platform.Architecture.Parse(cfg.Arch); err != nil {
			return fmt.Errorf("parsing architecture: %w", err)
		}
	}

	return nil
}
