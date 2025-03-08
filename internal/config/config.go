package config

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"

	"github.com/idelchi/godyl/internal/tools"
	"github.com/idelchi/godyl/internal/tools/sources"
	"github.com/idelchi/godyl/pkg/file"
)

// ErrUsage is returned when there is an error in the configuration.
var ErrUsage = errors.New("usage error")

// Update holds the configuration options for updating the built binary itself.
type Update struct {
	// Strategy to use for updating tools
	Strategy tools.Strategy `mapstructure:"strategy" validate:"oneof=none upgrade force"`
	// Update the tools
	Update bool `mapstructure:"update"`
}

// Tokens holds the configuration options for authentication tokens.
type Tokens struct {
	// GitHub token for authentication
	GitHub string `mapstructure:"github-token" mask:"fixed"`

	// GitLab token for authentication
	GitLab string `mapstructure:"gitlab-token" mask:"fixed"`

	// URL token for authentication
	URL string `mapstructure:"url-token" mask:"fixed"`
}

// Dump holds the configuration options for showing various configurations.
type Dump struct {
	// Show the parsed configuration and exit
	Config bool

	// Show the parsed environment variables and exit
	Env bool

	// Show the parsed default configuration and exit
	Defaults bool

	// Detect the platform and exit
	Platform bool

	// Show available tools
	Tools bool
}

// Config holds all the configuration options for godyl.
type Config struct {
	// Show enables output display
	Show bool

	// Show help message and exit
	Help bool

	// Show version information and exit
	Version bool

	// Path to .env file
	DotEnv file.File `mapstructure:"env-file"`

	// Run without making any changes (dry run)
	Dry bool

	// Log level (debug, info, warn, error, always, silent)
	Log string

	// Number of parallel downloads (>= 0)
	Parallel int `validate:"gte=0"`

	// Skip SSL verification
	NoVerifySSL bool `mapstructure:"no-verify-ssl"`

	// Path to defaults file
	Defaults file.File

	// Path to tools configuration file
	Tools string

	// Output path for the downloaded tools
	Output string

	// Tags to filter tools by
	Tags []string

	// Source from which to install the tools
	Source sources.Type `validate:"oneof=github url go command"`

	// Strategy to use for updating tools
	Strategy tools.Strategy `mapstructure:"strategy"`

	// Operating system to install the tools for
	OS string `mapstructure:"os"`

	// Architecture to install the tools for
	Arch string `mapstructure:"arch"`

	// Tokens for authentication
	Tokens Tokens `mapstructure:",squash"`

	// Dump various configurations
	Dump Dump

	// Update the tool itself
	Update Update `mapstructure:",squash"`
}

// Display returns the value of the Show field.
func (c *Config) Display() bool {
	return c.Show
}

// Validate checks the configuration for errors.
func (c *Config) Validate(_ any) error {
	if IsSet("defaults") && !c.Defaults.Exists() {
		return fmt.Errorf("%w: defaults file %q does not exist", ErrUsage, c.Defaults)
	}

	if IsSet("env-file") && !c.DotEnv.Exists() {
		return fmt.Errorf("%w: env-file file %q does not exist", ErrUsage, c.DotEnv)
	}

	validate := validator.New()
	if err := validate.Struct(c); err != nil {
		return fmt.Errorf("validating config: %w", err)
	}

	return nil
}

// IsSet checks if a flag is set in viper.
func IsSet(flag string) bool {
	return viper.IsSet(flag)
}

// NewConfig creates a new Config with default values.
func NewConfig() *Config {
	return &Config{
		DotEnv:   file.File(".env"),
		Log:      "info",
		Tools:    "tools.yml",
		Output:   "./bin",
		Tags:     []string{"!native"},
		Source:   sources.GITHUB,
		Strategy: tools.None,
	}
}
