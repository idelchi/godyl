package install

import (
	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/tools/strategy"
)

// Flags adds the flags for the `godyl install` command to the provided Cobra command.
func Flags(cmd *cobra.Command) {
	cmd.Flags().SortFlags = false

	cmd.Flags().StringP("output", "o", "./bin", "Output path for the downloaded tools")
	cmd.Flags().String("source", "github", "source from which to install the tools (github, gitlab, url, go, none)")
	cmd.Flags().String("os", "", "Operating system to install the tools for")
	cmd.Flags().String("arch", "", "Architecture to install the tools for")

	cmd.Flags().
		StringSliceP("tags", "t", []string{"!native"}, "Tags to filter tools by. Prefix with '!' to exclude")
	cmd.Flags().
		String("strategy", strategy.Sync.String(), "Strategy to use for updating tools (none, sync, force)")

	cmd.Flags().Bool("dry", false, "Perform a dry run without downloading tools")
}
