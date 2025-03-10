// Package version provides the subcommand for printing the tool version.
package version

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewCommand(version string) *cobra.Command {
	return &cobra.Command{
		Use:               "version",
		Short:             "Print the version number of godyl",
		Long:              `All software has versions. This is godyl's`,
		Args:              cobra.MaximumNArgs(1),
		PersistentPreRunE: func(_ *cobra.Command, _ []string) error { return nil },
		Run: func(_ *cobra.Command, _ []string) {
			fmt.Println(version)
		},
	}
}
