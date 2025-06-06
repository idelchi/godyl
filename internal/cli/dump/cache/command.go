package cache

import (
	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/cli/common"
	"github.com/idelchi/godyl/internal/config"
)

// Command returns the `dump cache` command.
func Command(global *config.Config, local any) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "cache [name]",
		Short:   "Dump cache information",
		Aliases: []string{"ls"},
		Args:    cobra.MaximumNArgs(1),

		RunE: func(_ *cobra.Command, args []string) error {
			// Exit early if the command is run with `--show/-s` flag.
			if common.ExitOnShow(global.ShowFunc) {
				return nil
			}

			return run(*global, args)
		},
	}

	common.SetSubcommandDefaults(cmd, local, global.ShowFunc)

	return cmd
}
