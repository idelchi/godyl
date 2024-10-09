package tools

import (
	"github.com/idelchi/godyl/internal/detect"
	"github.com/idelchi/godyl/internal/detect/platform"
	"github.com/idelchi/godyl/internal/match"
	"github.com/idelchi/godyl/internal/tools/sources"
)

type Defaults struct {
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

	return nil
}

func (t *Tool) ApplyDefaults(d Defaults) {
	// Apply default for Output if empty
	if t.Output == "" {
		t.Output = d.Output
	}

	// Apply default for Source
	if t.Source.Type == "" {
		t.Source = d.Source
	}
	if t.Source.Github.Token == "" {
		t.Source.Github.Token = d.Source.Github.Token
	}

	// Apply default for Strategy if empty
	if t.Strategy == "" {
		t.Strategy = d.Strategy
	}

	t.Platform.Merge(d.Platform)

	if t.Extensions == nil {
		t.Extensions = d.Extensions
	}

	if t.SkipTemplate == "" {
		t.SkipTemplate = "false"
	}

	t.Hints.Add(d.Hints)
}
