package dump

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/config"
	"github.com/idelchi/godyl/internal/core/defaults"
	"github.com/idelchi/godyl/internal/tools"
	iutils "github.com/idelchi/godyl/internal/utils"
)

func NewDefaultsCommand(cfg *config.Config, files config.Embedded) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "defaults",
		Short: "Display default configuration settings",
		Args:  cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			return commonPreRunE(cmd, &cfg.Dump)
		},
		RunE: func(_ *cobra.Command, _ []string) error {
			c, err := getDefaults(cfg, files)
			if err != nil {
				return err
			}

			iutils.Print(cfg.Dump.Format, c)

			return nil
		},
	}

	return cmd
}

func getDefaults(cfg *config.Config, files config.Embedded) (*tools.Defaults, error) {
	tools := &tools.Defaults{}
	if err := defaults.LoadDefaults(tools, cfg.Root.Defaults.Name(), files.Defaults, *cfg); err != nil {
		return nil, fmt.Errorf("loading defaults: %w", err)
	}

	return tools, nil
}
