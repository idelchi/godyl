package env

import (
	"fmt"
	"strings"

	"github.com/idelchi/godyl/internal/cli/common"
	"github.com/idelchi/godyl/internal/iutils"
)

// run executes the `dump env` command.
func run(input common.Input) error {
	_, _, context, _, _ := input.Unpack()

	startsWithGodyl := func(k, _ string) bool {
		return strings.HasPrefix(k, "GODYL_")
	}

	dotenv := *context.DotEnv
	env := context.Env.GetWithPredicates(startsWithGodyl)

	if len(dotenv) > 0 {
		fmt.Println("***** from .env file(s) *****")
		iutils.Print(iutils.ENV, dotenv)
	}

	if len(env) > 0 {
		fmt.Println("***** from environment variables *****")
		iutils.Print(iutils.ENV, env.GetWithPredicates(startsWithGodyl))
	}

	return nil
}
