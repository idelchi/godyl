package paths

import (
	"github.com/idelchi/godyl/internal/cli/core"
	"github.com/idelchi/godyl/internal/data"
	"github.com/idelchi/godyl/pkg/pretty"
)

// run executes the `paths` command.
func run(input core.Input) error {
	cfg, _, _, _, _ := input.Unpack()

	paths := struct {
		Config string
		Cache  string
		Go     string
		Temp   string
	}{
		Config: cfg.ConfigFile.Path(),
		Cache:  data.CacheFile(cfg.Cache.Dir).Path(),
		Go:     data.GoDir().Path(),
		Temp:   data.DownloadDir().Path() + "-" + data.Prefix(),
	}

	pretty.PrintYAML(paths)

	return nil
}
