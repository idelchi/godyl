package remove

import (
	"fmt"

	"github.com/idelchi/godyl/internal/cli/core"
	"github.com/idelchi/godyl/pkg/editor"
)

// run executes the `config remove` command.
func run(input core.Input) error {
	cfg, _, context, _, args := input.Unpack()

	configuration := context.Config

	logger, err := core.SetupLogger(cfg.LogLevel)
	if err != nil {
		return err
	}

	switch len(args) {
	case 0:
		configuration.Delete("")

		fmt.Println("Removing all configuration keys.")
	default:
		for _, key := range args {
			if !configuration.Exists(key) {
				logger.Warnf("Key %q does not exist in the configuration file.", key)

				continue
			}

			configuration.Delete(key)
			logger.Infof("Removed key %q from the configuration file.", key)
		}
	}

	if err := editor.New(cfg.ConfigFile).Write(configuration.Map()); err != nil {
		return err
	}

	return nil
}
