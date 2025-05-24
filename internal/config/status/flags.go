package status

import "github.com/spf13/cobra"

func Flags(cmd *cobra.Command) {
	cmd.Flags().StringSliceP("tags", "t", []string{"!native"}, "Tags to filter tools by. Prefix with '!' to exclude")
}
