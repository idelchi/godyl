package config

import (
	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/cli/common"
	"github.com/idelchi/godyl/internal/config"
	"github.com/idelchi/godyl/pkg/cobraext"
)

// Command returns the `config` command.
func Command(global *config.Config, local any) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "config [command]",
		Short:   "Interact with the config",
		Long:    "Interact with the config",
		Aliases: []string{"cfg"},
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Since the command is allowed to run with `--show/-s` flag,
			// we should suppress the default error message for unknown subcommands.
			if common.ExitOnShow(global.Root.ShowFunc, args...) {
				return nil
			}

			return cobraext.UnknownSubcommandAction(cmd, args)
		},
	}

	common.SetSubcommandDefaults(cmd, nil, global.Root.ShowFunc)

	subcommands(cmd, global)

	return cmd
}
