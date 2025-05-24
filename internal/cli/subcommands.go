package cli

import (
	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/cli/cache"
	"github.com/idelchi/godyl/internal/cli/common"
	cconfig "github.com/idelchi/godyl/internal/cli/config"
	"github.com/idelchi/godyl/internal/cli/download"
	"github.com/idelchi/godyl/internal/cli/dump"
	"github.com/idelchi/godyl/internal/cli/install"
	"github.com/idelchi/godyl/internal/cli/status"
	"github.com/idelchi/godyl/internal/cli/update"
	"github.com/idelchi/godyl/internal/cli/validate"
	"github.com/idelchi/godyl/internal/cli/version"
	"github.com/idelchi/godyl/internal/config"
)

// Subcommands adds all subcommands to the root command.
func subcommands(cmd *cobra.Command, global *config.Config, embedded *common.Embedded) {
	cmd.AddCommand(
		install.Command(global, &global.Install, embedded),
		download.Command(global, &global.Download, embedded),
		status.Command(global, &global.Status, embedded),
		dump.Command(global, &global.Dump, embedded),
		update.Command(global, &global.Update, embedded),
		cache.Command(global, &global.Cache),
		cconfig.Command(global, &global.Config),
		validate.Command(global, nil, embedded),

		version.Command(global, nil),
	)
}
