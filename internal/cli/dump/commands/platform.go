package dump

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/config"
	"github.com/idelchi/godyl/internal/detect"
	iutils "github.com/idelchi/godyl/internal/utils"
)

func NewPlatformCommand(cfg *config.Config, files config.Embedded) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "platform",
		Short: "Display platform information",
		Args:  cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			return commonPreRunE(cmd, &cfg.Dump)
		},
		RunE: func(_ *cobra.Command, _ []string) error {
			c, err := getPlatform()
			if err != nil {
				return err
			}

			iutils.Print(cfg.Dump.Format, c)

			return nil
		},
	}

	return cmd
}

func getPlatform() (*detect.Platform, error) {
	platform := &detect.Platform{}
	if err := platform.Detect(); err != nil {
		return nil, fmt.Errorf("detecting platform: %w", err)
	}

	return platform, nil
}
