package download

import (
	"github.com/idelchi/godyl/internal/config/common"
	"github.com/idelchi/godyl/internal/tools/sources"
	"github.com/idelchi/godyl/internal/tools/strategy"
)

type Download struct {
	Version string       `mapstructure:"version" yaml:"version"`
	Source  sources.Type `mapstructure:"source"  yaml:"source"  validate:"oneof=github gitlab url go command"`
	OS      string       `mapstructure:"os"      yaml:"os"`
	Arch    string       `mapstructure:"arch"    yaml:"arch"`
	Output  string       `mapstructure:"output"  yaml:"output"`
	Hints   []string     `mapstructure:"hints"   yaml:"hints"`
	Dry     bool         `mapstructure:"dry"     yaml:"dry"`

	common.Tracker `mapstructure:"-" yaml:"-"`
}

func (d Download) ToCommon() common.Common {
	return common.Common{
		Output:   d.Output,
		Strategy: strategy.Strategy(strategy.Force),
		Source:   d.Source,
		OS:       d.OS,
		Arch:     d.Arch,
		Hints:    d.Hints,

		Tracker: d.Tracker,
	}
}
