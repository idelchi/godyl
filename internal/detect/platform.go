package detect

import (
	"github.com/idelchi/godyl/internal/detect/platform"
	stringlike "github.com/idelchi/godyl/internal/generic"
)

type Info map[string]string

type Platform struct {
	OS           platform.OS
	Architecture platform.Architecture
	Library      platform.Library
	Extension    platform.Extension
	Distribution platform.Distribution
}

func (p *Platform) Parse(name string) {
	p.OS.Parse(name)
	p.Architecture.Parse(name)
	p.Library.Parse(name)
}

func (p *Platform) Default() {
	p.OS = p.OS.Default()
	p.Architecture = p.Architecture.Default()
	p.Library = p.Library.Default(p.OS, p.Distribution)
	p.Extension = p.Extension.Default(p.OS)
	p.Distribution = p.Distribution.Default()
}

func (p *Platform) Merge(other Platform) {
	stringlike.SetIfEmpty(&p.OS, other.OS)
	stringlike.SetIfEmpty(&p.Architecture.Type, other.Architecture.Type)
	stringlike.SetIfEmpty(&p.Architecture.Version, other.Architecture.Version)
	stringlike.SetIfEmpty(&p.Library, other.Library)
	stringlike.SetIfEmpty(&p.Extension, other.Extension)
	stringlike.SetIfEmpty(&p.Distribution, other.Distribution)
}

func (p *Platform) CommonExtensions() []string {
	switch p.OS {
	case platform.Windows:
		return []string{
			".zip",
			".exe",
			".gz",
		}
	default:
		return []string{
			".gz",
			"",
		}
	}
}
