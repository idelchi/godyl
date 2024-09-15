package config

import (
	"fmt"

	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/structs"

	"github.com/idelchi/godyl/internal/cli/common"
	"github.com/idelchi/godyl/internal/config/root"
	"github.com/idelchi/godyl/internal/iutils"
	"github.com/idelchi/godyl/pkg/koanfx"
)

// run executes the `dump config` command.
func run(input common.Input) error {
	cfg, _, context, _, args := input.Unpack()

	file := cfg.ConfigFile
	if !file.Exists() {
		fmt.Printf("Config file %q doesn't exist\n", file)

		return nil
	}

	configuration := context.Config

	if len(args) > 0 {
		if cfg.Dump.Config.Full {
			emptyConfig := root.Config{}
			k := koanfx.New()

			if err := k.Load(structs.Provider(emptyConfig, "json"), nil); err != nil {
				return err
			}

			k.Load(confmap.Provider(configuration.All(), "."), nil)

			configuration = k
		}

		for _, path := range args {
			val := configuration.Get(path)
			if val == nil {
				return fmt.Errorf("value %q not found in config", path)
			}

			if len(args) > 1 {
				fmt.Printf(" ---- %s ----\n", path)
			}

			iutils.Print(iutils.YAML, val)

			if len(args) > 1 {
				fmt.Println()
			}
		}

		return nil
	}

	if cfg.Dump.Config.Full {
		var raw root.Config

		if err := configuration.Unmarshal(&raw); err != nil {
			return err
		}

		iutils.Print(iutils.YAML, raw)

		return nil
	}

	iutils.Print(iutils.YAML, configuration.Map())

	return nil
}
