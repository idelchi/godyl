// Package defaults implements the defaults dump subcommand for godyl.
// It displays the application's default configuration settings.
package defaults

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/cli/flags"
	"github.com/idelchi/godyl/internal/config"
	"github.com/idelchi/godyl/internal/defaults"
	iutils "github.com/idelchi/godyl/internal/utils"
)

// Command encapsulates the defaults dump command with its associated configuration.
type Command struct {
	// Command is the defaults cobra.Command instance
	Command *cobra.Command
}

// Flags adds defaults-specific flags to the command.
func (cmd *Command) Flags() {
	// No specific flags for this command
}

// NewDefaultsCommand creates a Command for displaying default configuration settings.
func NewDefaultsCommand(cfg *config.Config, embedded config.Embedded) *Command {
	cmd := &cobra.Command{
		Use:   "defaults",
		Short: "Display default configuration settings",
		Args:  cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			return flags.ChainPreRun(cmd, nil)
		},
		RunE: func(_ *cobra.Command, _ []string) error {
			c, err := getDefaults(embedded)
			if err != nil {
				return err
			}

			iutils.Print(cfg.Dump.Format, c)

			return nil
		},
	}

	return &Command{
		Command: cmd,
	}
}

// NewCommand creates a cobra.Command instance for the defaults dump subcommand.
func NewCommand(cfg *config.Config, files config.Embedded) *cobra.Command {
	// Create the defaults command
	cmd := NewDefaultsCommand(cfg, files)

	// Add defaults-specific flags
	cmd.Flags()

	return cmd.Command
}

// getDefaults loads and returns the application's default settings.
func getDefaults(files config.Embedded) (*defaults.Defaults, error) {
	d := defaults.NewDefaults()

	if err := d.Unmarshal(files.Defaults); err != nil {
		return d, fmt.Errorf("setting defaults: %w", err)
	}

	return d, nil
}
