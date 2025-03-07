package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/config"
	"github.com/idelchi/godyl/internal/core/updater"
)

// NewUpdateCommand creates the update command for updating the application.
func NewUpdateCommand(cfg *config.Config, emb EmbeddedFiles) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update",
		Aliases: []string{"upgrade", "up"},
		Short:   "Update the application",
		Long:    "Update the godyl application to the latest version",
		RunE: func(cmd *cobra.Command, args []string) error {
			return processUpdate(*cfg)
		},
	}

	return cmd
}

// processUpdate handles the update process based on the configuration.
func processUpdate(cfg config.Config) error {
	appUpdater := updater.Updater{
		Strategy:    cfg.Update.Strategy,
		NoVerifySSL: cfg.NoVerifySSL,
	}

	if err := appUpdater.Update(""); err != nil {
		return fmt.Errorf("updating: %w", err)
	}

	return nil
}
