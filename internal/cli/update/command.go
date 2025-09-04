// Package update contains the subcommand definition for `update`.
package update

import (
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/cli/core"
	"github.com/idelchi/godyl/internal/config/root"
	"github.com/idelchi/godyl/internal/config/update"
)

// Command returns the `update` command.
func Command(global *root.Config, local any, embedded *core.Embedded) *cobra.Command {
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
		RunE: func(cmd *cobra.Command, args []string) error {
			// Exit early if the command is run with `--show/-s` flag.
			if core.ExitOnShow(global.ShowFunc) {
				return nil
			}

			return run(core.Input{Global: global, Cmd: cmd, Args: args, Embedded: embedded})
		},
	}
	core.SetSubcommandDefaults(cmd, local, global.ShowFunc)

	update.Flags(cmd)

	return cmd
}
