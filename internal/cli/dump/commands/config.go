package dump

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/config"
	"github.com/idelchi/godyl/internal/core/defaults"
	iutils "github.com/idelchi/godyl/internal/utils"
)

func NewConfigCommand(cfg *config.Config, files config.Embedded) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Display root configuration information",
		Args:  cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			return commonPreRunE(cmd, &cfg.Dump)
		},
		RunE: func(_ *cobra.Command, _ []string) error {
			c, err := getConfig(cfg, files)
			if err != nil {
				return err
			}

			iutils.Print(cfg.Dump.Format, c)

			return nil
		},
	}

	return cmd
}

func getConfig(cfg *config.Config, files config.Embedded) (*config.Config, error) {
	defs := &defaults.Defaults{}
	if err := defs.Load(cfg.Root.Defaults.Name(), files.Defaults); err != nil {
		return nil, fmt.Errorf("error loading defaults: %w", err)
	}

	if err := defs.Merge(*cfg); err != nil {
		return nil, fmt.Errorf("error merging defaults: %w", err)
	}

	return cfg, nil
}
