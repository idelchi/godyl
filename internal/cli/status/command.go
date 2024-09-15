package status

import (
	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/cli/common"
	"github.com/idelchi/godyl/internal/config"
)

func Command(global *config.Config, local any, embedded *config.Embedded) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "status [tools.yml]...",
		Aliases: []string{"diff", "s"},
		Short:   "Status of installed tools as specified in the YAML file(s).",
		Long:    "Status of installed tools as specified in the YAML file(s).",
		Args:    cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Exit early if the command is run with `--show/-s` flag.
			if common.ExitOnShow(global.Root.ShowFunc) {
				return nil
			}

			return run(*global, *embedded, args...)
		},
	}

	common.SetSubcommandDefaults(cmd, local, global.Root.ShowFunc)

	flags(cmd)

	return cmd
}
