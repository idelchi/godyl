package install

import (
	"github.com/idelchi/godyl/internal/tools/strategy"
)

func (cmd *Command) Flags() {
	cmd.Command.Flags().StringP("output", "o", "./bin", "Output path for the downloaded tools")
	cmd.Command.Flags().String("os", "", "Operating system to install the tools for")
	cmd.Command.Flags().String("arch", "", "Architecture to install the tools for")

	cmd.Command.Flags().StringSliceP("tags", "t", []string{"!native"}, "Tags to filter tools by. Prefix with '!' to exclude")
	cmd.Command.Flags().String("strategy", strategy.Sync.String(), "Strategy to use for updating tools (none, sync, force)")
}
