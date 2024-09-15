// Package download provides configuration for download operations.
package download

import (
	"github.com/idelchi/godyl/internal/config/shared"
	"github.com/idelchi/godyl/internal/tools/sources"
	"github.com/idelchi/godyl/internal/tools/strategy"
)

// Download represents configuration options for file download operations.
type Download struct {
	shared.Tracker `mapstructure:"-" yaml:"-"`

	Version string       `mapstructure:"version" yaml:"version"`
	Source  sources.Type `mapstructure:"source"  validate:"oneof=github gitlab url" yaml:"source"`
	OS      string       `mapstructure:"os"      yaml:"os"`
	Arch    string       `mapstructure:"arch"    yaml:"arch"`
	Output  string       `mapstructure:"output"  yaml:"output"`
	Hints   []string     `mapstructure:"hints"   yaml:"hints"`
	Dry     bool         `mapstructure:"dry"     yaml:"dry"`
	Pre     bool         `mapstructure:"pre"     yaml:"pre"`
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
