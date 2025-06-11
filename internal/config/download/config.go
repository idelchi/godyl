package download

import (
	"github.com/idelchi/godyl/internal/config/common"
	"github.com/idelchi/godyl/internal/tools/sources"
	"github.com/idelchi/godyl/internal/tools/strategy"
)

type Download struct {
	Version string       `yaml:"version" mapstructure:"version"`
	Source  sources.Type `yaml:"source"  mapstructure:"source"  validate:"oneof=github gitlab url go command"`
	OS      string       `yaml:"os"      mapstructure:"os"`
	Arch    string       `yaml:"arch"    mapstructure:"arch"`
	Output  string       `yaml:"output"  mapstructure:"output"`
	Hints   []string     `yaml:"hints"   mapstructure:"hints"`
	Dry     bool         `yaml:"dry"     mapstructure:"dry"`

	common.Tracker `yaml:"-" mapstructure:"-"`
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
