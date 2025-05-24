package set

import (
	"fmt"

	"github.com/idelchi/godyl/internal/config"
	"github.com/idelchi/godyl/pkg/koanfx"
	"github.com/idelchi/godyl/pkg/path/file"
	"github.com/idelchi/godyl/pkg/pretty"
	"github.com/idelchi/godyl/pkg/unmarshal"
)

func run(file file.File, koanf *koanfx.KoanfWithTracker, path, value string) error {
	if !file.Exists() {
		return fmt.Errorf("config file %q doesn't exist", file)
	}

	// First unmarshal the value string as YAML to get proper typing
	var typedValue any
	if err := unmarshal.Lax([]byte(value), &typedValue); err != nil {
		return fmt.Errorf("parsing value as YAML: %w", err)
	}

	if err := koanf.Set(path, typedValue); err != nil {
		return fmt.Errorf("setting value in config file %q: %w", file, err)
	}

	var cfg config.Config

	if err := koanf.Unmarshal(&cfg, koanfx.WithErrorUnused()); err != nil {
		return err
	}

	mapAny, err := koanf.AsMapAny()
	if err != nil {
		return fmt.Errorf("getting config as map: %w", err)
	}

	// Write the updated config back to the file
	if err := file.Write([]byte(pretty.YAML(mapAny))); err != nil {
		return fmt.Errorf("writing config file %q: %w", file, err)
	} else {
		fmt.Printf("Updated config key %q\n", path)
	}

	return nil
}
