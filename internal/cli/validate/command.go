// Package validate contains the subcommand definition for `validate`.
package validate

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/cli/common"
	"github.com/idelchi/godyl/internal/config/root"
)

// Command returns the `validate` command.
func Command(global *root.Config, local any, embedded *common.Embedded) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate the configuration for all subcommands",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			if err := run(common.Input{Global: global, Cmd: cmd, Args: args, Embedded: embedded}); err != nil {
				color.Red(err.Error())

				return
			}

			color.Green("Validation passed!")
		},
	}

	common.SetSubcommandDefaults(cmd, local, global.ShowFunc)

	return cmd
}
