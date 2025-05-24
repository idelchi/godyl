package install

import (
	"github.com/idelchi/godyl/internal/config/common"
	"github.com/idelchi/godyl/internal/tools/strategy"
)

type Install struct {
	common.Tracker `json:"-"        mapstructure:"-"`
	Strategy       strategy.Strategy `json:"strategy" mapstructure:"strategy" validate:"oneof=none sync force"`
	OS             string            `json:"os"       mapstructure:"os"`
	Arch           string            `json:"arch"     mapstructure:"arch"`
	Output         string            `json:"output"   mapstructure:"output"`
	Tags           []string          `json:"tags"     mapstructure:"tags"`
	Dry            bool              `json:"dry"      mapstructure:"dry"`
}

func (i Install) ToCommon() common.Common {
	return common.Common{
		Output:   i.Output,
		Strategy: i.Strategy,
		OS:       i.OS,
		Arch:     i.Arch,

		Tracker: i.Tracker,
	}
}
