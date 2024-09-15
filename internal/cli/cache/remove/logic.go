package remove

import (
	"errors"
	"fmt"

	"github.com/idelchi/godyl/internal/cache"
	"github.com/idelchi/godyl/internal/cli/common"
	"github.com/idelchi/godyl/internal/tmp"
)

// run executes the `cache remove` command.
func run(input common.Input) error {
	cfg, _, _, _, args := input.Unpack()

	logger, err := common.SetupLogger(cfg.LogLevel)
	if err != nil {
		return err
	}

	cacheFile := tmp.CacheFile(cfg.Cache.Dir)
	if !cacheFile.Exists() {
		return fmt.Errorf("cache file %q does not exist", cacheFile)
	}

	c := cache.New(cacheFile)

	if err := c.Load(); err != nil {
		return err
	}

	if c.IsEmpty() {
		logger.Info("Cache is already empty.")

		return nil
	}

	switch len(args) {
	case 0:
		if err := c.Delete(); err != nil {
			return fmt.Errorf("removing cache entries: %w", err)
		}

		logger.Info("All cache entries have been removed.")
	default:
		for _, name := range args {
			if err := c.DeleteByName(name); errors.Is(err, cache.ErrItemNotFound) {
				logger.Warnf("Cache entry %q does not exist.", name)
			} else {
				logger.Infof("Removing cache entry %q.", name)
			}
		}
	}

	return nil
}
