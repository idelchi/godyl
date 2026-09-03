// Package install provides configuration for the install command.
package install

import (
	"github.com/idelchi/godyl/internal/config/shared"
	"github.com/idelchi/godyl/internal/tools/sources"
	"github.com/idelchi/godyl/internal/tools/strategy"
)

// Install provides the configuration for the `install` command.
type Install struct {
	// Tracker records which install options were supplied explicitly.
	shared.Tracker `mapstructure:"-" yaml:"-"`

	// Strategy defines how the installation should be performed.
	Strategy strategy.Strategy `mapstructure:"strategy" validate:"oneof=none sync existing force" yaml:"strategy"`

	// OS defines the target operating system for the installation.
	OS string `mapstructure:"os" yaml:"os"`

	// Arch defines the target architecture for the installation.
	Arch string `mapstructure:"arch" yaml:"arch"`

	// Output specifies the output directory for the installation.
	Output string `mapstructure:"output" yaml:"output"`

	// Tags filter the tools selected for installation.
	Tags []string `mapstructure:"tags" yaml:"tags"`

	// Dry reports planned installation work without applying it.
	Dry bool `mapstructure:"dry" yaml:"dry"`

	// Pre allows prerelease versions.
	Pre bool `mapstructure:"pre" yaml:"pre"`

	// Source restricts installations to a release provider.
	Source sources.Type `mapstructure:"source" validate:"oneof=github gitlab url none go" yaml:"source"`
}

// ToCommon converts the Install configuration to a shared.Common instance.
func (i Install) ToCommon() shared.Common {
	return shared.Common{
		Output:   i.Output,
		Strategy: i.Strategy,
		OS:       i.OS,
		Arch:     i.Arch,
		Source:   i.Source,
		Pre:      i.Pre,

		Tracker: i.Tracker,
	}
}
