// Package cache contains the subcommand definition for `dump cache`.
package cache

import (
	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/cli/core"
	"github.com/idelchi/godyl/internal/config/root"
)

// Command returns the `dump cache` command.
func Command(global *root.Config, local any) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cache [name...]",
		Short: "Display cache information",
		Args:  cobra.ArbitraryArgs,

		RunE: func(cmd *cobra.Command, args []string) error {
			// Exit early if the command is run with `--show/-s` flag.
			if core.ExitOnShow(global.ShowFunc) {
				return nil
			}

			return run(core.Input{Global: global, Cmd: cmd, Args: args})
		},
	}

	core.SetSubcommandDefaults(cmd, local, global.ShowFunc)

	return cmd
}
