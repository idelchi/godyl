// Package validate provides configuration validation functionality.
// Wraps the gogen validator package to provide a user-friendly interface
// for validating configuration options and presenting validation errors
// in a clear, formatted manner with helpful error messages.
package validate

import (
	"errors"
	"fmt"
	"strings"

	"github.com/idelchi/gogen/pkg/validator"
)

// ErrValidation indicates one or more validation checks failed.
// Used as the base error for all validation failures.
var ErrValidation = errors.New("validation error")

// Validate performs validation checks on multiple values.
// Takes a variadic list of values to validate, runs all checks,
// and returns a formatted error message if any validations fail.
// The error includes a bulleted list of failures and help text.
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
