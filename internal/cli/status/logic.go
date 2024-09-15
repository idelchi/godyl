package status

import (
	"fmt"

	"github.com/idelchi/godyl/internal/cli/common"
	"github.com/idelchi/godyl/internal/iutils"
	"github.com/idelchi/godyl/internal/processor"
	"github.com/idelchi/godyl/internal/tools"
	"github.com/idelchi/godyl/internal/tools/tool"
	"github.com/idelchi/godyl/pkg/unmarshal"
)

// TODO(Idelchi): Presentation must look different for status (green -> yellow, yellow -> green, red -> red).

// run executes the `status` command.
func run(input common.Input) error {
	cfg, embedded, _, _, args := input.Unpack()

	// Always set the verbose level to 1 for the status command
	cfg.Verbose = 1

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
	cfg.Common = cfg.Status.ToCommon()

	runner := common.NewHandler(*cfg, *embedded)
	if err := runner.SetupLogger(cfg.LogLevel); err != nil {
		return fmt.Errorf("setting up logger: %w", err)
	}

	if err := runner.Resolve(cfg.Defaults, &tools); err != nil {
		return err
	}
	// At this point, all tools have been resolved and can be processed by the processor
	proc := processor.New(tools, *cfg, runner.Logger())
	proc.NoDownload = true
	proc.Options = []tool.ResolveOption{tool.WithoutURL()}

	if err := proc.Process(iutils.SplitTags(cfg.Status.Tags)); err != nil {
		return fmt.Errorf("processing tools: %w", err)
	}

	return nil
}
