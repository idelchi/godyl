// Package status contains the subcommand definition for `auth status`.
package status

import (
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/cli/core"
	"github.com/idelchi/godyl/internal/config/root"
)

// Command returns the `auth status` command.
func Command(global *root.Config, local any) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status [token...]",
		Short: "Show the status of authentication tokens.",
		Long: heredoc.Docf(`
			Shows which authentication tokens have non empty values for the current configuration,
			after parsing the full configuration.
		`),
		Example: heredoc.Doc(`
			# View which authentication tokens are set
			$ godyl auth status
		`),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Exit early if the command is run with `--show/-s` flag.
			if core.ExitOnShow(global.ShowFunc) {
				return nil
			}

			return run(core.Input{Global: global, Embedded: nil, Cmd: cmd, Args: args})
		},
	}

	core.SetSubcommandDefaults(cmd, local, global.ShowFunc)

	return cmd
}
