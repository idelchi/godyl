package update

import "github.com/spf13/cobra"

func Flags(cmd *cobra.Command) {
	cmd.Flags().String("version", "", "Version of the tool to install. Empty means latest.")
	cmd.Flags().Bool("pre", false, "Enable pre-release versions")
	cmd.Flags().Bool("force", false, "Force update even if the current version is the latest")
	cmd.Flags().Bool("check", false, "Check for updates only")
	cmd.Flags().Bool("cleanup", false, "Cleanup after update (only valid for Windows)")
}
