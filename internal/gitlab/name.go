package gitlab

import (
	"fmt"
	"strings"
)

// SplitName splits a full repository name in the format "owner/repo" into its owner and repo components.
// It returns an error if the input name is not in the correct format.
func SplitName(name string) (parts [2]string, err error) {
	// Split the name by the first '/'
	split := strings.Split(name, "/")

	// Check if the name is in the correct format
	const expectedParts = 2
	if len(split) != expectedParts {
		return parts, fmt.Errorf("invalid source name: %s", name)
	}

	// Set parts to the split values: owner and repo
	parts[0] = split[0] // owner
	parts[1] = split[1] // repo

	return parts, nil
}