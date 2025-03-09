package dump

import (
	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/config"
	iutils "github.com/idelchi/godyl/internal/utils"
	"github.com/idelchi/godyl/pkg/env"
)

func NewEnvCommand(cfg *config.Config, files config.Embedded) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "env",
		Short: "Display environment information",
		Args:  cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			return commonPreRunE(cmd, &cfg.Dump)
		},
		RunE: func(_ *cobra.Command, _ []string) error {
			c, err := getEnv()
			if err != nil {
				return err
			}

			iutils.Print(cfg.Dump.Format, c)

			return nil
		},
	}

	return cmd
}

func getEnv() (env.Env, error) {
	return env.FromEnv(), nil
}
