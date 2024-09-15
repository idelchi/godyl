// Package path contains the subcommand definition for `config path`.
package path

import (
	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/cli/core"
	"github.com/idelchi/godyl/internal/config/root"
)

// Command returns the `config path` command.
func Command(global *root.Config, local any) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "path",
		Short: "Display the config path",
		Long:  "Display the path where the configuration is stored, regardless if it exists or not.",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			// Exit early if the command is run with `--show/-s` flag.
			if core.ExitOnShow(global.ShowFunc) {
				return nil
			}

			return run(core.Input{Global: global})
		},
	}

	core.SetSubcommandDefaults(cmd, local, global.ShowFunc)

	return cmd
}
