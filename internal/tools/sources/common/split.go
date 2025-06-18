package common

import (
	"fmt"
	"strings"
)

// SplitName splits a repository name into owner and repository components.
// Expects input in the format "owner/repo" and returns an error if not properly formatted.
func SplitName(name string) (first, second string, err error) {
	// Split the name by the first '/'
	split := strings.Split(name, "/")

	// Check if the name is in the correct format
	const expectedParts = 2
	if len(split) != expectedParts {
		return first, second, fmt.Errorf("invalid name: %s", name)
	}

	return split[0], split[1], nil
}

// CutName splits a repository path into its first component and remaining path.
// Uses strings.Cut to split on the first "/" character. Returns an error if
// the input does not contain a "/" character.
func CutName(name string) (first, second string, err error) {
	// Split the name by the first '/'
	first, second, found := strings.Cut(name, "/")

	if !found {
		return first, second, fmt.Errorf("invalid name: %s", name)
	}

	return first, second, nil
}
