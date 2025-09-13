package cli

import (
	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/cli/auth"
	"github.com/idelchi/godyl/internal/cli/cache"
	cconfig "github.com/idelchi/godyl/internal/cli/config"
	"github.com/idelchi/godyl/internal/cli/core"
	"github.com/idelchi/godyl/internal/cli/download"
	"github.com/idelchi/godyl/internal/cli/dump"
	"github.com/idelchi/godyl/internal/cli/install"
	"github.com/idelchi/godyl/internal/cli/paths"
	"github.com/idelchi/godyl/internal/cli/status"
	"github.com/idelchi/godyl/internal/cli/update"
	"github.com/idelchi/godyl/internal/cli/validate"
	"github.com/idelchi/godyl/internal/cli/version"
	"github.com/idelchi/godyl/internal/config/root"
)

// subcommands adds all subcommands to the root command.
func subcommands(cmd *cobra.Command, global *root.Config, embedded *core.Embedded) {
	cmd.AddCommand(
		install.Command(global, &global.Install, embedded),
		download.Command(global, &global.Download, embedded),
		status.Command(global, &global.Status, embedded),
		dump.Command(global, nil, embedded),
		update.Command(global, &global.Update, embedded),
		cache.Command(global, nil),
		cconfig.Command(global, nil),
		validate.Command(global, nil, embedded),
		auth.Command(global, nil),
		paths.Command(global, nil),

		version.Command(global, nil),
	)
}
