package status

import "github.com/spf13/cobra"

// Flags adds the flags for the `godyl status` command to the provided Cobra command.
func Flags(cmd *cobra.Command) {
	cmd.Flags().StringSliceP("tags", "t", []string{"!native"}, "Tags to filter tools by. Prefix with '!' to exclude")
}
