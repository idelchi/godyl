// Package auth contains the subcommand definition for `auth`.
package auth

import (
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/cli/core"
	"github.com/idelchi/godyl/internal/config/root"
	"github.com/idelchi/godyl/pkg/cobraext"
)

// Command returns the `auth` command.
func Command(global *root.Config, local any) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth",
		Short: "Interact with authentication tokens.",
		Long: heredoc.Doc(`
			Store or remove authentication tokens, either in the configuration file or in the keyring.
			Check status of authentication tokens.
		`),
		Example: heredoc.Doc(`
			$ godyl auth store
			$ godyl --keyring auth store
			$ godyl auth remove
			$ godyl auth status
		`),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Since the command is allowed to run with `--show/-s` flag,
			// we should suppress the default error message for unknown subcommands.
			if core.ExitOnShow(global.ShowFunc, args...) {
				return nil
			}

			return cobraext.UnknownSubcommandAction(cmd, args)
		},
	}

	core.SetSubcommandDefaults(cmd, local, global.ShowFunc)

	subcommands(cmd, global)

	return cmd
}
