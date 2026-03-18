package download

import "github.com/spf13/cobra"

// Flags configures the command-line flags for the download command.
func Flags(cmd *cobra.Command) {
	cmd.Flags().SortFlags = false

	cmd.Flags().StringP("output", "o", "./bin", "output path for the downloaded tools")
	cmd.Flags().String("source", "github", "source from which to download the tools (github, gitlab, url)")
	cmd.Flags().String("os", "", "override the OS to match")
	cmd.Flags().String("arch", "", "override the architecture to match")
	cmd.Flags().StringSlice("hints", []string{""}, "hint patterns with weight 1 and type glob")
	cmd.Flags().String("version", "", "version to download, leave empty for latest")
	cmd.Flags().Bool("dry", false, "dry run, show what would be done without downloading")
	cmd.Flags().Bool("pre", false, "consider pre-releases when downloading tools")
}
