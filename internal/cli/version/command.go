// Package version provides the subcommand for printing the tool version.
package version

import (
	"fmt"

	"github.com/idelchi/godyl/internal/cli/common"
	"github.com/idelchi/godyl/internal/config"
	"github.com/spf13/cobra"
)

// Command encapsulates the version cobra command with its version string.
type Command struct {
	// Command is the version cobra.Command instance
	Command *cobra.Command
}

// NewVersionCommand creates a Command for displaying the application version.
func NewVersionCommand(cfg *config.Config) *Command {
	cmd := &cobra.Command{
		Use:               "version",
		Short:             "Print the version number of godyl",
		Long:              `Print the version number of godyl along with available updates`,
		Aliases:           []string{"v"},
		Args:              cobra.NoArgs,
		PersistentPreRunE: common.KCreateSubcommandPreRunE("version", nil, &cfg.Root.Show),
		Run: func(cmd *cobra.Command, args []string) {
			if cfg.Root.Show {
				return
			}
			fmt.Println(cmd.Root().Version)
		},
	}

	return &Command{
		Command: cmd,
	}
}

// NewCommand creates a cobra.Command instance for the version command.
func NewCommand(cfg *config.Config) *cobra.Command {
	// Create the version command
	cmd := NewVersionCommand(cfg)

	// Add version-specific flags
	cmd.Flags()

	return cmd.Command
}
