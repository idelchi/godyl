package path

import (
	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/cli/common"
	"github.com/idelchi/godyl/internal/config"
)

// Command returns the `cache path` command.
func Command(global *config.Config, local any) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "path",
		Short: "Display the cache path",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			// Exit early if the command is run with `--show/-s` flag.
			if common.ExitOnShow(global.Root.ShowFunc) {
				return nil
			}

			return run(global.Root.Cache.Dir)
		},
	}

	common.SetSubcommandDefaults(cmd, local, global.Root.ShowFunc)

	return cmd
}
