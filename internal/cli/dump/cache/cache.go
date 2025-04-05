package cache

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/cache"
	"github.com/idelchi/godyl/internal/cli/flags"
	"github.com/idelchi/godyl/internal/config"
	"github.com/idelchi/godyl/internal/tmp"
	iutils "github.com/idelchi/godyl/internal/utils"
	"github.com/idelchi/godyl/pkg/path/folder"
)

// Command encapsulates the tools dump command with its associated configuration.
type Command struct {
	// Command is the tools cobra.Command instance
	Command *cobra.Command
	// Config contains application configuration
	Config *config.Config
}

// Flags adds defaults-specific flags to the command.
func (cmd *Command) Flags() {
	cmd.Command.Flags().BoolP("file", "f", false, "Show the path to the cache file")
}

// NewCacheCommand creates a Command for displaying tools information.
func NewCacheCommand(cfg *config.Config) *Command {
	cmd := &cobra.Command{
		Use:   "cache [name]",
		Short: "Display cache information",
		Args:  cobra.MaximumNArgs(1),
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			return flags.ChainPreRun(cmd, &cfg.Dump.Cache)
		},
		RunE: func(_ *cobra.Command, args []string) error {
			var folder folder.Folder

			// TODO(Idelchi): This setting of flags to correct values should be centralized. Like flags.Defaults().
			switch {
			case cfg.Root.IsSet("cache-dir"):
				folder = cfg.Root.Cache.Dir
			default:
				folder = tmp.CacheDir()
			}

			if cfg.Dump.Cache.File {
				fmt.Println(cache.File(folder, cfg.Root.Cache.Type))

				return nil
			}

			var name string
			if len(args) > 0 {
				name = args[0]
			}

			c, err := getCache(folder, cfg.Root.Cache.Type, name)
			if err != nil {
				return err
			}

			iutils.Print(cfg.Dump.Format, c)

			return nil
		},
	}

	return &Command{
		Command: cmd,
		Config:  cfg,
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
func getCache(folder folder.Folder, cacheType string, name string) (any, error) {
	cache, err := cache.New(folder, cacheType)
	if err != nil {
		return nil, fmt.Errorf("failed to create cache: %w", err)
	}
	defer cache.Close()

	// err = cache.Save("go", "1.18.0")
	// if err != nil {
	// 	log.Fatalf("Failed to save item: %v", err)
	// }

	// err = cache.Save("nodejs", "16.13.1")
	// if err != nil {
	// 	log.Fatalf("Failed to save item: %v", err)
	// }

	var content any

	if name != "" {
		content, err = cache.Get(name)
	} else {
		content, err = cache.GetAll()
	}

	if err != nil {
		return nil, fmt.Errorf("failed to display cache: %w", err)
	}

	return content, nil
}
