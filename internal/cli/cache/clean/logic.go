package clean

import (
	"fmt"
	"time"

	"github.com/idelchi/godyl/internal/cache"
	"github.com/idelchi/godyl/internal/cli/core"
	"github.com/idelchi/godyl/internal/config/root"
	"github.com/idelchi/godyl/internal/tmp"
	"github.com/idelchi/godyl/pkg/executable"
	"github.com/idelchi/godyl/pkg/logger"
	"github.com/idelchi/godyl/pkg/path/file"
	"github.com/idelchi/godyl/pkg/version"
)

// run executes the `cache clean` command.
func run(input core.Input) error {
	cfg, _, _, _, _ := input.Unpack()

	logger, cacheHandler, err := setup(cfg)
	if err != nil {
		return err
	}

	tools, err := cacheHandler.Get()
	if err != nil {
		return fmt.Errorf("getting tools from cache: %w", err)
	}

	for _, tool := range tools {
		cleanTool(cacheHandler, tool, logger)
	}

	if !cacheHandler.Touched() {
		logger.Info("no changes necessary")
	}

	return nil
}

// setup initializes the logger and cache handler.
func setup(cfg *root.Config) (*logger.Logger, *cache.Cache, error) {
	logger, err := core.SetupLogger(cfg.LogLevel)
	if err != nil {
		return nil, nil, err
	}

	cacheFile := tmp.CacheFile(cfg.Cache.Dir)

	if !cacheFile.Exists() {
		return logger, nil, fmt.Errorf("cache file %q does not exist", cacheFile)
	}

	cacheHandler := cache.New(cacheFile)
	if err = cacheHandler.Load(); err != nil {
		return nil, nil, fmt.Errorf("loading cache: %w", err)
	}

	return logger, cacheHandler, nil
}

// cleanTool removes missing tools from cache or updates version for existing tools.
func cleanTool(cacheHandler *cache.Cache, tool *cache.Item, logger *logger.Logger) {
	exe := executable.New(tool.Path)

	if !file.New(exe.String()).Exists() {
		if err := cacheHandler.Delete(tool.ID); err != nil {
			logger.Warnf("failed to delete cache for id %q: %v", tool.ID, err)
		} else {
			logger.Warnf("cache deleted for %q: executable %q has been removed from system", tool.Name, tool.Path)
		}

		return
	}

	if tool.Version.Commands == nil {
		return
	}

	parser := &executable.Parser{
		Patterns: *tool.Version.Patterns,
		Commands: *tool.Version.Commands,
	}

	parsed, err := exe.Parse(parser)
	if err != nil {
		logger.Warnf("failed to parse version for %q: %v", tool.Name, err)

		return
	}

	if version.Equal(parsed, tool.Version.Version) {
		return
	}

	tool.Version.Version = parsed
	tool.Updated = time.Now()

	if err := cacheHandler.Add(tool); err != nil {
		logger.Warnf("failed to save cache for %q: %v", tool.Name, err)
	} else {
		logger.Infof("cache updated for %q: version %q parsed", tool.Name, tool.Version.Version)
	}
}
