package install

import (
	"github.com/idelchi/godyl/internal/config/common"
	"github.com/idelchi/godyl/internal/tools/sources"
	"github.com/idelchi/godyl/internal/tools/strategy"
)

// Package install provides the configuration for the `install` command.
type Install struct {
	// Strategy defines how the installation should be performed
	Strategy strategy.Strategy `mapstructure:"strategy" validate:"oneof=none sync force" yaml:"strategy"`

	// OS defines the target operating system for the installation
	OS string `mapstructure:"os" yaml:"os"`

	// Arch defines the target architecture for the installation
	Arch string `mapstructure:"arch" yaml:"arch"`

	// Output specifies the output directory for the installation
	Output string `mapstructure:"output" yaml:"output"`

	// Tags are used to filter the installation based on specific criteria
	Tags []string `mapstructure:"tags" yaml:"tags"`

	// Dry indicates whether the installation should be performed in dry-run mode
	Dry bool `mapstructure:"dry" yaml:"dry"`

	Source sources.Type `mapstructure:"source" yaml:"source" validate:"oneof=github gitlab url none go"`

	// Tracker embed the common tracker configuration, allowing to tracker
	// whether configuration values have been explicitly set or defaulted
	common.Tracker `mapstructure:"-" yaml:"-"`
}

// ToCommon converts the Install configuration to a common.Common instance.
func (i Install) ToCommon() common.Common {
	return common.Common{
		Output:   i.Output,
		Strategy: i.Strategy,
		OS:       i.OS,
		Arch:     i.Arch,
		Source:   i.Source,

		Tracker: i.Tracker,
	}
}
