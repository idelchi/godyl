package path

import (
	"fmt"

	"github.com/idelchi/godyl/internal/cli/common"
)

// run executes the `config path` command.
func run(input common.Input) error {
	cfg, _, _, _, _ := input.Unpack()

	fmt.Println(cfg.ConfigFile.Absolute())

	return nil
}
