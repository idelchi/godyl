package tool

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/config"
	"github.com/idelchi/godyl/internal/core/updater"
)

func NewUpdateCommand(cfg *config.Config, files config.Embedded) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update",
		Aliases: []string{"upgrade", "up"},
		Short:   "Update the application",
		Long:    "Update the godyl application to the latest version",
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			return commonPreRunE(cmd, &cfg.Tool)
		},
		RunE: func(_ *cobra.Command, _ []string) error {
			appUpdater := updater.Updater{
				Strategy:    cfg.Tool.Strategy,
				NoVerifySSL: cfg.Tool.NoVerifySSL,
				Template:    files.Template,
			}

			if err := appUpdater.Update(""); err != nil {
				return fmt.Errorf("updating: %w", err)
			}

			return nil
		},
	}

	addUpdateFlags(cmd)

	return cmd
}
