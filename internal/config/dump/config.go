// Package dump provides the configuration for the `dump` command.
package dump

import (
	"github.com/idelchi/godyl/internal/config/dump/tools"
	"github.com/idelchi/godyl/internal/config/shared"
)

// Dump holds the configuration for the `dump` command.
type Dump struct {
	// Tracker records which dump options were supplied explicitly.
	shared.Tracker `mapstructure:"-" yaml:"-"`

	// Tools contains the configuration for the `godyl dump tools` command.
	Tools tools.Tools `mapstructure:"tools" validate:"-" yaml:"tools"`
}
