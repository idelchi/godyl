// Package version contains the subcommand definition for `version`.
package version

import (
	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/cli/core"
	"github.com/idelchi/godyl/internal/config/root"
)

// Command returns the `version` command.
func Command(global *root.Config, local any) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "version",
		Short:   "Print the version number of godyl",
		Long:    `Print the version number of godyl along with available updates`,
		Aliases: []string{"v"},
		Args:    cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			// Exit early if the command is run with `--show/-s` flag.
			if core.ExitOnShow(global.ShowFunc) {
				return
			}

			run(core.Input{Global: global, Cmd: cmd, Args: args})
		},
	}

	core.SetSubcommandDefaults(cmd, local, global.ShowFunc)

	return cmd
}
