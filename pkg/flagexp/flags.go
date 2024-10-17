// Package flagexp provides some experimental features to provide suggestions for typos in flag names.
// Requires the `pflag` package and is highly sensitive to changes in the package,
// as it's based on the error messages returned by the package.
package flagexp

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/spf13/pflag"

	"github.com/idelchi/godyl/pkg/stringman"
)

// ParseWithSuggestions parses the command-line arguments and suggests the closest flag name if a typo is detected.
func ParseWithSuggestions(args []string) (err error) {
	// Continue on error so that the error can be parsed
	pflag.CommandLine.Init("", pflag.ContinueOnError)

	err = pflag.CommandLine.Parse(args)

	// If the error is nil or a help message, skip remaining steps
	if err == nil || errors.Is(err, pflag.ErrHelp) {
		return nil
	}

	// Extract the flag name from the error message
	flagName, errExtract := extractFlagValue(err.Error())

	if errExtract != nil {
		return fmt.Errorf("%w: extracting flag value: %w", err, errExtract)
	}

	// Word sizes for the closest match algorithm
	// TODO(Idelchi): Should be flexible.
	sizes := []int{2, 3, 4, 5, 6, 7, 8, 9, 10}

	var closestMatch string

	// For long flags, use `ClosestString`, for short flags use `ClosestLevensteinString
	long, short := All()

	if len(flagName) > 1 {
		closestMatch = stringman.ClosestString(flagName, long, sizes)
	} else {
		closestMatch = stringman.ClosestLevensteinString(flagName, short)
	}

	if closestMatch != "" {
		return fmt.Errorf("%w: did you mean %q?", err, closestMatch)
	}

	return fmt.Errorf("%w: no suggestions available, see `--help`", err)
}

// extractFlagValue extracts the flag value from a given error message.
//
// TODO(Idelchi): Experimental function, there should be a nicer way of extracting the flag name from the error message.
func extractFlagValue(errMsg string) (string, error) {
	expression := `bad flag syntax: (.+)|unknown flag: --(.+)|unknown flag: -(.+)|unknown shorthand flag: '([^'])' in -.*`
	// Define regex pattern to match the error messages and extract the dynamic part
	re := regexp.MustCompile(expression)

	// FindSubmatch returns the slices of submatches, the first element is the full match
	matches := re.FindStringSubmatch(errMsg)
	if len(matches) > 2 { //nolint:mnd // Clear from the context above.
		for _, match := range matches[1:] {
			if match != "" {
				return match, nil
			}
		}
	}

	//nolint:goerr113 	// Once in a while is OK.
	return "", fmt.Errorf("no match found for %q in %q", errMsg, expression)
}

// All returns all flags defined.
func All() ([]string, []string) {
	var (
		longFlags  []string
		shortFlags []string
	)

	pflag.VisitAll(func(flag *pflag.Flag) {
		longFlags = append(longFlags, flag.Name)

		if flag.Shorthand != "" {
			shortFlags = append(shortFlags, flag.Shorthand)
		}
	})

	return longFlags, shortFlags
}
