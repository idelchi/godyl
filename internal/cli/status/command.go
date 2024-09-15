// Package status contains the subcommand definition for `status`.
package status

import (
	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/cli/core"
	"github.com/idelchi/godyl/internal/config/root"
	"github.com/idelchi/godyl/internal/config/status"
)

// Command returns the `status` command.
func Command(global *root.Config, local any, embedded *core.Embedded) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "status [tools.yml|-]...",
		Aliases: []string{"diff", "s"},
		Short:   "Status of installed tools as specified in the YAML file(s).",
		Long:    "Status of installed tools as specified in the YAML file(s).",
		Args:    cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Exit early if the command is run with `--show/-s` flag.
			if core.ExitOnShow(global.ShowFunc) {
				return nil
			}

			return run(core.Input{Global: global, Cmd: cmd, Args: args, Embedded: embedded})
		},
	}

	core.SetSubcommandDefaults(cmd, local, global.ShowFunc)

	status.Flags(cmd)

	return cmd
}
