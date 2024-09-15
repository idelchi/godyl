// Package cache contains the subcommand definition for `cache`.
package cache

import (
	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/cli/cache/clean"
	"github.com/idelchi/godyl/internal/cli/cache/path"
	"github.com/idelchi/godyl/internal/cli/cache/remove"
	"github.com/idelchi/godyl/internal/config"
)

// subcommands for the `cache` command.
func subcommands(cmd *cobra.Command, global *config.Config) {
	cmd.AddCommand(
		path.Command(global, nil),
		remove.Command(global, nil),
		clean.Command(global, nil),
	)
}
