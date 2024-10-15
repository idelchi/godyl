package tools

import (
	"github.com/idelchi/godyl/internal/detect"
	"github.com/idelchi/godyl/internal/match"
	"github.com/idelchi/godyl/internal/tools/sources"
	"github.com/idelchi/godyl/pkg/env"
	"github.com/idelchi/godyl/pkg/utils"
)

type Defaults struct {
	Exe        Exe
	Output     string
	Platform   detect.Platform
	Values     map[string]any
	Fallbacks  []string
	Hints      match.Hints
	Source     sources.Source
	Tags       Tags
	Strategy   Strategy
	Extensions []string
	Env        env.Env
	Mode       Mode
}

func (d *Defaults) Initialize() error {
	p := detect.Platform{}
	if err := p.Detect(); err != nil {
		return err
	}

	d.Platform.Merge(p)

	utils.SetSliceIfNil(&d.Extensions, p.CommonExtensions()...)

	return nil
}
