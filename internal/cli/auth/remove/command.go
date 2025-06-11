// Package remove contains the subcommand definition for `auth remove`.
package remove

import (
	"strings"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/cli/common"
	"github.com/idelchi/godyl/internal/config/root"
	"github.com/idelchi/godyl/internal/iutils"
	"github.com/idelchi/godyl/pkg/pretty"
)

// Command returns the `auth remove` command.
func Command(global *root.Config, local any) *cobra.Command {
	tokens, _ := iutils.StructToKoanf(global.Tokens)

	cmd := &cobra.Command{
		Use:   "remove [token...]",
		Short: "Remove authentication tokens.",
		Long: heredoc.Docf(`
			Remove all or the specified tokens, either in the configuration file or in the keyring.

			Allowed values are:

			%v
		`, strings.TrimSpace(pretty.YAML(tokens.Keys()))),
		Example: heredoc.Doc(`
			# Remove all tokens
			$ godyl auth remove

			# Remove only the GitLab token
			$ godyl --keyring auth remove gitlab-token
		`),
		Aliases:   []string{"rm"},
		Args:      cobra.OnlyValidArgs,
		ValidArgs: tokens.Keys(),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Exit early if the command is run with `--show/-s` flag.
			if global.ShowFunc() != nil {
				return nil
			}

			return run(common.Input{Global: global, Embedded: nil, Cmd: cmd, Args: args})
		},
	}

	common.SetSubcommandDefaults(cmd, local, global.ShowFunc)

	return cmd
}
