package set

import (
	"errors"

	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/cli/common"
	"github.com/idelchi/godyl/internal/config"
	"github.com/idelchi/godyl/pkg/koanfx"
)

// Command returns the `cache path` command.
func Command(global *config.Config, local any) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set <key> <value>",
		Short: "Interact with the config",
		Long:  "Interact with the config.",
		Args:  cobra.ExactArgs(2), //nolint:mnd	// The command takes 2 arguments as documented.
		RunE: func(cmd *cobra.Command, args []string) error {
			// Exit early if the command is run with `--show/-s` flag.
			if common.ExitOnShow(global.ShowFunc) {
				return nil
			}

			koanf, ok := cmd.Root().Context().Value("config").(*koanfx.KoanfWithTracker)
			if !ok {
				return errors.New("failed to get config from context")
			}

			return run(global.ConfigFile.Absolute(), koanf, args[0], args[1])
		},
	}

	common.SetSubcommandDefaults(cmd, local, global.ShowFunc)

	return cmd
}
