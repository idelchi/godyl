package update

import (
	"fmt"

	"github.com/idelchi/godyl/internal/cli/common"
	"github.com/idelchi/godyl/internal/config"
	"github.com/idelchi/godyl/internal/tools"
	"github.com/idelchi/godyl/internal/updater"
)

// run executes the `cache dump` command.
func run(cfg config.Config, embedded common.Embedded, version string) error {
	// Generate a common configuration for the command
	cfg.Common = cfg.Update.ToCommon()

	godyl := updater.NewGodyl(version, &cfg)

	cfg.Inherit = "default"

	runner := common.NewHandler(cfg, embedded)
	if err := runner.SetupLogger(cfg.LogLevel); err != nil {
		return fmt.Errorf("setting up logger: %w", err)
	}

	if err := runner.Resolve("", &tools.Tools{godyl.Tool}); err != nil {
		return err
	}

	if !cfg.Update.Cleanup {
		embedded.Template = nil
	}

	updater := updater.New(&godyl, embedded.Template, runner.Logger())

	return updater.Update(cfg.Update.Check)
}
