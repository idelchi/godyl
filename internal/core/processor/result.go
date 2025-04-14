package processor

import (
	"github.com/idelchi/godyl/internal/tools"
	"github.com/idelchi/godyl/pkg/path/file"
)

// Result holds the outcome of processing a tool.
type result struct {
	// Tool is the tool that was processed.
	Tool *tools.Tool

	// Result contains the processing outcome and any error information.
	Result tools.Result
}

// UpdateCache attempts to update the cache with tool version information.
func (p *Processor) UpdateCache(tool *tools.Tool) {
	if tool.NoCache {
		return
	}

	if tool.Version.Version == "" {
		return
	}

	toolPath := file.New(tool.Output, tool.Exe.Name).Path()
	if err := p.cache.Save(toolPath, tool.Version.Version); err != nil {
		p.log.Error("  failed to save cache: %v", err)
	}
}
