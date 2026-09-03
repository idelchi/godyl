// Package download provides configuration for download operations.
package download

import (
	"github.com/idelchi/godyl/internal/config/shared"
	"github.com/idelchi/godyl/internal/tools/sources"
	"github.com/idelchi/godyl/internal/tools/strategy"
)

// Download represents configuration options for file download operations.
type Download struct {
	// Tracker records which download options were supplied explicitly.
	shared.Tracker `mapstructure:"-" yaml:"-"`

	// Version identifies the release to download.
	Version string `mapstructure:"version" yaml:"version"`
	// Source selects the release provider.
	Source sources.Type `mapstructure:"source"  validate:"oneof=github gitlab url" yaml:"source"`
	// OS selects the target operating system.
	OS string `mapstructure:"os"      yaml:"os"`
	// Arch selects the target architecture.
	Arch string `mapstructure:"arch"    yaml:"arch"`
	// Output specifies the destination path.
	Output string `mapstructure:"output"  yaml:"output"`
	// Hints help rank matching release assets.
	Hints []string `mapstructure:"hints"   yaml:"hints"`
	// Dry reports what would be downloaded without writing it.
	Dry bool `mapstructure:"dry"     yaml:"dry"`
	// Pre allows prerelease versions.
	Pre bool `mapstructure:"pre"     yaml:"pre"`
}

// ToCommon converts Download configuration to a common configuration structure.
func (d Download) ToCommon() shared.Common {
	return shared.Common{
		Output:   d.Output,
		Strategy: strategy.Force,
		Source:   d.Source,
		OS:       d.OS,
		Arch:     d.Arch,
		Hints:    d.Hints,
		Pre:      d.Pre,

		Tracker: d.Tracker,
	}
}
