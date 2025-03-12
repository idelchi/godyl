// Package configuration implements the config dump subcommand for godyl.
// It displays the application's current configuration settings.
package configuration

import (
	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/cli/flags"
	"github.com/idelchi/godyl/internal/config"
	iutils "github.com/idelchi/godyl/internal/utils"
)

// Command encapsulates the config dump command with its associated configuration.
type Command struct {
	// Command is the config cobra.Command instance
	Command *cobra.Command
	// Config contains application configuration
	Config *config.Config
	// Files contains the embedded configuration files and templates
	Files config.Embedded
}

// Flags adds config-specific flags to the command.
func (cmd *Command) Flags() {
	// No specific flags for this command
}

// NewConfigCommand creates a Command for displaying application configuration.
func NewConfigCommand(cfg *config.Config, files config.Embedded) *Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Display root configuration information",
		Args:  cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			return flags.Bind(cmd.Parent(), &cfg.Dump, cmd.Root().Name(), cmd.Parent().Name())
		},
		RunE: func(_ *cobra.Command, _ []string) error {
			c, err := getConfig(cfg)
			if err != nil {
				return err
			}

			iutils.Print(cfg.Dump.Format, c)

			return nil
		},
	}

	return &Command{
		Command: cmd,
		Config:  cfg,
		Files:   files,
	}
}

// NewCommand creates a cobra.Command instance for the config dump subcommand.
func NewCommand(cfg *config.Config, files config.Embedded) *cobra.Command {
	// Create the config command
	cmd := NewConfigCommand(cfg, files)

	// Add config-specific flags
	cmd.Flags()

	return cmd.Command
}

// getConfig returns the current application configuration.
func getConfig(cfg *config.Config) (*config.Config, error) {
	return cfg, nil
}
