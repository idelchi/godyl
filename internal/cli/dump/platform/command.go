// Package platform contains the subcommand definition for `dump platform`.
package platform

import (
	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/cli/common"
	"github.com/idelchi/godyl/internal/config/root"
)

// Command returns the `dump platform` command.
func Command(global *root.Config, local any) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "platform",
		Short: "Display platform information",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Exit early if the command is run with `--show/-s` flag.
			if global.ShowFunc() != nil {
				return nil
			}

			return run(common.Input{Global: global, Cmd: cmd, Args: args, Embedded: nil})
		},
	}
	common.SetSubcommandDefaults(cmd, local, global.ShowFunc)

	return cmd
}
