// Package cache contains the subcommand definition for `cache`.
package cache

import (
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/cli/core"
	"github.com/idelchi/godyl/internal/config/root"
	"github.com/idelchi/godyl/pkg/cobraext"
)

// Command returns the `cache` command.
func Command(global *root.Config, local any) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cache",
		Short: "Interact with the cache.",
		Long: heredoc.Doc(`
			Display the path, clean, or remove entries from the cache.
		`),
		Example: heredoc.Doc(`
			# Display the cache path
			$ godyl cache path

			# Clean the cache
			$ godyl cache clean

			# Remove all entries from the cache
			$ godyl cache remove

			# Remove a specific entry from the cache
			$ godyl cache remove idelchi/envprof
		`),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Since the command is allowed to run with `--show/-s` flag,
			// we should suppress the default error message for unknown subcommands.
			if core.ExitOnShow(global.ShowFunc, args...) {
				return nil
			}

			return cobraext.UnknownSubcommandAction(cmd, args)
		},
	}

	core.SetSubcommandDefaults(cmd, local, global.ShowFunc)

	subcommands(cmd, global)

	return cmd
}
