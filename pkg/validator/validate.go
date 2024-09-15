package validator

import (
	"fmt"
	"strings"
)

// Validate performs validation checks on multiple values.
// Takes a variadic list of values to validate, runs all checks,
// and returns a formatted error message if any validations fail.
// The error includes a bulleted list of failures and help text.
func Validate(validations ...any) error {
	validator := New()

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
	return fmt.Errorf("%w\n%s", ErrValidation, errList.String())
}
