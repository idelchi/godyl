package config

import (
	"fmt"

	"github.com/idelchi/godyl/pkg/file"
	"github.com/idelchi/godyl/pkg/validate"
)

// Root holds the root configuration options.
type Root struct {
	// Run without making any changes
	Dry bool

	// Log level (DEBUG, INFO, WARN, ERROR, SILENT)
	Log string `validate:"oneof=DEBUG INFO WARN ERROR SILENT"`

	// Path to .env file
	EnvFile file.File `mapstructure:"env-file"`

	// Path to defaults file
	Defaults file.File

	// Tokens for authentication
	Tokens Tokens `mapstructure:",squash"`

	// Viper instance
	viperable `mapstructure:"-" yaml:"-" json:"-"`
}

// Tokens holds the configuration options for authentication tokens.
type Tokens struct {
	// GitHub token for authentication
	GitHub string `mapstructure:"github-token" mask:"fixed"`
}

// Validate checks the configuration for errors.
func (r *Root) Validate() error {
	if r.IsSet("defaults") && !r.Defaults.Expanded().Exists() {
		return fmt.Errorf("%w: defaults file %q does not exist", ErrUsage, r.Defaults.Expanded())
	}

	if r.IsSet("env-file") && !r.EnvFile.Expanded().Exists() {
		return fmt.Errorf("%w: env-file file %q does not exist", ErrUsage, r.EnvFile.Expanded())
	}

	return validate.Validate(r)
}
