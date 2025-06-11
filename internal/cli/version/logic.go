package version

import (
	"fmt"

	"github.com/idelchi/godyl/internal/cli/common"
)

// run executes the `version` command.
func run(input common.Input) {
	_, _, _, cmd, _ := input.Unpack()

	fmt.Println(cmd.Root().Version)
}
