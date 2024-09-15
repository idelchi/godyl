package processor

import (
	"time"

	"github.com/idelchi/godyl/internal/cache"
	"github.com/idelchi/godyl/internal/tools/tool"
)

// UpdateCache attempts to update the cache with tool version information.
func (p *Processor) UpdateCache(tool *tool.Tool) {
	if p.config.Root.Cache.Disabled {
		return
	}

	if tool.Version.Version == "" {
		return
	}

	now := time.Now()

	item := &cache.Item{
		ID:         tool.ID(),
		Name:       tool.Name,
		Path:       tool.AbsPath(),
		Version:    tool.Version,
		Downloaded: now,
		Updated:    now,
		Type:       tool.Source.Type.String(),
	}

	if err := p.cache.Save(item); err != nil {
		p.log.Errorf("  failed to save cache: %v", err)
	}
}
