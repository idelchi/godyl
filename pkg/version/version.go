// Package version provides some utilities for working with version strings.
package version

import (
	"unicode"

	"github.com/Masterminds/semver/v3"
)

// To attempts to convert a version string to a semantic version.
func To(version string) *semver.Version {
	for index := range len(version) {
		candidate := version[index:]
		if startsWithNonDigit(candidate) {
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
	aVersion := To(a)
	bVersion := To(b)

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
	aVersion := To(a)
	bVersion := To(b)

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
