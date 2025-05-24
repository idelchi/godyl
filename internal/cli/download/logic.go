package download

import (
	"fmt"

	"github.com/idelchi/godyl/internal/cli/common"
	"github.com/idelchi/godyl/internal/config"
	"github.com/idelchi/godyl/internal/processor"
	"github.com/idelchi/godyl/internal/tools"
	"github.com/idelchi/godyl/internal/tools/mode"
	"github.com/idelchi/godyl/internal/tools/sources"
	"github.com/idelchi/godyl/internal/tools/strategy"
	"github.com/idelchi/godyl/internal/tools/tags"
	"github.com/idelchi/godyl/internal/tools/tool"
	"github.com/idelchi/godyl/internal/tools/version"
	"github.com/idelchi/godyl/pkg/path/file"
	"github.com/idelchi/godyl/pkg/utils"
)

func run(global config.Config, embedded common.Embedded, args ...string) error {
	tools := tools.Tools{}

	for _, name := range args {
		tool := &tool.Tool{
			Mode:     mode.Extract,
			Strategy: strategy.Force,
			Version: version.Version{
				Version: global.Download.Version,
			},
		}

		if utils.IsURL(name) {
			tool.Name = file.New(name).Base()
			tool.URL = name
			tool.Source.Type = sources.URL
		} else {
			tool.Name = name
			tool.Source.Type = sources.GITHUB
		}

		tools.Append(tool)
	}

	// Generate a common configuration for the command
	global.Common = global.Download.ToCommon()
	global.Root.Cache.Disabled = true

	runner := common.NewHandler(global, embedded)
	if err := runner.SetupLogger(global.Root.LogLevel); err != nil {
		return fmt.Errorf("setting up logger: %w", err)
	}

	if err := runner.Resolve(global.Root.Defaults, &tools); err != nil {
		return err
	}

	// Process tools
	proc := processor.New(tools, global, runner.Logger())

	proc.NoDownload = global.Download.Dry
	if err := proc.Process(tags.IncludeTags{}); err != nil {
		return fmt.Errorf("processing tools: %w", err)
	}

	return nil
}
