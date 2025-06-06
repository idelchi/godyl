package config

import (
	"errors"

	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/cli/common"
	"github.com/idelchi/godyl/internal/config"
	cconfig "github.com/idelchi/godyl/internal/config/dump/config"
	"github.com/idelchi/godyl/pkg/koanfx"
)

// Command returns the `config dump` command.
func Command(global *config.Config, local any) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "config [items...]",
		Short:   "Dump cache information",
		Aliases: []string{"ls"},
		Args:    cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Exit early if the command is run with `--show/-s` flag.
			if common.ExitOnShow(global.ShowFunc) {
				return nil
			}

			koanf, ok := cmd.Root().Context().Value("config").(*koanfx.KoanfWithTracker)
			if !ok {
				return errors.New("failed to get config from context")
			}

			return run(*global, koanf, args...)
		},
	}

	common.SetSubcommandDefaults(cmd, local, global.ShowFunc)

	cconfig.Flags(cmd)

	return cmd
}
