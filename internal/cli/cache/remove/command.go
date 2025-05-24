package remove

import (
	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/cli/common"
	"github.com/idelchi/godyl/internal/config"
)

// Command returns the `cache remove` command.
func Command(global *config.Config, local any) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "remove",
		Short:   "Remove the cache",
		Long:    "Remove the cache.",
		Aliases: []string{"rm"},
		Args:    cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			// Exit early if the command is run with `--show/-s` flag.
			if global.Root.ShowFunc() != nil {
				return nil
			}

			return run(global.Root.Cache.Dir)
		},
	}

	common.SetSubcommandDefaults(cmd, local, global.Root.ShowFunc)

	return cmd
}
