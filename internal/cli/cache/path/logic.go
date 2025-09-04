package path

import (
	"fmt"

	"github.com/idelchi/godyl/internal/cli/core"
	"github.com/idelchi/godyl/internal/tmp"
)

// run executes the `cache path` command.
func run(input core.Input) error {
	cfg, _, _, _, _ := input.Unpack()

	fmt.Println(tmp.CacheFile(cfg.Cache.Dir))

	return nil
}
