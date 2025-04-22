package path

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/cache"
	"github.com/idelchi/godyl/internal/cli/flags"
	"github.com/idelchi/godyl/internal/config"
)

type Command struct {
	// Command is the tools cobra.Command instance
	Command *cobra.Command
}

// Flags adds defaults-specific flags to the command.
func (cmd *Command) Flags() {
}

func NewPathCommand(cfg *config.Config) *Command {
	cmd := &cobra.Command{
		Use:   "path",
		Short: "Display the cache path",
		Args:  cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			return flags.ChainPreRun(cmd, nil)
		},
		Run: func(_ *cobra.Command, _ []string) {
			cacheFile := cache.File(cfg.Root.Cache.Dir)

			if !cacheFile.Exists() {
				fmt.Println("Cache file doesn't exist")
			} else {
				fmt.Println(cacheFile)
			}
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
