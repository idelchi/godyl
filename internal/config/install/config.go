package install

import (
	"github.com/idelchi/godyl/internal/config/common"
	"github.com/idelchi/godyl/internal/tools/strategy"
)

// Package install provides the configuration for the `install` command.
type Install struct {
	// Strategy defines how the installation should be performed
	Strategy strategy.Strategy `json:"strategy" mapstructure:"strategy" validate:"oneof=none sync force"`

	// OS defines the target operating system for the installation
	OS string `json:"os" mapstructure:"os"`

	// Arch defines the target architecture for the installation
	Arch string `json:"arch" mapstructure:"arch"`

	// Output specifies the output directory for the installation
	Output string `json:"output" mapstructure:"output"`

	// Tags are used to filter the installation based on specific criteria
	Tags []string `json:"tags" mapstructure:"tags"`

	// Dry indicates whether the installation should be performed in dry-run mode
	Dry bool `json:"dry" mapstructure:"dry"`

	// Tracker embed the common tracker configuration, allowing to tracker
	// whether configuration values have been explicitly set or defaulted
	common.Tracker `json:"-" mapstructure:"-"`
}

// ToCommon converts the Install configuration to a common.Common instance.
func (i Install) ToCommon() common.Common {
	return common.Common{
		Output:   i.Output,
		Strategy: i.Strategy,
		OS:       i.OS,
		Arch:     i.Arch,

		Tracker: i.Tracker,
	}
}
