package path

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/cache"
	"github.com/idelchi/godyl/internal/cli/common"
	"github.com/idelchi/godyl/internal/config"
)

type Command struct {
	// Command is the tools cobra.Command instance
	Command *cobra.Command
}

func NewPathCommand(cfg *config.Config) *Command {
	cmd := &cobra.Command{
		Use:               "path",
		Short:             "Display the cache path",
		Args:              cobra.NoArgs,
		PersistentPreRunE: common.KCreateSubcommandPreRunE("path", nil, &cfg.Root.Show),
		Run: func(cmd *cobra.Command, args []string) {
			if cfg.Root.Show {
				return
			}

			cacheFile := cache.File(cfg.Root.Cache.Dir)

			if !cacheFile.Exists() {
				fmt.Println("Cache file doesn't exist")

				return
			}

			fmt.Println(cacheFile)
		},
	}

	return &Command{
		Command: cmd,
	}
}

func NewCommand(cfg *config.Config) *cobra.Command {
	// Create the tools command
	cmd := NewPathCommand(cfg)

	// Add tools-specific flags
	cmd.Flags()

	return cmd.Command
}
