package common

import (
	"github.com/idelchi/godyl/internal/tools/sources"
	"github.com/idelchi/godyl/internal/tools/strategy"
)

type Common struct {
	Tracker  `mapstructure:"-" yaml:"-"`
	Output   string
	Strategy strategy.Strategy
	Source   sources.Type
	OS       string
	Arch     string
	Hints    []string
}
