package update

// Update adds update-related command flags to the provided cobra command.
// These flags control the self-update.
func (cmd *Command) Flags() {
	cmd.Command.Flags().String("version", "", "Version of the tool to install. Empty means latest.")
	cmd.Command.Flags().Bool("pre", false, "Enable pre-release versions")
	cmd.Command.Flags().Bool("force", false, "Force update even if the current version is the latest")
	cmd.Command.Flags().Bool("check", false, "Check for updates only")
	cmd.Command.Flags().Bool("cleanup", false, "Cleanup after update (only valid for Windows)")
}
