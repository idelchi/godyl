package dump

import (
	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/cli/core"
	"github.com/idelchi/godyl/internal/cli/dump/auth"
	"github.com/idelchi/godyl/internal/cli/dump/cache"
	cconfig "github.com/idelchi/godyl/internal/cli/dump/config"
	"github.com/idelchi/godyl/internal/cli/dump/defaults"
	"github.com/idelchi/godyl/internal/cli/dump/env"
	"github.com/idelchi/godyl/internal/cli/dump/platform"
	"github.com/idelchi/godyl/internal/cli/dump/tools"
	"github.com/idelchi/godyl/internal/config/root"
)

// subcommands for the `dump` command.
func subcommands(cmd *cobra.Command, global *root.Config, embedded *core.Embedded) {
	cmd.AddCommand(
		defaults.Command(global, nil, embedded),
		env.Command(global, nil),
		platform.Command(global, nil),
		tools.Command(global, &global.Dump.Tools, embedded),
		cache.Command(global, nil),
		cconfig.Command(global, nil),
		auth.Command(global, nil),
	)
}
