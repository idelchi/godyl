package remove

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

func NewRemoveCommand(cfg *config.Config) *Command {
	cmd := &cobra.Command{
		Use:               "remove",
		Short:             "Remove the cache",
		Long:              "Remove the cache.",
		Aliases:           []string{"rm"},
		Args:              cobra.NoArgs,
		PersistentPreRunE: common.KCreateSubcommandPreRunE("remove", nil, &cfg.Root.Show),
		RunE: func(cmd *cobra.Command, args []string) error {
			if cfg.Root.Show {
				return nil
			}

			cacheFile := cache.File(cfg.Root.Cache.Dir)

			if !cacheFile.Exists() {
				fmt.Println("Cache file doesn't exist")

				return nil
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
