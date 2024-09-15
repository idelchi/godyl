package remove

import (
	"errors"

	"github.com/idelchi/godyl/internal/cli/common"
	"github.com/idelchi/godyl/internal/iutils"
	"github.com/idelchi/godyl/internal/tokenstore"
	"github.com/idelchi/godyl/pkg/editor"
)

// run executes the `auth remove` command.
func run(input common.Input) error {
	cfg, _, context, _, args := input.Unpack()

	logger, err := common.SetupLogger(cfg.LogLevel)
	if err != nil {
		return err
	}

	var errs []error

	switch cfg.Keyring {
	case true:
		store := tokenstore.New()
		if ok, err := store.Available(); !ok {
			return err
		}

		if len(args) == 0 {
			if err := store.Delete(); err != nil {
				return err
			}

			logger.Info("All tokens successfully deleted from the keyring.")
		}

		for _, key := range args {
			if err := store.Delete(key); err != nil {
				logger.Warnf("%s: secret not found in keyring", key)

				continue
			}

			logger.Infof("Token %q successfully deleted from the keyring.", key)
		}
	case false:
		configuration := context.Config

		tokens, _ := iutils.StructToKoanf(cfg.Tokens)

		keys := tokens.Keys()
		if len(args) > 0 {
			keys = args
		}

		for _, key := range keys {
			if !configuration.Exists(key) {
				logger.Warnf("%s: secret not found in configuration file %q", key, cfg.ConfigFile)

				continue
			}

			configuration.Delete(key)

			if err := editor.New(cfg.ConfigFile).Write(configuration.Map()); err != nil {
				return err
			}

			logger.Infof("Token %q successfully deleted from the configuration file %q.", key, cfg.ConfigFile)
		}
	}

	return errors.Join(errs...)
}
