package clean

import (
	"fmt"
	"time"

	"github.com/idelchi/godyl/internal/cache"
	"github.com/idelchi/godyl/internal/config"
	"github.com/idelchi/godyl/pkg/executable"
	"github.com/idelchi/godyl/pkg/logger"
	"github.com/idelchi/godyl/pkg/path/file"
	"github.com/idelchi/godyl/pkg/version"
)

// run executes the `cache clean` command.
func run(global config.Config) error {
	root := global

	cacheFile, err := cache.File(root.Cache.Dir)
	if err != nil {
		return err
	}

	lvl, err := logger.LevelString(root.LogLevel)
	if err != nil {
		return fmt.Errorf("parsing log level: %w", err)
	}

	log, err := logger.New(lvl)
	if err != nil {
		return fmt.Errorf("creating logger: %w", err)
	}

	cacheHandler := cache.New(cacheFile)
	if err = cacheHandler.Load(); err != nil {
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

		if !file.New(exe.String()).Exists() {
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

		if version.Equal(parsed, tool.Version.Version) {
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
}
