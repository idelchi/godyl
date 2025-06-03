package dump

import (
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/cli/common"
	"github.com/idelchi/godyl/internal/config"
	"github.com/idelchi/godyl/pkg/cobraext"
)

// Command creates a Command for displaying configuration information.
func Command(global *config.Config, local any, embedded *common.Embedded) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "dump",
		Short:   "Dump configuration information",
		Aliases: []string{"ls"},
		Example: heredoc.Doc(`
			$ godyl dump --format json defaults
			$ godyl dump env
			$ godyl dump platform
			$ godyl dump tools
		`),
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
