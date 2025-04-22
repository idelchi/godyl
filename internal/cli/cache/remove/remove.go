package remove

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

func NewRemoveCommand(cfg *config.Config) *Command {
	cmd := &cobra.Command{
		Use:     "remove",
		Short:   "Remove the cache",
		Long:    "Remove the cache.",
		Aliases: []string{"rm"},
		Args:    cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			return flags.ChainPreRun(cmd, nil)
		},
		RunE: func(_ *cobra.Command, _ []string) error {
			cacheFile := cache.File(cfg.Root.Cache.Dir)

			if !cacheFile.Exists() {
				fmt.Println("Cache file doesn't exist")
			}

			if err := cacheFile.Remove(); err != nil {
				return fmt.Errorf("removing cache: %w", err)
			}

			fmt.Printf("Cache file %q removed\n", cacheFile)

			return nil
		},
	}

	return &Command{
		Command: cmd,
	}
}

func NewCommand(cfg *config.Config) *cobra.Command {
	// Create the tools command
	cmd := NewRemoveCommand(cfg)

	// Add tools-specific flags
	cmd.Flags()

	return cmd.Command
}
