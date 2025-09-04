// Package remove contains the subcommand definition for `cache remove`.
package remove

import (
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/cli/core"
	"github.com/idelchi/godyl/internal/config/root"
)

// Command returns the `cache remove` command.
func Command(global *root.Config, local any) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove [name...]",
		Short: "Remove cache entries.",
		Long:  "Remove all or specific entries.",
		Example: heredoc.Doc(`
			# Remove all entries from the cache
			$ godyl cache remove

			# Remove a specific entry from the cache
			$ godyl cache remove idelchi/envprof
		`),
		Aliases: []string{"rm"},
		Args:    cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Exit early if the command is run with `--show/-s` flag.
			if global.ShowFunc() != nil {
				return nil
			}

			return run(core.Input{Global: global, Embedded: nil, Cmd: cmd, Args: args})
		},
	}

	core.SetSubcommandDefaults(cmd, local, global.ShowFunc)

	return cmd
}
