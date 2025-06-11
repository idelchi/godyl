package store

import (
	"fmt"

	"github.com/idelchi/godyl/internal/cli/common"
	"github.com/idelchi/godyl/internal/ierrors"
	"github.com/idelchi/godyl/internal/iutils"
	"github.com/idelchi/godyl/internal/tokenstore"
	"github.com/idelchi/godyl/pkg/editor"
)

// run executes the `auth store` command.
func run(input common.Input) error {
	cfg, _, _, _, args := input.Unpack()

	kTokens, _ := iutils.StructToKoanf(cfg.Tokens)
	tokens := kTokens.Map()

	selected := make(map[string]any, len(args))

	if len(args) == 0 {
		selected = tokens
	}

	for _, arg := range args {
		selected[arg] = tokens[arg]
	}

	for key := range selected {
		if !cfg.IsSet(key) {
			delete(selected, key)
		}
	}

	if len(selected) == 0 {
		return fmt.Errorf("%w: no token values provided", ierrors.ErrUsage)
	}

	logger, err := common.SetupLogger(cfg.LogLevel)
	if err != nil {
		return err
	}

	switch cfg.Keyring {
	case true:
		store := tokenstore.New()

		if ok, err := store.Available(); !ok {
			return err
		}

		err := store.SetAll(selected)
		if err != nil {
			return err
		}

		logger.Info("tokens successfully set in the keyring.")
	case false:
		if err := editor.New(cfg.ConfigFile).Merge(selected); err != nil {
			return err
		}

		logger.Info("tokens successfully set in the configuration file.")
	}

	return nil
}
