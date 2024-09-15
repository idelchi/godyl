// Package cache contains the subcommand definition for `cache`.
package config

import (
	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/cli/config/path"
	"github.com/idelchi/godyl/internal/cli/config/set"
	"github.com/idelchi/godyl/internal/config"
)

// subcommands for the `config` command.
func subcommands(cmd *cobra.Command, cfg *config.Config) {
	cmd.AddCommand(
		path.Command(cfg, nil),
		set.Command(cfg, nil),
	)
}
