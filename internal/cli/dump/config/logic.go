package config

import (
	"fmt"
	"strings"

	"github.com/idelchi/godyl/internal/cli/core"
	"github.com/idelchi/godyl/internal/iutils"
	"github.com/idelchi/godyl/pkg/pretty"
)

// run executes the `dump config` command.
func run(input core.Input) error {
	cfg, _, _, _, args := input.Unpack()

	if len(args) == 0 {
		pretty.PrintYAMLMasked(cfg)

		return nil
	}

	configuration, err := iutils.StructToKoanf(cfg)
	if err != nil {
		return err
	}

	for _, key := range args {
		val := configuration.Get(key)
		if val == nil {
			return fmt.Errorf("value %q not found in config", key)
		}

		if len(args) > 1 {
			fmt.Printf(" ---- %s ----\n", key)
		}

		if strings.HasSuffix(key, "-token") {
			val = "********"
		}

		iutils.Print(iutils.YAML, val)

		if len(args) > 1 {
			fmt.Println()
		}
	}

	return nil
}
