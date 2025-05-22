package validate

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/cli/common"
	"github.com/idelchi/godyl/internal/config"
)

func Command(global *config.Config, local any, embedded *common.Embedded) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate the configuration for all subcommands",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			if err := run(*global, cmd, args); err != nil {
				color.Red(err.Error())

				return
			}

			color.Green("Validation passed!")
		},
	}

	common.SetSubcommandDefaults(cmd, local, global.Root.ShowFunc)

	return cmd
}
