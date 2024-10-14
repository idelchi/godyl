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

func (d *Defaults) Defaults() error {
	p := detect.Platform{}
	if err := p.Detect(); err != nil {
		return err
	}

	d.Platform.Merge(p)

	utils.SetSliceIfNil(&d.Extensions, p.CommonExtensions()...)
	// stringlike.SetSliceIfNil(&d.Exe.Patterns, "{{ .Exe.Name }}.*")

	return nil
}

func (t *Tool) ApplyDefaults(d Defaults) {
	utils.SetIfEmpty(&t.Output, d.Output)
	utils.SetIfEmpty(&t.Source.Type, d.Source.Type)
	utils.SetIfEmpty(&t.Source.Github.Token, d.Source.Github.Token)
	utils.SetIfEmpty(&t.Strategy, d.Strategy)
	utils.SetIfEmpty(&t.Skip.Template, "false")
	utils.SetIfEmpty(&t.Mode, d.Mode)
	utils.SetSliceIfNil(&t.Exe.Patterns, d.Exe.Patterns...)
	utils.SetSliceIfNil(&t.Extensions, d.Extensions...)

	t.Platform.Merge(d.Platform)
	t.Hints.Add(d.Hints)
}
