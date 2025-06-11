// Package dump provides the configuration for the `dump` command.
package dump

import (
	"github.com/idelchi/godyl/internal/config/common"
	"github.com/idelchi/godyl/internal/config/dump/tools"
)

// Dump holds the configuration for the `dump` command.
type Dump struct {
	// Tools contains the configuration for the `godyl dump tools` command.
	Tools tools.Tools `yaml:"tools" mapstructure:"tools" validate:"-"`

	// Tracker embed the common tracker configuration, allowing to tracker
	// whether configuration values have been explicitly set or defaulted
	common.Tracker `yaml:"-" mapstructure:"-"`
}
