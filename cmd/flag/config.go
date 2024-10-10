package main

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type Update struct {
	Strategy string `mapstructure:"update-strategy"`
	Now      bool   `mapstructure:"update"`
}

// Config holds all the configuration options for godyl.
type Config struct {
	Tools  string
	Update Update `mapstructure:",squash"`
}

// Validate the configuration.
func (c *Config) Validate() error {
	validate := validator.New()
	if err := validate.Struct(c); err != nil {
		return fmt.Errorf("validating config: %w", err)
	}
	return nil
}
