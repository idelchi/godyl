// Package update provides configuration and flags for the `godyl update` command.
package update

import (
	"github.com/idelchi/godyl/internal/config/common"
	"github.com/idelchi/godyl/internal/tools/strategy"
)

// Update represents the configuration for the `update` command.
type Update struct {
	// Version is the version to update to
	Version string `json:"version" mapstructure:"version"`

	// Pre indicates whether to allow pre-release versions
	Pre bool `json:"pre" mapstructure:"pre"`

	// Check indicates whether to only check for updates without applying them
	Check bool `json:"check" mapstructure:"check"`

	// Cleanup indicates whether to clean up old versions after updating
	Cleanup bool `json:"cleanup" mapstructure:"cleanup"`

	// Force indicates whether to force the update, ignoring any checks
	Force bool `json:"force" mapstructure:"force"`

	// Tracker embed the common tracker configuration, allowing to tracker
	// whether configuration values have been explicitly set or defaulted
	common.Tracker `json:"-" mapstructure:"-"`
}

// ToCommon converts the Update configuration to a common.Common instance.
func (u Update) ToCommon() common.Common {
	s := strategy.Sync
	if u.Force {
		s = strategy.Force
	}

	return common.Common{
		Strategy: s,

		Tracker: u.Tracker,
	}
}
