package config

import (
	"github.com/idelchi/godyl/internal/tools/sources"
)

type Status struct {
	viperable `json:"-" mapstructure:"-" yaml:"-"`
	Source    sources.Type `validate:"oneof=github gitlab url go command"`
	Tags      []string
}
