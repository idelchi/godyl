package tools

import "github.com/spf13/cobra"

// Flags configures the command-line flags for the dump tools command.
func Flags(cmd *cobra.Command) {
	cmd.Flags().SortFlags = false

	cmd.Flags().BoolP("embedded", "e", false, "Show the embedded tools")
	cmd.Flags().StringSliceP("tags", "t", []string{""}, "Tags to filter tools by. Prefix with '!' to exclude")
	cmd.Flags().BoolP("full", "f", false, "Show the tools in full syntax")
}
