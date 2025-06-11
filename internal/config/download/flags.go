package download

import "github.com/spf13/cobra"

func Flags(cmd *cobra.Command) {
	cmd.Flags().SortFlags = false

	cmd.Flags().StringP("output", "o", "./bin", "output path for the downloaded tools")
	cmd.Flags().String("source", "github", "source from which to install the tools (github, gitlab, url, go, none)")
	cmd.Flags().String("os", "", "operating system to install the tools for")
	cmd.Flags().String("arch", "", "architecture to install the tools for")
	cmd.Flags().StringSlice("hints", []string{""}, "hints to use for tool resolution")
	cmd.Flags().String("version", "", "version of the tool to install, leave empty for latest")
	cmd.Flags().Bool("dry", false, "perform a dry run without downloading tools")
}
