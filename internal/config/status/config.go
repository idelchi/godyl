// Package status provides configuration and flags for the `godyl status` command.
package status

import "github.com/idelchi/godyl/internal/config/shared"

// Status represents the configuration for the `status` command.
type Status struct {
	// Tracker records which status options were supplied explicitly.
	shared.Tracker `mapstructure:"-" yaml:"-"`

	// Output specifies the directory searched for installed tools.
	Output string `mapstructure:"output" yaml:"output"`

	// Tags are the tags to consider when checking the status.
	Tags []string `mapstructure:"tags" yaml:"tags"`
}

// ToCommon converts the Status configuration to a shared.Common instance.
func (s Status) ToCommon() shared.Common {
	return shared.Common{
		Output:  s.Output,
		Tracker: s.Tracker,
	}
}
