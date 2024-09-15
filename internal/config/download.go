package config

import (
	"github.com/idelchi/godyl/internal/tools/sources"
	"github.com/idelchi/godyl/internal/tools/strategy"
)

type Download struct {
	Version string
	Source  sources.Type `validate:"oneof=github gitlab url go command"`
	OS      string
	Arch    string
	Output  string
	Hints   []string

	trackable `json:"-" mapstructure:"-" yaml:"-"`
}

func (d Download) ToCommon() Common {
	return Common{
		Output:   d.Output,
		Strategy: strategy.Strategy(strategy.Force),
		Source:   d.Source,
		OS:       d.OS,
		Arch:     d.Arch,
		Hints:    d.Hints,

		trackable: d.trackable,
	}
}
