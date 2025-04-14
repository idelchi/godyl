package processor

import (
	"github.com/hashicorp/go-getter/v2"
	"github.com/idelchi/godyl/internal/tools"
	"github.com/idelchi/godyl/pkg/pretty"
)

// ProcessOneTool handles the processing of a single tool, including configuration,
// resolution, and downloading. It sends the processing result to the result channel.
func (p *Processor) processOneTool(
	tool *tools.Tool,
	tags tools.IncludeTags,
	resultCh chan<- result,
	progressTracker getter.ProgressTracker,
	dry bool,
) {
	// Apply defaults and resolve tool configuration
	tool.ApplyDefaults(p.defaults)
	tool.Cache(p.cache)

	if res := tool.Resolve(tags); !res.Successful() || dry {
		resultCh <- result{Tool: tool, Result: res}

		return
	}

	// Apply SSL verification setting
	if p.config.Tool.NoVerifySSL {
		tool.NoVerifySSL = true
	}

	// Log tool configuration in debug mode
	p.logToolConfiguration(tool)

	// Download the tool
	resultCh <- result{Tool: tool, Result: tool.Download(progressTracker)}
}

// LogToolConfiguration logs the complete tool configuration at debug level.
func (p *Processor) logToolConfiguration(tool *tools.Tool) {
	p.log.Debug("Tool configuration:")
	p.log.Debug("-------")
	p.log.Debug("%s", pretty.YAMLMasked(tool))
	p.log.Debug("-------")
}
