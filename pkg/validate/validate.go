// Package validate provides functionality for validating configuration options.
// It wraps the `gogen` validator package to provide a more user-friendly interface.
package validate

import (
	"errors"
	"fmt"
	"strings"

	"github.com/idelchi/gogen/pkg/validator"
)

var ErrValidation = errors.New("validation error")

// Validate performs validation checks.
func Validate(validations ...any) error {
	validator := validator.NewValidator()

	var allErrors []error

	for _, v := range validations {
		errs := validator.Validate(v)
		if errs != nil {
			allErrors = append(allErrors, errs...)
		}
	}

	if len(allErrors) == 0 {
		return nil
	}

	// Create a bulleted list for the errors
	var errList strings.Builder
	for _, err := range allErrors {
		errList.WriteString(fmt.Sprintf("  • %s\n", err))
	}

	// Add help text
	return fmt.Errorf("%w\n%sSee --help for more info on usage", ErrValidation, errList.String())
}
