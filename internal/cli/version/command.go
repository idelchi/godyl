// Package version provides the subcommand for printing the tool version.
package version

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/cli/common"
	"github.com/idelchi/godyl/internal/config"
)

func Command(global *config.Config, local any) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "version",
		Short:   "Print the version number of godyl",
		Long:    `Print the version number of godyl along with available updates`,
		Aliases: []string{"v"},
		Args:    cobra.NoArgs,
		Run: func(cmd *cobra.Command, _ []string) {
			// Exit early if the command is run with `--show/-s` flag.
			if common.ExitOnShow(global.ShowFunc) {
				return
			}

			fmt.Println(cmd.Root().Version)

			return
		},
	}

	common.SetSubcommandDefaults(cmd, local, global.ShowFunc)

	return cmd
}
