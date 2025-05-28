package status

import (
	"fmt"

	"github.com/idelchi/godyl/internal/cli/common"
	"github.com/idelchi/godyl/internal/config"
	"github.com/idelchi/godyl/internal/iutils"
	"github.com/idelchi/godyl/internal/processor"
	"github.com/idelchi/godyl/internal/tools"
	"github.com/idelchi/godyl/internal/tools/tool"
	"github.com/idelchi/godyl/pkg/unmarshal"
)

func run(global config.Config, embedded common.Embedded, args ...string) error {
	// Load the tools from the source as []byte
	data, err := iutils.ReadPathsOrDefault("tools.yml", args...)
	if err != nil {
		return fmt.Errorf("reading arguments %v: %w", args, err)
	}

	// The tools can now be unmarshalled into a tools.Tools instance
	var tools tools.Tools
	if err := unmarshal.Strict(data, &tools); err != nil {
		return fmt.Errorf("unmarshalling tools: %w", err)
	}

	// Generate a common configuration for the command
	global.Common = global.Status.ToCommon()

	runner := common.NewHandler(global, embedded)
	if err := runner.SetupLogger(global.LogLevel); err != nil {
		return fmt.Errorf("setting up logger: %w", err)
	}

	if err := runner.Resolve(global.Defaults, &tools); err != nil {
		return err
	}
	// At this point, all tools have been resolved and can be processed by the processor
	proc := processor.New(tools, global, runner.Logger())
	proc.NoDownload = true
	proc.Options = []tool.ResolveOption{tool.WithoutURL()}

	if err := proc.Process(iutils.SplitTags(global.Status.Tags)); err != nil {
		return fmt.Errorf("processing tools: %w", err)
	}

	return nil
}
