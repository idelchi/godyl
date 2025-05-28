package path

import (
	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/cli/common"
	"github.com/idelchi/godyl/internal/config"
)

func Command(global *config.Config, local any) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "path",
		Short: "Display the config path",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Exit early if the command is run with `--show/-s` flag.
			if common.ExitOnShow(global.ShowFunc) {
				return nil
			}

			return run(global.ConfigFile.Absolute())
		},
	}

	common.SetSubcommandDefaults(cmd, local, global.ShowFunc)

	return cmd
}
