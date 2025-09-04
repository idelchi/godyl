package auth

import (
	"fmt"

	"github.com/idelchi/godyl/internal/cli/core"
	"github.com/idelchi/godyl/internal/iutils"
	"github.com/idelchi/godyl/internal/tokenstore"
	"github.com/idelchi/godyl/pkg/pretty"
)

// run executes the `dump auth` command.
func run(input core.Input) error {
	cfg, _, context, _, _ := input.Unpack()

	tokens, _ := iutils.StructToKoanf(cfg.Tokens)

	var output any

	switch cfg.Keyring {
	case true:
		store := tokenstore.New()

		if ok, err := store.Available(); !ok {
			return err
		}

		values, err := store.GetAll(tokens.Keys()...)
		if err != nil {
			return err
		}

		if len(values) == 0 {
			fmt.Println("No tokens found in the keyring.")

			return nil
		}

		output = values

	case false:
		configuration := context.Config.Filtered(tokens.Keys()...)

		output = configuration.Map()
	}

	pretty.PrintYAML(output)

	return nil
}
