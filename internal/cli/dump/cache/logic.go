package cache

import (
	"fmt"

	"github.com/idelchi/godyl/internal/cache"
	"github.com/idelchi/godyl/internal/cli/core"
	"github.com/idelchi/godyl/internal/iutils"
	"github.com/idelchi/godyl/internal/tmp"
	"github.com/idelchi/godyl/pkg/path/file"
)

// run executes the `dump cache` command.
func run(input core.Input) error {
	cfg, _, _, _, args := input.Unpack()

	cacheFile := tmp.CacheFile(cfg.Cache.Dir)

	if !cacheFile.Exists() {
		return fmt.Errorf("cache file %q does not exist", cacheFile)
	}

	var name string

	if len(args) > 0 {
		name = args[0]
	}

	cache, err := getCache(cacheFile, name)
	if err != nil {
		return err
	}

	iutils.Print(iutils.YAML, cache)

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
		content, err = cache.Get()
	}

	if err != nil {
		return nil, fmt.Errorf("failed to display cache: %w", err)
	}

	return content, nil
}
