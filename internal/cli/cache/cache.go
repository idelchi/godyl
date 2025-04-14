package cache

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/cache"
	"github.com/idelchi/godyl/internal/cli/flags"
	"github.com/idelchi/godyl/internal/config"
)

// Command encapsulates the Cache cobra command with its associated config and embedded files.
type Command struct {
	// Command is the Cache cobra.Command instance
	Command *cobra.Command
}

// Flags adds Cache-specific flags to the command.
func (cmd *Command) Flags() {
	flags.Cache(cmd.Command)
}

// NewCacheCommand creates a Command for updating the application to the latest version.
func NewCacheCommand(cfg *config.Config) *Command {
	cmd := &cobra.Command{
		Use:   "cache",
		Short: "Interact with the cache",
		Long:  "Interact with the cache. Displays the path.",
		Args:  cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			return flags.ChainPreRun(cmd, &cfg.Cache)
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			cacheFile := cache.File(cfg.Root.Cache.Dir)

			if !cfg.Cache.Delete || !cacheFile.Exists() {
				return nil
			}

			if err := cacheFile.Remove(); err != nil {
				return fmt.Errorf("removing cache: %w", err)
			}

			fmt.Println(cacheFile)

			return nil
		},
	}

	return &Command{
		Command: cmd,
	}
}

// NewCommand creates a cobra.Command instance containing the Cache command.
func NewCommand(cfg *config.Config) *cobra.Command {
	// Create the Cache command
	cmd := NewCacheCommand(cfg)

	// Add Cache-specific flags
	cmd.Flags()

	return cmd.Command
}
