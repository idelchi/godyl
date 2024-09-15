package dump

import "github.com/spf13/cobra"

func flags(cmd *cobra.Command) {
	cmd.Flags().StringP("format", "f", "yaml", "Output format (json or yaml)")
}
