// Package status provides configuration and flags for the `godyl status` command.
package status

import "github.com/idelchi/godyl/internal/config/shared"

// Status represents the configuration for the `status` command.
type Status struct {
	// Tracker embed the common tracker configuration, allowing to tracker
	// whether configuration values have been explicitly set or defaulted
	shared.Tracker `mapstructure:"-" yaml:"-"`

	// Tags are the tags to consider when checking the status.
	Tags []string `mapstructure:"tags" yaml:"tags"`
}

// ToCommon converts the Status configuration to a shared.Common instance.
func (s Status) ToCommon() shared.Common {
	return shared.Common{
		Tracker: s.Tracker,
	}
}
