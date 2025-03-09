package update

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/config"
	"github.com/idelchi/godyl/internal/core/updater"
)

func NewCommand(cfg *config.Config, files config.Embedded) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update",
		Aliases: []string{"upgrade", "up"},
		Short:   "Update the application",
		Long:    "Update the godyl application to the latest version",
		RunE: func(_ *cobra.Command, _ []string) error {
			appUpdater := updater.Updater{
				Strategy:    cfg.Update.Strategy,
				NoVerifySSL: cfg.Update.NoVerifySSL,
				Template:    files.Template,
			}

			if err := appUpdater.Update(""); err != nil {
				return fmt.Errorf("updating: %w", err)
			}

			return nil
		},
	}

	cmd.Flags().String("strategy", "none", "Strategy to use for updating tools (none upgrade force)")
	cmd.Flags().String("github-token", os.Getenv("GODYL_GITHUB_TOKEN"), "GitHub token for authentication")
	cmd.Flags().BoolP("no-verify-ssl", "k", false, "Skip SSL verification")

	return cmd
}
