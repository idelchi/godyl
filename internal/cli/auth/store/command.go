// Package store contains the subcommand definition for `auth store`.
package store

import (
	"strings"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/cli/common"
	"github.com/idelchi/godyl/internal/config/root"
	"github.com/idelchi/godyl/internal/iutils"
	"github.com/idelchi/godyl/pkg/pretty"
)

// Command returns the `auth store` command.
func Command(global *root.Config, local any) *cobra.Command {
	tokens, _ := iutils.StructToKoanf(global.Tokens)

	cmd := &cobra.Command{
		Use:   "store [token...]",
		Short: "Store authentication tokens.",
		Long: heredoc.Docf(`
			Store all or the specified tokens, either in the configuration file or in the keyring.

			Allowed values are:

			%v
		`, strings.TrimSpace(pretty.YAML(tokens.Keys()))),
		Example: heredoc.Doc(`
			# Store all tokens in the default storage
			$ godyl auth store

			# Explicitly store all tokens in the keyring
			$ godyl --keyring auth store

			# Store only the GitHub and GitLab tokens, loading a custom environment file
			$ godyl --env-file=~/tokens.env auth store github-token gitlab-token

			# Store the GitHub token in the keyring, using an environment variable
			$ GODYL_GITHUB_TOKEN=token godyl auth store github-token
		`),
		Args:      cobra.OnlyValidArgs,
		ValidArgs: tokens.Keys(),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Exit early if the command is run with `--show/-s` flag.
			if common.ExitOnShow(global.ShowFunc) {
				return nil
			}

			return run(common.Input{Global: global, Embedded: nil, Cmd: cmd, Args: args})
		},
	}

	common.SetSubcommandDefaults(cmd, local, global.ShowFunc)

	return cmd
}
