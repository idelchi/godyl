package config

import (
	"fmt"
	"strings"

	"github.com/idelchi/gogen/pkg/validator"
)

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

	if len(allErrors) == 1 {
		return fmt.Errorf("%w: %w\nSee --help for more info on usage", ErrUsage, allErrors[0])
	}

	// Create a bulleted list for multiple errors
	var errList strings.Builder
	errList.WriteString("validation errors:\n")
	for _, err := range allErrors {
		errList.WriteString(fmt.Sprintf("  • %v\n", err))
	}
	errList.WriteString("See --help for more info on usage")

	return fmt.Errorf("%w: %s", ErrUsage, errList.String())
}
