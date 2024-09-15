package auth

import (
	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/cli/auth/remove"
	"github.com/idelchi/godyl/internal/cli/auth/store"
	"github.com/idelchi/godyl/internal/config/root"
)

// subcommands for the `auth` command.
func subcommands(cmd *cobra.Command, global *root.Config) {
	cmd.AddCommand(
		remove.Command(global, nil),
		store.Command(global, nil),
	)
}
