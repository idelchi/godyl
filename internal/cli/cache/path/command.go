// Package path contains the subcommand definition for `cache path`.
package path

import (
	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/cli/common"
	"github.com/idelchi/godyl/internal/config/root"
)

// Command returns the `cache path` command.
func Command(global *root.Config, local any) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "path",
		Short: "Display the cache path.",
		Long:  "Display the path where the cache is stored, regardless if it exists or not.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Exit early if the command is run with `--show/-s` flag.
			if common.ExitOnShow(global.ShowFunc) {
				return nil
			}

			return run(common.Input{Global: global, Embedded: nil, Cmd: cmd, Args: args})
		},
	}

	common.SetSubcommandDefaults(cmd, local, global.ShowFunc)

	return cmd
}
