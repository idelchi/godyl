package version

import (
	"fmt"

	"github.com/idelchi/godyl/internal/cli/core"
)

// run prints the version stored on the root command.
func run(input core.Input) {
	_, _, _, cmd, _ := input.Unpack()

	fmt.Println(cmd.Root().Version)
}
