package status

import (
	"fmt"
	"maps"
	"slices"

	"github.com/idelchi/godyl/internal/cli/core"
	"github.com/idelchi/godyl/internal/iutils"
)

// run executes the `auth status` command.
func run(input core.Input) error {
	cfg, _, _, _, _ := input.Unpack()

	kTokens, _ := iutils.StructToKoanf(cfg.Tokens)

	tokens := kTokens.Map()

	if !slices.ContainsFunc(slices.Collect(maps.Values(tokens)), func(v any) bool {
		return v != ""
	}) {
		fmt.Println("No authentication tokens are set in the current configuration.")

		return nil
	}

	for key, value := range tokens {
		set := "set"

		if value == "" {
			set = "unset"
		}

		fmt.Printf("%s: %s\n", key, set)
	}

	return nil
}
