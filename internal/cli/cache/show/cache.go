package show

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/cache"
	"github.com/idelchi/godyl/internal/cli/flags"
	"github.com/idelchi/godyl/internal/config"
	iutils "github.com/idelchi/godyl/internal/utils"
	"github.com/idelchi/godyl/pkg/path/folder"
)

// Command encapsulates the tools dump command with its associated configuration.
type Command struct {
	// Command is the tools cobra.Command instance
	Command *cobra.Command
}

func (cmd *Command) Flags() {
	cmd.Command.Flags().StringP("format", "f", "yaml", "Output format (json or yaml)")
}

// NewCacheCommand creates a Command for displaying tools information.
func NewCacheCommand(cfg *config.Config) *Command {
	cmd := &cobra.Command{
		Use:   "show [name]",
		Short: "Display cache information",
		Args:  cobra.MaximumNArgs(1),
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			return flags.ChainPreRun(cmd, &cfg.Cache)
		},
		RunE: func(_ *cobra.Command, args []string) error {
			var name string
			if len(args) > 0 {
				name = args[0]
			}

			c, err := getCache(cfg.Root.Cache.Dir, name)
			if err != nil {
				return err
			}

			iutils.Print(cfg.Cache.Format, c)

			return nil
		},
	}

	return &Command{
		Command: cmd,
	}
}

// NewCommand creates a cobra.Command instance for the tools dump subcommand.
func NewCommand(cfg *config.Config) *cobra.Command {
	// Create the tools command
	cmd := NewCacheCommand(cfg)

	// Add tools-specific flags
	cmd.Flags()

	return cmd.Command
}

// getCache retrieves the cache from the specified folder and cache type and returns the content.
func getCache(folder folder.Folder, name string) (content any, err error) {
	cache := cache.New(folder)
	if err = cache.Load(); err != nil {
		return nil, fmt.Errorf("failed to load cache: %w", err)
	}

	if name != "" {
		content, err = cache.GetByProperty("name", name)
	} else {
		content, err = cache.GetAll()
	}

	if err != nil {
		return nil, fmt.Errorf("failed to display cache: %w", err)
	}

	return content, nil
}
