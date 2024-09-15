// Package update implements the update command for godyl.
// It provides functionality to update the application itself to the latest version.
package update

import (
	"fmt"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/cli/common"
	"github.com/idelchi/godyl/internal/config"
	"github.com/idelchi/godyl/internal/tools"
	"github.com/idelchi/godyl/internal/updater"
)

// Command encapsulates the update cobra command with its associated config and embedded files.
type Command struct {
	// Command is the update cobra.Command instance
	Command *cobra.Command
}

// NewUpdateCommand creates a Command for updating the application to the latest version.
func NewUpdateCommand(cfg *config.Config, embedded *config.Embedded) *Command {
	cmd := &cobra.Command{
		Use:     "update",
		Aliases: []string{"upgrade", "u"},
		Short:   "Update the application",
		Example: heredoc.Doc(`
			Update to a specific version:
			$ godyl update --version=v0.0.10
			Update to the latest pre-release version:
			$ godyl update --pre
			Check the latest version available (including pre-releases):
			$ godyl update --pre --check
			Update to the latest version with cleanup (windows only):
			$ godyl update --cleanup
		`),

		PersistentPreRunE: common.KCreateSubcommandPreRunE("update", &cfg.Update, &cfg.Root.Show),
		RunE: func(cmd *cobra.Command, _ []string) error {
			if cfg.Root.Show {
				return nil
			}
			// if cfg.Update.Check {
			// 	latest := updater.Latest{}
			// 	if err := latest.Get(cfg.Update.Pre); err != nil {
			// 		return fmt.Errorf("checking for updates: %w", err)
			// 	}

			// 	if !version.Compare(cmd.Root().Version, latest.Version) {
			// 		log.Info("")

			// 		log.Infof("A new version %q is available!", latest.Version)
			// 		log.Info(latest.Changelog)
			// 		log.Info("")
			// 		log.Info("You can update with:\n\n  godyl update [--pre]")
			// 	} else {
			// 		log.Infof("You are using the latest version %q", cmd.Root().Version)
			// 	}

			// 	return nil
			// }

			// Generate a common configuration for the command
			cfg.SetCommon(cfg.Update.ToCommon())

			godyl := updater.NewGodyl(cmd.Root().Version, cfg)

			cfg.Root.Inherit = "default"

			runner := common.NewHandler(cfg, embedded)
			if err := runner.SetupLogger(cfg.Root.LogLevel); err != nil {
				return fmt.Errorf("setting up logger: %w", err)
			}
			if err := runner.Resolve(nil, &tools.Tools{godyl.Tool}); err != nil {
				return err
			}

			if !cfg.Update.Cleanup {
				embedded.Template = nil
			}

			updater := updater.New(&godyl, embedded.Template, runner.Logger())

			return updater.Update(cfg.Update.Check)
		},
	}

	return &Command{
		Command: cmd,
	}
}

// NewCommand creates a cobra.Command instance containing the update command.
func NewCommand(cfg *config.Config, files *config.Embedded) *cobra.Command {
	// Create the update command
	cmd := NewUpdateCommand(cfg, files)

	// Add update-specific flags
	cmd.Flags()

	return cmd.Command
}
