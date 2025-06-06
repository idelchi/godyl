// Package status provides configuration and flags for the `godyl status` command.
package status

import "github.com/idelchi/godyl/internal/config/common"

// Status represents the configuration for the `status` command.
type Status struct {
	// Tags are the tags to consider when checking the status.
	Tags []string `json:"tags" mapstructure:"tags"`

	// Tracker embed the common tracker configuration, allowing to tracker
	// whether configuration values have been explicitly set or defaulted
	common.Tracker `json:"-" mapstructure:"-"`
}

// ToCommon converts the Status configuration to a common.Common instance.
func (s Status) ToCommon() common.Common {
	return common.Common{
		Tracker: s.Tracker,
	}
}
