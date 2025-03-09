package version

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewCommand(version string) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "version",
		Short:             "Print the version number of godyl",
		Long:              `All software has versions. This is godyl's`,
		Args:              cobra.MaximumNArgs(1),
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error { return nil },
		Run: func(_ *cobra.Command, args []string) {
			fmt.Println(version)
		},
	}

	return cmd
}
