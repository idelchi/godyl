package config

import (
	"fmt"

	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/structs"

	"github.com/idelchi/godyl/internal/config"
	"github.com/idelchi/godyl/internal/iutils"
	"github.com/idelchi/godyl/pkg/koanfx"
)

// run executes the `config dump` command.
func run(cfg config.Config, koanf *koanfx.KoanfWithTracker, paths ...string) error {
	file := cfg.ConfigFile
	if !file.Exists() {
		fmt.Printf("Config file %q doesn't exist\n", file)

		return nil
	}

	if len(paths) > 0 {
		if cfg.Dump.Config.All {
			emptyConfig := config.Config{}
			k := koanfx.New()

			if err := k.Load(structs.Provider(emptyConfig, "yaml"), nil); err != nil {
				return err
			}

			k.Load(confmap.Provider(koanf.All(), "."), nil)

			koanf = k
		}

		for _, path := range paths {
			val := koanf.Get(path)
			if val == nil {
				return fmt.Errorf("value %q not found in config", path)
			}

			if len(paths) > 1 {
				fmt.Printf(" ---- %s ----\n", path)
			}

			iutils.Print(iutils.YAML, val)

			if len(paths) > 1 {
				fmt.Println()
			}
		}

		return nil
	}

	if cfg.Dump.Config.All {
		var raw config.Config

		koanf.Unmarshal(&raw)

		iutils.Print(iutils.YAML, raw)

		return nil
	}

	iutils.Print(iutils.YAML, koanf.Raw())

	return nil
}
