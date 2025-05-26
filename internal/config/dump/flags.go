package dump

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Flags adds the flags for the `godyl dump` command to the provided Cobra command.
func Flags(cmd *cobra.Command) {
	cmd.Flags().StringP("format", "f", "yaml", fmt.Sprintf("Output format (%v)", "[json yaml]"))
}
