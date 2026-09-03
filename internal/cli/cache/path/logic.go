package path

import (
	"fmt"

	"github.com/idelchi/godyl/internal/cli/core"
	"github.com/idelchi/godyl/internal/data"
)

// run prints the cache file derived from the configured cache directory.
func run(input core.Input) error {
	cfg, _, _, _, _ := input.Unpack()

	fmt.Println(data.CacheFile(cfg.Cache.Dir))

	return nil
}
