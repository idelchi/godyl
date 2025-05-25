package processor

import (
	"time"

	"github.com/idelchi/godyl/internal/cache"
	"github.com/idelchi/godyl/internal/tools/tool"
)

// UpdateCache attempts to update the cache with tool version information.
func (p *Processor) UpdateCache(tool *tool.Tool) {
	if tool.NoCache || tool.Version.Version == "" || p.cache == nil {
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

	if err := p.cache.Set(item.ID, item); err != nil {
		p.log.Errorf("  failed to save cache: %v", err)
	}
}
