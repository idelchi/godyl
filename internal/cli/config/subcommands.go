// Package config contains the subcommand definition for `config`.
package config

import (
	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/cli/config/path"
	"github.com/idelchi/godyl/internal/cli/config/remove"
	"github.com/idelchi/godyl/internal/cli/config/set"
	"github.com/idelchi/godyl/internal/config/root"
)

// subcommands for the `config` command.
func subcommands(cmd *cobra.Command, cfg *root.Config) {
	cmd.AddCommand(
		path.Command(cfg, nil),
		set.Command(cfg, nil),
		remove.Command(cfg, nil),
	)
}
