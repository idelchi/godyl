// Package dump contains the subcommand definition for `dump`.
package dump

import (
	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/cli/common"
	"github.com/idelchi/godyl/internal/config/root"
	"github.com/idelchi/godyl/pkg/cobraext"
)

// Command returns the `dump` command.
func Command(global *root.Config, local any, embedded *common.Embedded) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "dump",
		Short:   "Dump configuration information",
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			// Since the command is allowed to run with `--show/-s` flag,
			// we should suppress the default error message for unknown subcommands.
			if common.ExitOnShow(global.ShowFunc, args...) {
				return nil
			}

			return cobraext.UnknownSubcommandAction(cmd, args)
		},
	}

	common.SetSubcommandDefaults(cmd, local, global.ShowFunc)

	subcommands(cmd, global, embedded)

	return cmd
}
