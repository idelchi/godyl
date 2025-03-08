package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/config"
	"github.com/idelchi/godyl/internal/core/updater"
)

// NewUpdateCommand creates the update command for updating the application.
func NewUpdateCommand(cfg *config.Config, files Embedded) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update",
		Aliases: []string{"upgrade", "up"},
		Short:   "Update the application",
		Long:    "Update the godyl application to the latest version",
		RunE: func(_ *cobra.Command, _ []string) error {
			appUpdater := updater.Updater{
				Strategy:    cfg.Update.Strategy,
				NoVerifySSL: cfg.NoVerifySSL,
				Template:    files.Template,
			}

			if err := appUpdater.Update(""); err != nil {
				return fmt.Errorf("updating: %w", err)
			}

			return nil
		},
	}

	return cmd
}
