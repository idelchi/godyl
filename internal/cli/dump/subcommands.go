package dump

import (
	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/cli/common"
	"github.com/idelchi/godyl/internal/cli/dump/cache"
	cconfig "github.com/idelchi/godyl/internal/cli/dump/config"
	"github.com/idelchi/godyl/internal/cli/dump/defaults"
	"github.com/idelchi/godyl/internal/cli/dump/env"
	"github.com/idelchi/godyl/internal/cli/dump/platform"
	"github.com/idelchi/godyl/internal/cli/dump/tools"
	"github.com/idelchi/godyl/internal/config"
)

func subcommands(cmd *cobra.Command, global *config.Config, embedded *common.Embedded) {
	cmd.AddCommand(
		defaults.Command(global, nil, embedded),
		env.Command(global, nil),
		platform.Command(global, nil),
		tools.Command(global, &global.Dump.Tools, embedded),
		cache.Command(global, nil),
		cconfig.Command(global, &global.Dump.Config),
	)
}
