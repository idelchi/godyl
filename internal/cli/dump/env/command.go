// Package env implements the env dump subcommand for godyl.
// It displays information about the environment variables.
package env

import (
	"errors"

	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/cli/common"
	"github.com/idelchi/godyl/internal/config"
	"github.com/idelchi/godyl/pkg/env"
)

func Command(global *config.Config, local any) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "env",
		Short: "Display environment information",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			// Exit early if the command is run with `--show/-s` flag.
			if common.ExitOnShow(global.Root.ShowFunc) {
				return nil
			}

			dotenv, ok := cmd.Root().Context().Value("dotenv").(env.Env)
			if !ok {
				return errors.New("failed to get dotenv from context")
			}

			return run(dotenv)
		},
	}

	common.SetSubcommandDefaults(cmd, local, global.Root.ShowFunc)

	return cmd
}
