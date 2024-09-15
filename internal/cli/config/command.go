// Package config contains the subcommand definition for `config`.
package config

import (
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/cli/core"
	"github.com/idelchi/godyl/internal/config/root"
	"github.com/idelchi/godyl/pkg/cobraext"
)

// Command returns the `config` command.
func Command(global *root.Config, local any) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config [command]",
		Short: "Interact with the config file.",
		Long: heredoc.Doc(`
			Display the path, remove, or set entries in the configuration file.
		`),
		Example: heredoc.Doc(`
			# Display the configuration file path
			$ godyl config path

			# Remove all entries from the configuration file
			$ godyl config remove

			# Remove specific keys from the configuration file
			$ godyl config remove no-progress dump update.cleanup
		`),
		Aliases: []string{"cfg"},
		Args:    cobra.NoArgs,
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
