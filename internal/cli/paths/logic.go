package paths

import (
	"fmt"

	"github.com/idelchi/godyl/internal/cli/core"
	"github.com/idelchi/godyl/internal/data"
)

// run executes the `paths` command.
func run(input core.Input) error {
	cfg, _, _, _, _ := input.Unpack()

	fmt.Printf("config path: %s\n", cfg.ConfigFile)
	fmt.Printf("cache path: %s\n", data.CacheFile(cfg.Cache.Dir))
	fmt.Printf("go path: %s\n", data.GoDir())
	fmt.Printf("temp download path: %s\n", data.DownloadDir().Path()+"-"+data.Prefix())

	return nil
}
