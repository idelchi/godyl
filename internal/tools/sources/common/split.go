package common

import (
	"fmt"
	"strings"
)

// SplitName splits a full repository name in the format "X/Y" into its respective components.
// It returns an error if the input name is not in the correct format.
func SplitName(name string) (first, second string, err error) {
	// Split the name by the first '/'
	split := strings.Split(name, "/")

	// Check if the name is in the correct format
	const expectedParts = 2
	if len(split) != expectedParts {
		return first, second, fmt.Errorf("invalid source name: %s", name)
	}

	return split[0], split[1], nil
}
