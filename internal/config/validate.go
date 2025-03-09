package config

import (
	"errors"
	"fmt"

	"github.com/idelchi/gogen/pkg/validator"
)

// ErrExitGracefully is an error that signals the program to exit gracefully.
var ErrExitGracefully = errors.New("gracefully exiting")

// Validate performs validation checks.
func Validate(validations ...any) error {
	validator := validator.NewValidator()

	// TODO(Idelchi): Collect all and return them as single error.
	for _, v := range validations {
		errs := validator.Validate(v)

		switch {
		case errs == nil:
			return nil
		case len(errs) == 1:
			return fmt.Errorf("%w: %w\nSee --help for more info on usage", ErrUsage, errs[0])
		case len(errs) > 1:
			return fmt.Errorf("%ws:\n%w\nSee --help for more info on usage", ErrUsage, errors.Join(errs...))
		}
	}

	return nil
}
