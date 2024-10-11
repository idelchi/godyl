package tools

import (
	"github.com/idelchi/godyl/internal/detect"
	stringlike "github.com/idelchi/godyl/internal/generic"
	"github.com/idelchi/godyl/internal/match"
	"github.com/idelchi/godyl/internal/tools/sources"
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
}

func (d *Defaults) Defaults() error {
	p := detect.Platform{}
	if err := p.Detect(); err != nil {
		return err
	}

	d.Platform.Merge(p)

	stringlike.SetSliceIfNil(&d.Extensions, p.CommonExtensions()...)
	// stringlike.SetSliceIfNil(&d.Exe.Patterns, "{{ .Exe.Name }}.*")

	return nil
}

func (t *Tool) ApplyDefaults(d Defaults) {
	stringlike.SetIfEmpty(&t.Output, d.Output)
	stringlike.SetIfEmpty(&t.Source.Type, d.Source.Type)
	stringlike.SetIfEmpty(&t.Source.Github.Token, d.Source.Github.Token)
	stringlike.SetIfEmpty(&t.Strategy, d.Strategy)
	stringlike.SetIfEmpty(&t.Skip.Template, "false")
	stringlike.SetSliceIfNil(&t.Exe.Patterns, d.Exe.Patterns...)
	stringlike.SetSliceIfNil(&t.Extensions, d.Extensions...)

	t.Platform.Merge(d.Platform)
	t.Hints.Add(d.Hints)
}
