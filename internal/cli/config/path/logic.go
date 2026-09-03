package path

import (
	"fmt"

	"github.com/idelchi/godyl/internal/cli/core"
)

// run prints the resolved configuration file path.
func run(input core.Input) error {
	cfg, _, _, _, _ := input.Unpack()

	fmt.Println(cfg.ConfigFile.Absolute())

	return nil
}
