// Package defaults contains the subcommand definition for `dump defaults`.
package defaults

import (
	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/cli/common"
	"github.com/idelchi/godyl/internal/config/root"
)

// Command returns the `dump defaults` command.
func Command(global *root.Config, local any, embedded *common.Embedded) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "defaults [default...]",
		Short: "Display default configuration settings",
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Exit early if the command is run with `--show/-s` flag.
			if common.ExitOnShow(global.ShowFunc) {
				return nil
			}

			return run(common.Input{Global: global, Cmd: cmd, Args: args, Embedded: embedded})
		},
	}

	common.SetSubcommandDefaults(cmd, local, global.ShowFunc)

	return cmd
}
