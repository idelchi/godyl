package tools

import (
	"github.com/idelchi/godyl/internal/detect"
	"github.com/idelchi/godyl/internal/detect/platform"
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

	if len(d.Extensions) == 0 {
		switch d.Platform.OS {
		case platform.Windows:
			d.Extensions = []string{
				".zip",
				".exe",
				".gz",
			}
		default:
			d.Extensions = []string{
				".gz",
				"",
			}
		}
	}

	if d.Exe.Patterns == nil {
		d.Exe.Patterns = []string{"{{ .Exe.Name }}.*"}
	}

	return nil
}

func (t *Tool) ApplyDefaults(d Defaults) {
	SetStringIfEmpty(&t.Output, d.Output)
	SetStringIfEmpty(&t.Source.Type, d.Source.Type)
	SetStringIfEmpty(&t.Source.Github.Token, d.Source.Github.Token)
	SetStringIfEmpty(&t.Strategy, d.Strategy)
	SetStringIfEmpty(&t.SkipTemplate, "false")
	SetStringSliceIfNil(&t.Exe.Patterns, d.Exe.Patterns...)
	SetStringSliceIfNil(&t.Extensions, d.Extensions...)

	t.Platform.Merge(d.Platform)
	t.Hints.Add(d.Hints)
}
