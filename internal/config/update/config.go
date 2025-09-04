// Package update provides configuration and flags for the `godyl update` command.
package update

import (
	"github.com/idelchi/godyl/internal/config/shared"
	"github.com/idelchi/godyl/internal/tools/strategy"
)

// Update represents the configuration for the `update` command.
type Update struct {
	// Tracker embed the common tracker configuration, allowing to tracker
	// whether configuration values have been explicitly set or defaulted
	shared.Tracker `mapstructure:"-" yaml:"-"`

	// Version is the version to update to
	Version string `mapstructure:"version" yaml:"version"`

	// Pre indicates whether to allow pre-release versions
	Pre bool `mapstructure:"pre" yaml:"pre"`

	// Check indicates whether to only check for updates without applying them
	Check bool `mapstructure:"check" yaml:"check"`

	// Cleanup indicates whether to clean up old versions after updating
	Cleanup bool `mapstructure:"cleanup" yaml:"cleanup"`

	// Force indicates whether to force the update, ignoring any checks
	Force bool `mapstructure:"force" yaml:"force"`
}

// ToCommon converts the Update configuration to a shared.Common instance.
func (u Update) ToCommon() shared.Common {
	s := strategy.Sync

	if u.Force {
		s = strategy.Force
	}

	return shared.Common{
		Strategy: s,

		Tracker: u.Tracker,
	}
}
