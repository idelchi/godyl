// Package version provides the subcommand for printing the tool version.
package version

import (
	"fmt"

	"github.com/idelchi/godyl/internal/core/updater"
	"github.com/spf13/cobra"
)

// Command encapsulates the version cobra command with its version string.
type Command struct {
	// Command is the version cobra.Command instance
	Command *cobra.Command
}

// Flags adds version-specific flags to the command.
func (cmd *Command) Flags() {
	// No specific flags for this command
}

// NewVersionCommand creates a Command for displaying the application version.
func NewVersionCommand() *Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version number of godyl",
		Long:  `Print the version number of godyl along with available updates`,
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, _ []string) {
			current := cmd.Root().Version

			fmt.Println(current)

			all, err := updater.AllVersions()
			if err == nil {
				fmt.Printf("latest versions available: %q\n", all)

				fmt.Printf("Install with\n\n  godyl upgrade --version %s\n", all[0])
			}
		},
	}

	return &Command{
		Command: cmd,
	}
}

// NewCommand creates a cobra.Command instance for the version command.
func NewCommand() *cobra.Command {
	// Create the version command
	cmd := NewVersionCommand()

	// Add version-specific flags
	cmd.Flags()

	return cmd.Command
}
