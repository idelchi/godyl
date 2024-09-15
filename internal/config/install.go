package config

import (
	"github.com/idelchi/godyl/internal/tools/strategy"
)

type Install struct {
	Strategy strategy.Strategy `validate:"oneof=none sync force"`
	OS       string
	Arch     string
	Output   string
	Tags     []string

	trackable `json:"-" mapstructure:"-" yaml:"-"`
}

func (i Install) ToCommon() Common {
	return Common{
		Output:   i.Output,
		Strategy: i.Strategy,
		// Source is not valid for install
		// Source:
		OS:   i.OS,
		Arch: i.Arch,
		// Hints is not valid for install
		// Hints:

		trackable: i.trackable,
	}
}
