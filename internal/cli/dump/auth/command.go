// Package auth contains the subcommand definition for `dump auth`.
package auth

import (
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/cli/core"
	"github.com/idelchi/godyl/internal/config/root"
)

// Command returns the `dump auth` command.
func Command(global *root.Config, local any) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth",
		Short: "Display authentication tokens.",
		Example: heredoc.Doc(`
			$ godyl dump auth
			$ godyl --keyring dump auth
		`),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Exit early if the command is run with `--show/-s` flag.
			if global.ShowFunc() != nil {
				return nil
			}

			return run(core.Input{Global: global, Embedded: nil, Cmd: cmd, Args: args})
		},
	}
	core.SetSubcommandDefaults(cmd, local, global.ShowFunc)

	return cmd
}
