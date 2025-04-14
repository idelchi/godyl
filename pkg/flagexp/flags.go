// Package flagexp provides command-line flag typo detection.
// Extends the pflag package with fuzzy matching to suggest correct
// flag names when users make typos. Note: This is experimental and
// depends on pflag's error message format, which may change.
package flagexp

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/spf13/pflag"

	"github.com/idelchi/godyl/pkg/stringman"
)

// ParseWithSuggestions parses args with typo detection.
// Attempts to parse command-line arguments using pflag. If an
// unknown flag is encountered, uses fuzzy matching to suggest
// the closest matching valid flag name.
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

// extractFlagValue parses flag names from error messages.
// Uses regex to extract flag names from pflag error messages.
// Handles various error formats including bad syntax, unknown
// flags, and unknown shorthand flags.
// TODO(Idelchi): Find a more robust way to extract flag names.
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

// All returns all defined flag names.
// Returns two slices: one containing all long flag names,
// and another containing all shorthand flag characters.
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
