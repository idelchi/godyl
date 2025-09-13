package status

import "github.com/spf13/cobra"

// Flags adds the flags for the `godyl status` command to the provided Cobra command.
func Flags(cmd *cobra.Command) {
	cmd.Flags().SortFlags = false

	cmd.Flags().StringP("output", "o", "./bin", "Output path for the downloaded tools")
	cmd.Flags().StringSliceP("tags", "t", []string{"!native"}, "Tags to filter tools by. Prefix with '!' to exclude")
}
