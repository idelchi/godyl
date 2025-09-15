package install

import (
	"fmt"

	"github.com/idelchi/godyl/internal/cli/core"
	"github.com/idelchi/godyl/internal/iutils"
	"github.com/idelchi/godyl/internal/processor"
	"github.com/idelchi/godyl/internal/tools"
	"github.com/idelchi/godyl/pkg/unmarshal"
)

// run executes the `install` command.
func run(input core.Input) error {
	cfg, embedded, _, _, args := input.Unpack()

	if cfg.Install.Dry {
		cfg.Verbose = 1
	}

	// Load the tools from the source as []byte
	data, err := iutils.ReadPathsOrDefault(cfg.Tools, args...)
	if err != nil {
		return fmt.Errorf("reading tools file: %w", err)
	}

	// The tools can now be unmarshalled into a tools.Tools instance
	var tools tools.Tools
	if err := unmarshal.Strict(data, &tools); err != nil {
		return fmt.Errorf("unmarshalling tools: %w", err)
	}

	// Generate a common configuration for the command
	cfg.Common = cfg.Install.ToCommon()

	runner := core.NewHandler(*cfg, *embedded)
	if err := runner.SetupLogger(cfg.LogLevel); err != nil {
		return fmt.Errorf("setting up logger: %w", err)
	}

	if err := runner.Resolve(cfg.Defaults, &tools); err != nil {
		return err
	}

	// At this point, all tools have been resolved and can be processed by the processor
	proc := processor.New(tools, *cfg, runner.Logger())

	proc.NoDownload = cfg.Install.Dry

	if err := proc.Process(iutils.SplitTags(cfg.Install.Tags)); err != nil {
		return fmt.Errorf("processing tools: %w", err)
	}

	return nil
}
