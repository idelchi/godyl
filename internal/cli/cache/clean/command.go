package clean

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/cache"
	"github.com/idelchi/godyl/internal/cli/common"
	"github.com/idelchi/godyl/internal/config"
	"github.com/idelchi/godyl/pkg/executable"
	"github.com/idelchi/godyl/pkg/logger"
	"github.com/idelchi/godyl/pkg/version"
)

type Command struct {
	// Command is the Cache cobra.Command instance
	Command *cobra.Command
}

func NewCleanCommand(cfg *config.Config, embedded *config.Embedded) *Command {
	cmd := &cobra.Command{
		Use:   "clean",
		Short: "Interact with the cache",
		Long:  "Interact with the cache.",
		Args:  cobra.ArbitraryArgs,

		PersistentPreRunE: common.KCreateSubcommandPreRunE("clean", nil, &cfg.Root.Show),
		RunE: func(cmd *cobra.Command, args []string) error {
			if cfg.Root.Show {
				return nil
			}
			lvl, err := logger.LevelString(cfg.Root.LogLevel)
			if err != nil {
				return fmt.Errorf("parsing log level: %w", err)
			}

			log := logger.New(lvl)

			cacheHandler := cache.New(cfg.Root.Cache.Dir)
			if err := cacheHandler.Load(); err != nil {
				return fmt.Errorf("loading cache: %w", err)
			}

			tools, err := cacheHandler.GetAll()
			if err != nil {
				return fmt.Errorf("getting tools from cache: %w", err)
			}

			for _, tool := range tools {
				// Parse the version of the existing tool.
				exe := executable.New(tool.Path)
				commands := tool.Version

				if !exe.Exists() {
					if err := cacheHandler.Delete(tool.ID); err != nil {
						log.Warnf("failed to delete cache for id %q: %v", tool.ID, err)
					} else {
						log.Warnf("cache deleted for %q: executable %q has been removed from system", tool.Name, tool.Path)
					}

					continue
				}

				// Check if we have commands to determine version
				if commands.Commands == nil {
					continue
				}

				// Parse version using available commands
				parser := &executable.Parser{
					Patterns: *commands.Patterns,
					Commands: *commands.Commands,
				}

				parsed, err := exe.Parse(parser)
				if err != nil {
					log.Warnf("failed to parse version for %q: %v", tool.Name, err)

					continue
				}

				if version.Compare(parsed, tool.Version.Version) {
					continue
				}

				tool.Version.Version = parsed
				tool.Updated = time.Now()

				if err := cacheHandler.Save(tool); err != nil {
					log.Warnf("failed to save cache for %q: %v", tool.Name, err)
				} else {
					log.Infof("cache updated for %q: version %q parsed", tool.Name, tool.Version.Version)
				}
			}

			if !cacheHandler.Touched() {
				log.Info("no changes necessary")
			}

			return nil
		},
	}

	return &Command{
		Command: cmd,
	}
}

func NewCommand(cfg *config.Config, files *config.Embedded) *cobra.Command {
	cmd := NewCleanCommand(cfg, files)

	// Add Cache-specific flags
	cmd.Flags()

	return cmd.Command
}
