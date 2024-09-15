package show

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/cache"
	"github.com/idelchi/godyl/internal/cli/common"
	"github.com/idelchi/godyl/internal/config"
	iutils "github.com/idelchi/godyl/internal/utils"
	"github.com/idelchi/godyl/pkg/path/folder"
)

// Command encapsulates the tools dump command with its associated configuration.
type Command struct {
	// Command is the tools cobra.Command instance
	Command *cobra.Command
}

func NewShowCommand(cfg *config.Config) *Command {
	cmd := &cobra.Command{
		Use:               "dump [name]",
		Short:             "Dump cache information",
		Args:              cobra.MaximumNArgs(1),
		PersistentPreRunE: common.KCreateSubcommandPreRunE("dump", &cfg.Cache.Dump, &cfg.Root.Show),
		RunE: func(cmd *cobra.Command, args []string) error {
			if cfg.Root.Show {
				return nil
			}

			var name string
			if len(args) > 0 {
				name = args[0]
			}

			c, err := getCache(cfg.Root.Cache.Dir, name)
			if err != nil {
				return err
			}

			iutils.Print(cfg.Cache.Dump.Format, c)

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
	cmd := NewShowCommand(cfg)

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
