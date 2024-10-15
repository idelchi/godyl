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

func (d *Defaults) Load(data []byte) error {
	return yaml.Unmarshal(data, d)
}

func (d *Defaults) FromFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return d.Load(data)
}

func (d *Defaults) Default() error {
	return d.Load(defaultsFile)
}

// Validate the Defaultsuration.
func (d *Defaults) Validate() error {
	validate := validator.New()
	if err := validate.Struct(d); err != nil {
		return fmt.Errorf("validating Defaults: %w", err)
	}
	return nil
}

func (d Defaults) IsSet() bool {
	return IsSet("defaults")
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
