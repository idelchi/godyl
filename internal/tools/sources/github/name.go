package github

import (
	"fmt"
	"strings"
)

func SplitName(name string) (parts [2]string, err error) {
	// Split name by first '/'
	split := strings.Split(name, "/")

	// Check if the name is in the correct format
	if len(split) != 2 {
		return parts, fmt.Errorf("invalid source name: %s", name)
	}

	// Set parts to the split values
	parts[0] = split[0] // owner
	parts[1] = split[1] // repo

	return parts, nil
}
