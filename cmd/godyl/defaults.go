package main

import (
	"fmt"
	"os"

	_ "embed"

	"github.com/go-playground/validator/v10"
	"github.com/idelchi/godyl/internal/tools"
	"gopkg.in/yaml.v3"
)

//go:embed defaults.yml
var defaultsFile []byte

// Defaults holds all the Defaultsuration options for godyl.
type Defaults struct {
	// Defaults for tools. Allows setting a default subset of values for tools
	tools.Defaults `yaml:",inline"`
}

func (d *Defaults) Unmarshal(data []byte) error {
	return yaml.Unmarshal(data, d)
}

func (d *Defaults) FromFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return d.Unmarshal(data)
}

func (d *Defaults) Default() error {
	return d.Unmarshal(defaultsFile)
}

// Validate the Defaultsuration.
func (d *Defaults) Validate() error {
	validate := validator.New()
	if err := validate.Struct(d); err != nil {
		return fmt.Errorf("validating Defaults: %w", err)
	}
	return nil
}

func (d *Defaults) Merge(cfg Config) {
	if IsSet("output") {
		d.Output = cfg.Output
	}

	if IsSet("source") {
		d.Source.Type = cfg.Source
	}

	if IsSet("strategy") {
		d.Strategy = cfg.Update.Strategy
	}

	if IsSet("github-token") {
		d.Source.Github.Token = cfg.Tokens.GitHub
	}
}

func (d *Defaults) Load(path string) error {
	if IsSet("defaults") {
		if err := d.FromFile(path); err != nil {
			return fmt.Errorf("loading defaults from %q: %w", path, err)
		}
	} else {
		if err := d.Default(); err != nil {
			return fmt.Errorf("setting defaults: %w", err)
		}
	}

	if err := d.Initialize(); err != nil {
		return fmt.Errorf("setting tool defaults: %w", err)
	}

	return nil
}
