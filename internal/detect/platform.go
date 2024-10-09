package detect

import (
	"github.com/idelchi/godyl/internal/detect/platform"
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
	if p.OS == "" {
		p.OS = other.OS
	}
	if p.Architecture.Type == "" {
		p.Architecture.Type = other.Architecture.Type
	}
	if p.Architecture.Version == "" {
		p.Architecture.Version = other.Architecture.Version
	}
	if p.Library == "" {
		p.Library = other.Library
	}
	if p.Extension == "" {
		p.Extension = other.Extension
	}
	if p.Distribution == "" {
		p.Distribution = other.Distribution
	}
}

func (p Platform) ToInfo() Info {
	return Info{
		"os":           p.OS.String(),
		"architecture": p.Architecture.String(),
		"library":      p.Library.String(),
		"extension":    p.Extension.String(),
		"distribution": p.Distribution.String(),
	}
}

func (p *Platform) ToMap() map[string]string {
	return map[string]string{
		"os":           p.OS.String(),
		"architecture": p.Architecture.String(),
		"library":      p.Library.String(),
		"extension":    p.Extension.String(),
		"distribution": p.Distribution.String(),
	}
}
