// Package update implements the update command for godyl.
// It provides functionality to update the application itself to the latest version.
package update

import (
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/cli/common"
	"github.com/idelchi/godyl/internal/config"
	"github.com/idelchi/godyl/internal/config/update"
)

func Command(global *config.Config, local any, embedded *common.Embedded) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update",
		Aliases: []string{"upgrade", "u"},
		Short:   "Update the application",
		Example: heredoc.Doc(`
			Update to a specific version:
			$ godyl update --version=v0.0.10
			Update to the latest pre-release version:
			$ godyl update --pre
			Check the latest version available (including pre-releases):
			$ godyl update --pre --check
			Update to the latest version with cleanup (windows only):
			$ godyl update --cleanup
		`),
		RunE: func(cmd *cobra.Command, _ []string) error {
			// Exit early if the command is run with `--show/-s` flag.
			if common.ExitOnShow(global.Root.ShowFunc) {
				return nil
			}

			return run(*global, *embedded, cmd.Root().Version)
		},
	}
	common.SetSubcommandDefaults(cmd, local, global.Root.ShowFunc)

	update.Flags(cmd)

	return cmd
}
