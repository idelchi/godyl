// Package env contains the subcommand definition for `dump env`.
package env

import (
	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/cli/core"
	"github.com/idelchi/godyl/internal/config/root"
)

// Command returns the `dump env` command.
func Command(global *root.Config, local any) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "env",
		Short: "Display environment information",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Exit early if the command is run with `--show/-s` flag.
			if core.ExitOnShow(global.ShowFunc) {
				return nil
			}

			return run(core.Input{Global: global, Cmd: cmd, Args: args, Embedded: nil})
		},
	}

	core.SetSubcommandDefaults(cmd, local, global.ShowFunc)

	return cmd
}
