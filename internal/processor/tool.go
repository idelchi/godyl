package processor

import (
	"github.com/hashicorp/go-getter/v2"

	"github.com/idelchi/godyl/internal/tools/tags"
	"github.com/idelchi/godyl/internal/tools/tool"
	"github.com/idelchi/godyl/pkg/pretty"
)

// ProcessOneTool handles the processing of a single tool, including configuration,
// resolution, and downloading. It sends the processing result to the result channel.
func (p *Processor) processOneTool(
	tool *tool.Tool,
	tags tags.IncludeTags,
	resultCh chan<- *tool.Tool,
	progressTracker getter.ProgressTracker,
) {
	tool.EnableCache(p.cache)
	res := tool.Resolve(tags, p.Options...)

	p.logToolConfiguration(tool)
	tool.SetResult(res)

	if p.NoDownload || !res.IsOK() || p.Options != nil {
		tool.DisableCache()

		resultCh <- tool

		return
	}

	// Apply SSL verification setting
	if p.config.Root.NoVerifySSL {
		tool.NoVerifySSL = true
	}

	// Download the tool
	res = tool.Download(progressTracker)

	tool.SetResult(res)

	resultCh <- tool
}

// LogToolConfiguration logs the complete tool configuration at debug level.
func (p *Processor) logToolConfiguration(tool *tool.Tool) {
	p.log.Debug("Tool:")
	p.log.Debug("-------")
	p.log.Debugf("%s", pretty.YAML(tool))
	p.log.Debug("-------")
}
