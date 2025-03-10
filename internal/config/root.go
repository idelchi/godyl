package config

import (
	"fmt"

	"github.com/idelchi/godyl/pkg/cobraext"
	"github.com/idelchi/godyl/pkg/file"
	"github.com/idelchi/godyl/pkg/validate"
)

// Root holds the root configuration options.
type Root struct {
	// Show enables output display
	Show bool

	// Run without making any changes
	Dry bool

	// Log level (DEBUG, INFO, WARN, ERROR, SILENT)
	Log string `validate:"oneof=DEBUG INFO WARN ERROR SILENT"`

	// Path to .env file
	DotEnv file.File `mapstructure:"env-file"`

	// Path to defaults file
	Defaults file.File
}

// Validate checks the configuration for errors.
func (c *Root) Validate() error {
	if cobraext.IsSet("defaults") && !c.Defaults.Exists() {
		return fmt.Errorf("%w: defaults file %q does not exist", ErrUsage, c.Defaults)
	}

	if cobraext.IsSet("env-file") && !c.DotEnv.Exists() {
		return fmt.Errorf("%w: env-file file %q does not exist", ErrUsage, c.DotEnv)
	}

	return validate.Validate(c)
}
