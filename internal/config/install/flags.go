package install

import (
	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/tools/strategy"
)

// Flags adds the flags for the `godyl install` command to the provided Cobra command.
func Flags(cmd *cobra.Command) {
	cmd.Flags().SortFlags = false

	cmd.Flags().StringP("output", "o", "./bin", "output path for the downloaded tools")
	cmd.Flags().String("source", "github", "source from which to install the tools (github, gitlab, url, go, none)")
	cmd.Flags().String("os", "", "override the OS to match")
	cmd.Flags().String("arch", "", "override the architecture to match")

	cmd.Flags().StringSliceP("tags", "t", []string{"!native"}, "tags to filter tools by, prefix with '!' to exclude")
	cmd.Flags().
		String("strategy", strategy.Sync.String(), "strategy to use for updating tools (none, sync, existing, force)")

	cmd.Flags().Bool("dry", false, "dry run, show what would be done without downloading")
	cmd.Flags().Bool("pre", false, "consider pre-releases when installing tools")
}
