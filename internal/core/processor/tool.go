package processor

import (
	"fmt"

	"github.com/hashicorp/go-getter/v2"
	"github.com/idelchi/godyl/internal/tools"
	"github.com/idelchi/godyl/pkg/pretty"
)

// processOneTool processes an individual tool and sends the result to the result channel.
func (p *Processor) processOneTool(
	tool *tools.Tool,
	tags tools.IncludeTags,
	resultCh chan<- result,
	progressTracker getter.ProgressTracker,
	dry bool,
) {
	// Apply defaults and resolve tool configuration
	tool.ApplyDefaults(p.defaults, p.cache)

	if res := tool.Resolve(tags); res.Status != tools.Installed {
		resultCh <- result{Tool: tool, Details: details{err: res.Error()}}
		return
	}

	// Handle dry run
	if dry {
		resultCh <- result{Tool: tool, Details: details{err: fmt.Errorf("tool can be upgraded to %q", tool.Version.Version)}}
		return
	}

	// Apply SSL verification setting
	if p.config.Tool.NoVerifySSL {
		tool.NoVerifySSL = true
	}

	// Log tool configuration in debug mode
	p.logToolConfiguration(tool)

	// Download the tool
	msg, err := tool.Download(progressTracker)
	resultCh <- result{Tool: tool, Details: details{err: err, messages: []string{msg}}}
}

// logToolConfiguration logs the tool configuration at debug level
func (p *Processor) logToolConfiguration(tool *tools.Tool) {
	p.log.Debug("Tool configuration:")
	p.log.Debug("-------")
	p.log.Debug("%s", pretty.YAMLMasked(tool))
	p.log.Debug("-------")
}
