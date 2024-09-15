// Package platform implements the platform dump subcommand for godyl.
// It displays information about the detected system platform.
package platform

import (
	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/cli/common"
	"github.com/idelchi/godyl/internal/config"
)

func Command(global *config.Config, local any) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "platform",
		Short: "Display platform information",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			// Exit early if the command is run with `--show/-s` flag.
			if global.Root.ShowFunc() != nil {
				return nil
			}

			return run(global.Dump.Format)
		},
	}
	common.SetSubcommandDefaults(cmd, local, global.Root.ShowFunc)

	return cmd
}
