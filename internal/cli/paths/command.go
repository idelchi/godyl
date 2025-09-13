// Package paths contains the subcommand definition for `paths`.
package paths

import (
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/cli/core"
	"github.com/idelchi/godyl/internal/config/root"
)

// Command returns the `paths` command.
func Command(global *root.Config, local any) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "paths",
		Short: "Display all paths used by godyl",
		Example: heredoc.Doc(`
			$ godyl paths
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Exit early if the command is run with `--show/-s` flag.
			if core.ExitOnShow(global.ShowFunc) {
				return nil
			}

			return run(core.Input{Global: global, Cmd: cmd, Args: args})
		},
	}
	core.SetSubcommandDefaults(cmd, local, global.ShowFunc)

	return cmd
}
