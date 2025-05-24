package download

import (
	"github.com/idelchi/godyl/internal/config/common"
	"github.com/idelchi/godyl/internal/tools/sources"
	"github.com/idelchi/godyl/internal/tools/strategy"
)

type Download struct {
	common.Tracker `json:"-"       mapstructure:"-"`
	Version        string       `json:"version" mapstructure:"version"`
	Source         sources.Type `json:"source"  mapstructure:"source"  validate:"oneof=github gitlab url go command"`
	OS             string       `json:"os"      mapstructure:"os"`
	Arch           string       `json:"arch"    mapstructure:"arch"`
	Output         string       `json:"output"  mapstructure:"output"`
	Hints          []string     `json:"hints"   mapstructure:"hints"`
	Dry            bool         `json:"dry"     mapstructure:"dry"`
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
