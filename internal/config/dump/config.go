// Package dump provides the configuration for the `dump` command.
package dump

import (
	"github.com/idelchi/godyl/internal/config/dump/tools"
	"github.com/idelchi/godyl/internal/config/shared"
)

// Dump holds the configuration for the `dump` command.
type Dump struct {
	// Tracker embed the common tracker configuration, allowing to tracker
	// whether configuration values have been explicitly set or defaulted
	shared.Tracker `mapstructure:"-" yaml:"-"`

	// Tools contains the configuration for the `godyl dump tools` command.
	Tools tools.Tools `mapstructure:"tools" validate:"-" yaml:"tools"`
}
