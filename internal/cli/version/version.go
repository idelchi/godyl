// Package version provides the subcommand for printing the tool version.
package version

import (
	"fmt"

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
		Long:  `All software has versions. This is godyl's`,
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, _ []string) {
			fmt.Println(cmd.Root().Version)
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
