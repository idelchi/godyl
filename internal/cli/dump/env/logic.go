package env

import (
	"strings"

	"github.com/idelchi/godyl/internal/iutils"
	"github.com/idelchi/godyl/pkg/env"
)

func run(dotenv env.Env) error {
	env := env.FromEnv()
	env.Merge(dotenv)

	startsWithGodyl := func(k, _ string) bool {
		return strings.HasPrefix(k, "GODYL_")
	}

	iutils.Print(iutils.ENV, env.GetWithPredicates(startsWithGodyl))

	return nil
}
