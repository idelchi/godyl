// Package remove contains the subcommand definition for `config remove`.
package remove

import (
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/cli/core"
	"github.com/idelchi/godyl/internal/config/root"
)

// Command returns the `config remove` command.
func Command(global *root.Config, local any) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove [key...]",
		Short: "Remove configuration keys.",
		Long:  "Remove all or specific keys.",
		Example: heredoc.Doc(`
			# Remove all configuration keys
			$ godyl config remove

			# Remove a specific keys from the configuration
			$ godyl config remove idelchi/envprof dump.tools.full dump.config.full
		`),

		Aliases: []string{"rm"},
		Args:    cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Exit early if the command is run with `--show/-s` flag.
			if global.ShowFunc() != nil {
				return nil
			}

			return run(core.Input{Global: global, Cmd: cmd, Args: args})
		},
	}

	core.SetSubcommandDefaults(cmd, local, global.ShowFunc)

	return cmd
}
