package update

import "github.com/spf13/cobra"

// Flags adds the flags for the `godyl update` command to the provided Cobra command.
func Flags(cmd *cobra.Command) {
	cmd.Flags().SortFlags = false

	cmd.Flags().String("version", "", "Version of the tool to install. Empty means latest.")
	cmd.Flags().Bool("pre", false, "Enable pre-release versions")
	cmd.Flags().Bool("check", false, "Check for updates only")
	cmd.Flags().Bool("cleanup", false, "Cleanup after update (only valid for Windows)")
	cmd.Flags().Bool("force", false, "Force update even if the current version is the latest")
}
