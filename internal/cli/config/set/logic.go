package set

import (
	"fmt"

	"github.com/idelchi/godyl/internal/cli/common"
	"github.com/idelchi/godyl/internal/config/root"
	"github.com/idelchi/godyl/pkg/editor"
	"github.com/idelchi/godyl/pkg/koanfx"
	"github.com/idelchi/godyl/pkg/unmarshal"
)

// run executes the `config set` command.
func run(input common.Input) error {
	cfg, _, context, _, args := input.Unpack()

	koanf := context.Config

	if err := validate(koanf, args[0], args[1]); err != nil {
		return err
	}

	return editor.New(cfg.ConfigFile.Absolute()).Merge(koanf.Map())
}

// validate checks if the provided key and value are valid for the configuration.
func validate(koanf *koanfx.Koanf, path, value string) error {
	// First unmarshal the value string as YAML to get proper typing
	var typedValue any
	if err := unmarshal.Lax([]byte(value), &typedValue); err != nil {
		return fmt.Errorf("parsing value as YAML: %w", err)
	}

	if err := koanf.Set(path, typedValue); err != nil {
		return err
	}

	var cfg root.Config

	if err := koanf.Unmarshal(&cfg, koanfx.WithErrorUnused()); err != nil {
		return err
	}

	return nil
}
