package update

import (
	"fmt"

	"github.com/idelchi/godyl/internal/cli/common"
	"github.com/idelchi/godyl/internal/tools"
	"github.com/idelchi/godyl/internal/updater"
)

// run executes the `update` command.
func run(input common.Input) error {
	cfg, embedded, _, cmd, _ := input.Unpack()

	// Generate a common configuration for the command
	cfg.Common = cfg.Update.ToCommon()

	version := cmd.Root().Version
	godyl := updater.NewGodyl(version, cfg)

	cfg.Inherit = "default"

	handler := common.NewHandler(*cfg, *embedded)
	if err := handler.SetupLogger(cfg.LogLevel); err != nil {
		return fmt.Errorf("setting up logger: %w", err)
	}

	if err := handler.Resolve("", &tools.Tools{godyl.Tool}); err != nil {
		return err
	}

	if !cfg.Update.Cleanup {
		embedded.Template = nil
	}

	updater := updater.New(&godyl, embedded.Template, handler.Logger())

	return updater.Update(cfg.Update.Check)
}
