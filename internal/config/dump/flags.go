package dump

import "github.com/spf13/cobra"

func Flags(cmd *cobra.Command) {
	cmd.Flags().StringP("format", "f", "yaml", "Output format (json or yaml)")
}
