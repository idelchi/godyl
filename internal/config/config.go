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

// // addToolFlags adds tool-related flags to the command.
// func addToolFlags(cmd *cobra.Command) {
// 	// Tool flags
// 	cmd.Flags().StringP("output", "o", "./bin", "Output path for the downloaded tools")
// 	cmd.Flags().StringSliceP("tags", "t", []string{"!native"}, "Tags to filter tools by. Prefix with '!' to exclude")
// 	cmd.Flags().String("source", "github", "Source from which to install the tools")
// 	cmd.Flags().String("strategy", "none", "Strategy to use for updating tools")
// 	cmd.Flags().String("github-token", os.Getenv("GODYL_GITHUB_TOKEN"), "GitHub token for authentication")
// 	cmd.Flags().String("os", "", "Operating system to install the tools for")
// 	cmd.Flags().String("arch", "", "Architecture to install the tools for")
// }

type ToolConfiguration struct {
	Output      string
	Tags        []string
	Source      sources.Type   `validate:"oneof=github url go command"`
	Strategy    tools.Strategy `validate:"oneof=none upgrade force"`
	Tokens      Tokens         `mapstructure:",squash"`
	OS          string
	Arch        string
	Parallel    int  `validate:"gte=0"`
	NoVerifySSL bool `mapstructure:"no-verify-ssl"`
}

type RootConfiguration struct {
	Dry bool
	Log string `validate:"oneof=debug info warn error always silent"`

	DotEnv   file.File `mapstructure:"env-file"`
	Defaults file.File
}

type UpdateConfiguration struct {
	Strategy    tools.Strategy `validate:"upgrade force"`
	Tokens      Tokens         `mapstructure:",squash"`
	NoVerifySSL bool           `mapstructure:"no-verify-ssl"`
}

type DumpConfiguration struct {
	Type   string `validate:"oneof=config defaults env platform tools"`
	Format string `validate:"oneof=json yaml"`
}

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
}

type Dump struct {
	Type   string
	Format string // `validate:"oneof=json yaml"`
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

	// Dump configuration information
	Dump Dump `mapstructure:",squash"`

	// Update the tool itself
	Update Update `mapstructure:",squash"`
}

// Display returns the value of the Show field.
func (c *Config) Display() bool {
	return c.Show
}

// Validate checks the configuration for errors.
func (c *Config) Validate(config any) error {
	if IsSet("defaults") && !c.Defaults.Exists() {
		return fmt.Errorf("%w: defaults file %q does not exist", ErrUsage, c.Defaults)
	}

	if IsSet("env-file") && !c.DotEnv.Exists() {
		return fmt.Errorf("%w: env-file file %q does not exist", ErrUsage, c.DotEnv)
	}

	validate := validator.New()
	if err := validate.Struct(config); err != nil {
		return fmt.Errorf("validating config: %w", err)
	}

	return nil
}

// IsSet checks if a flag is set in viper.
func IsSet(flag string) bool {
	return viper.IsSet(flag)
}
