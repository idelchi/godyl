package auth

import (
	"fmt"

	"github.com/idelchi/godyl/internal/cli/common"
	"github.com/idelchi/godyl/internal/iutils"
	"github.com/idelchi/godyl/internal/tokenstore"
	"github.com/idelchi/godyl/pkg/pretty"
)

// run executes the `dump auth` command.
func run(input common.Input) error {
	cfg, _, context, _, _ := input.Unpack()

	tokens, _ := iutils.StructToKoanf(cfg.Tokens)

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

		pretty.PrintYAML(values)
	case false:
		configuration := context.Config.Filtered(tokens.Keys()...)

		pretty.PrintYAML(configuration.Map())
	}

	return nil
}
