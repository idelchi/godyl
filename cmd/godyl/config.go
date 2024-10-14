package main

import (
	"fmt"
	"slices"

	_ "embed"

	"github.com/go-playground/validator/v10"
	"github.com/idelchi/godyl/internal/tools"
	"github.com/idelchi/godyl/pkg/logger"
)

//go:embed config.yml
var defaultConfigFile []byte

type Update struct {
	Strategy tools.Strategy `mapstructure:"update-strategy"`
	Update   bool           `mapstructure:"update"`
}

type Tokens struct {
	GitHub string `mapstructure:"github-token"`
}

// Config holds all the configuration options for godyl.
type Flags struct {
	Output   string
	Tools    string
	Tags     []string
	Config   string `yaml:"-"`
	Update   Update `mapstructure:",squash"`
	Dry      bool
	Detect   bool
	Log      logger.Level
	Tokens   Tokens `mapstructure:",squash"`
	Mode     tools.Mode
	Source   string
	Strategy tools.Strategy
	// Show help message
	Help bool
	// Show parsed configuration
	Show bool
	// Show version information
	Version bool

	// Number of parallel downloads
	Parallel int `validate:"gte=0"`
}

// Config holds all the configuration options for godyl.
type Config struct {
	// Defaults for tools. Allows setting a default subset of values for tools
	Defaults tools.Defaults

	// Path to file to load tools from
	Tools string

	// Tags to consider when selecting tools
	Tags []string

	// Config file to load
	Config string `yaml:"-"`

	// Update the binary now
	Update Update `mapstructure:",squash"`

	Dry bool

	Detect bool

	Log logger.Level

	// Show help message
	Help bool
	// Show parsed configuration
	Show bool
	// Show version information
	Version bool

	// Number of parallel downloads
	Parallel int `validate:"gte=0"`
}

// Validate the configuration.
func (c *Config) Validate() error {
	allowedUpdateStrategies := []tools.Strategy{tools.Upgrade, tools.Force}
	if !slices.Contains(allowedUpdateStrategies, c.Update.Strategy) {
		return fmt.Errorf("invalid update strategy: %q: allowed are %v", c.Update.Strategy, allowedUpdateStrategies)
	}

	validate := validator.New()
	if err := validate.Struct(c); err != nil {
		return fmt.Errorf("validating config: %w", err)
	}
	return nil
}
