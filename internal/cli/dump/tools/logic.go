package tools

import (
	"fmt"

	"github.com/idelchi/godyl/internal/cli/common"
	"github.com/idelchi/godyl/internal/config"
	"github.com/idelchi/godyl/internal/iutils"
	"github.com/idelchi/godyl/internal/tools"
	"github.com/idelchi/godyl/internal/tools/tags"
	"github.com/idelchi/godyl/pkg/unmarshal"
	"github.com/idelchi/godyl/pkg/utils"
)

func run(cfg config.Config, embedded common.Embedded, args ...string) (err error) {
	tags := iutils.SplitTags(cfg.Dump.Tools.Tags)

	data := embedded.Tools

	// Load the tools from the source as []byte
	if !cfg.Dump.Tools.Embedded {
		data, err = iutils.ReadPathsOrDefault(cfg.Tools, args...)
		if err != nil {
			return fmt.Errorf("reading tools file: %w", err)
		}
	}

	c, err := getTools(data, cfg.Dump.Tools.Full, tags)
	if err != nil {
		return err
	}

	iutils.Print(iutils.YAML, c)

	return nil
}

// getTools returns the tools configuration from embedded files.
func getTools(embeddedTools []byte, rendered bool, tags tags.IncludeTags) (any, error) {
	tools := tools.Tools{}

	err := unmarshal.Strict(embeddedTools, &tools)
	if err != nil {
		return nil, err
	}

	var included []int //nolint:prealloc  		// Size is unknown

	for i, tool := range tools {
		tool.Tags.Append(tool.Name)

		if !tool.Tags.Include(tags.Include) || !tool.Tags.Exclude(tags.Exclude) {
			continue
		}

		included = append(included, i)
	}

	if !rendered {
		var tools []any

		err := unmarshal.Strict(embeddedTools, &tools)
		if err != nil {
			return nil, err
		}

		return utils.PickByIndices(tools, included), nil
	}

	return utils.PickByIndices(tools, included), nil
}
