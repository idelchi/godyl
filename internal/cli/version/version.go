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
	// Version is the application version string
	Version string
}

// Flags adds version-specific flags to the command.
func (cmd *Command) Flags() {
	// No specific flags for this command
}

// NewVersionCommand creates a Command for displaying the application version.
func NewVersionCommand(version string) *Command {
	cmd := &cobra.Command{
		Use:               "version",
		Short:             "Print the version number of godyl",
		Long:              `All software has versions. This is godyl's`,
		Args:              cobra.MaximumNArgs(1),
		PersistentPreRunE: func(_ *cobra.Command, _ []string) error { return nil },
		Run: func(_ *cobra.Command, _ []string) {
			fmt.Println(version)
		},
	}

	return &Command{
		Command: cmd,
		Version: version,
	}
}

// NewCommand creates a cobra.Command instance for the version command.
func NewCommand(version string) *cobra.Command {
	// Create the version command
	cmd := NewVersionCommand(version)

	// Add version-specific flags
	cmd.Flags()

	return cmd.Command
}
