// Package defaults implements the defaults dump subcommand for godyl.
// It displays the application's default configuration settings.
package defaults

import (
	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/cli/common"
	"github.com/idelchi/godyl/internal/config"
)

func Command(global *config.Config, local any, embedded *common.Embedded) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "defaults [name...]",
		Short: "Display default configuration settings",
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Exit early if the command is run with `--show/-s` flag.
			if common.ExitOnShow(global.ShowFunc) {
				return nil
			}

			return run(*global, *embedded, args...)
		},
	}

	common.SetSubcommandDefaults(cmd, local, global.ShowFunc)

	return cmd
}
