package version

import (
	"fmt"

	"github.com/idelchi/godyl/internal/cli/core"
)

// run executes the `version` command.
func run(input core.Input) {
	_, _, _, cmd, _ := input.Unpack()

	fmt.Println(cmd.Root().Version)
}
