package clean

import (
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/cli/common"
	"github.com/idelchi/godyl/internal/config"
)

// Command returns the `cache clean` command.
func Command(global *config.Config, local any) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "clean",
		Short: "Clean the cache",
		Long: heredoc.Doc(`Clean can be run to clear the cache from removed tools,
		as well as updating the recorded versions in case of mismatches.
			`),
		Example: heredoc.Doc(`
			$ godyl cache clean
			`),
		Args: cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			// Exit early if the command is run with `--show/-s` flag.
			if common.ExitOnShow(global.Root.ShowFunc) {
				return nil
			}

			return run(*global)
		},
	}

	common.SetSubcommandDefaults(cmd, local, global.Root.ShowFunc)

	return cmd
}
