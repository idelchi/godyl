package cobraext

import (
	"errors"
	"fmt"

	"github.com/spf13/viper"

	"github.com/idelchi/godyl/pkg/pretty"
)

// Validator is an interface for types that can show and validate their configuration.
type Validator interface {
	Validate(config any) error
	Display() bool
}

// ErrExitGracefully is an error that signals the program to exit gracefully.
var ErrExitGracefully = errors.New("exit")

// Validate unmarshals the configuration and performs validation checks.
// If cfg.Show is true, prints the configuration and exits.
func Validate(cfg Validator, validations ...any) error {
	if err := viper.Unmarshal(cfg); err != nil {
		return fmt.Errorf("unmarshalling config: %w", err)
	}

	if cfg.Display() {
		pretty.PrintJSONMasked(cfg)

		return ErrExitGracefully
	}

	for _, v := range validations {
		if err := cfg.Validate(v); err != nil {
			return fmt.Errorf("validating config: %w\nSee --help for more info on usage", err)
		}
	}

	return nil
}
