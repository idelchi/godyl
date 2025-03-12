// Package update implements the update command for godyl.
// It provides functionality to update the application itself to the latest version.
package update

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/cli/flags"
	"github.com/idelchi/godyl/internal/config"
	"github.com/idelchi/godyl/internal/core/defaults"
	"github.com/idelchi/godyl/internal/core/updater"
)

// Command encapsulates the update cobra command with its associated config and embedded files.
type Command struct {
	// Command is the update cobra.Command instance
	Command *cobra.Command
}

// Flags adds update-specific flags to the command.
func (cmd *Command) Flags() {
	flags.Update(cmd.Command)
}

// NewUpdateCommand creates a Command for updating the application to the latest version.
func NewUpdateCommand(cfg *config.Config, files config.Embedded) *Command {
	cmd := &cobra.Command{
		Use:     "update",
		Aliases: []string{"upgrade", "up"},
		Short:   "Update the application",
		Long:    "Update the godyl application to the latest version",
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			return flags.Bind(cmd, &cfg.Tool, cmd.Root().Name(), cmd.Name())
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			defaults, err := defaults.Load(cfg.Root.Defaults.Name(), files, *cfg)
			if err != nil {
				return fmt.Errorf("loading defaults: %w", err)
			}

			appUpdater := updater.New(defaults, cfg.Tool.NoVerifySSL, files.Template)

			versions := updater.Versions{
				Current:   cmd.Root().Version,
				Requested: cfg.Tool.Version,
			}

			if err := appUpdater.Update(versions); err != nil {
				return fmt.Errorf("updating: %w", err)
			}

			return nil
		},
	}

	return &Command{
		Command: cmd,
	}
}

// NewCommand creates a cobra.Command instance containing the update command.
func NewCommand(cfg *config.Config, files config.Embedded) *cobra.Command {
	// Create the update command
	cmd := NewUpdateCommand(cfg, files)

	// Add update-specific flags
	cmd.Flags()

	return cmd.Command
}
