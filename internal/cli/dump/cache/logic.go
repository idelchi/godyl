package cache

import (
	"fmt"

	"github.com/idelchi/godyl/internal/cache"
	"github.com/idelchi/godyl/internal/config"
	"github.com/idelchi/godyl/internal/iutils"
	"github.com/idelchi/godyl/pkg/path/file"
)

// run executes the `cache dump` command.
func run(cfg config.Config, args []string) error {
	cacheFile, err := cache.File(cfg.Root.Cache.Dir)
	if err != nil {
		return err
	}

	var name string
	if len(args) > 0 {
		name = args[0]
	}

	cache, err := getCache(cacheFile, name)
	if err != nil {
		return err
	}

	iutils.Print(cfg.Dump.Format, cache)

	return nil
}

// getCache retrieves the cache from the specified folder and cache type and returns the content.
func getCache(file file.File, name string) (content any, err error) {
	cache := cache.New(file)
	if err = cache.Load(); err != nil {
		return nil, fmt.Errorf("loading cache file %q: %w", file, err)
	}

	if name != "" {
		content, err = cache.GetByName(name)
	} else {
		content, err = cache.GetAll()
	}

	if err != nil {
		return nil, fmt.Errorf("failed to display cache: %w", err)
	}

	return content, nil
}
