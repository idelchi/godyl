// Package set contains the subcommand definition for `config set`.
package set

import (
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/cli/core"
	"github.com/idelchi/godyl/internal/config/root"
)

// Command returns the `config set` command.
func Command(global *root.Config, local any) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set <key> <value>",
		Short: "Set a key in the configuration",
		Example: heredoc.Doc(`
			$ godyl config set dump.tools.full true
		`),
		Args: cobra.ExactArgs(2), //nolint:mnd	// The command takes 2 arguments as documented.
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
