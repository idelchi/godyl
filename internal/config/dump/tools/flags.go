package tools

import "github.com/spf13/cobra"

func Flags(cmd *cobra.Command) {
	cmd.Flags().BoolP("full", "f", false, "Show the tools in full syntax")
	cmd.Flags().BoolP("embedded", "e", false, "Show the embedded tools")
	cmd.Flags().StringSliceP("tags", "t", []string{""}, "Tags to filter tools by. Prefix with '!' to exclude")
}
