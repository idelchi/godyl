// Package version provides some utilities for working with version strings.
package version

import (
	"strings"
	"unicode"

	"github.com/Masterminds/semver/v3"
)

// Parse attempts to extract the semantic version from a string.
// It iterates through the string, left to right, looking for a valid semantic version.
// If no valid version is found, it returns nil.
func Parse(name string) *semver.Version {
	for index := range len(name) {
		candidate := name[index:]
		if startsWithNonDigit(candidate) && !strings.HasPrefix(candidate, "v") {
			continue
		}

		if version, err := semver.NewVersion(candidate); err == nil {
			return version
		}
	}

	return nil
}

// Compare compares two version strings for equality.
// A failure will always return false.
func Compare(a, b string) bool {
	// Convert the version strings to semantic versions.
	aVersion := Parse(a)
	bVersion := Parse(b)

	// If either version is nil, return false.
	if aVersion == nil || bVersion == nil {
		return false
	}

	// Compare the two versions.
	return aVersion.Equal(bVersion)
}

// LessThan compares two version strings and returns true if the first version is less than the second.
func LessThan(a, b string) bool {
	// Convert the version strings to semantic versions.
	aVersion := Parse(a)
	bVersion := Parse(b)

	// If either version is nil, return true.
	if aVersion == nil || bVersion == nil {
		return true
	}

	// Compare the two versions.
	return aVersion.LessThan(bVersion)
}

// startsWithNonDigit checks if the string starts with a non-digit character.
func startsWithNonDigit(s string) bool {
	return len(s) > 0 && !unicode.IsDigit(rune(s[0]))
}
