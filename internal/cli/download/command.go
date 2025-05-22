// Package download implements the download command for godyl.
// It provides functionality to download and extract tools from various sources.
package download

import (
	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/cli/common"
	"github.com/idelchi/godyl/internal/config"
	"github.com/idelchi/godyl/internal/config/download"
)

func Command(global *config.Config, local any, embedded *common.Embedded) *cobra.Command {
	// Create the download command
	cmd := &cobra.Command{
		Use:     "download [tool]",
		Aliases: []string{"dl", "x"},
		Short:   "Download and extract tools",
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Exit early if the command is run with `--show/-s` flag.
			if common.ExitOnShow(global.Root.ShowFunc) {
				return nil
			}

			return run(*global, *embedded, args...)
		},
	}

	common.SetSubcommandDefaults(cmd, local, global.Root.ShowFunc)

	download.Flags(cmd)

	return cmd
}
