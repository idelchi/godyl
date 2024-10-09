// Package stringman provides utilities for working with strings.
package stringman

import (
	"math"
	"regexp"

	"github.com/agnivade/levenshtein"
	"github.com/schollz/closestmatch"
)

// TransformString replaces all non-alphanumeric characters with a dash.
func TransformString(input string) string {
	// Regular expression to match allowed characters
	reg := regexp.MustCompile("[^a-zA-Z0-9-]+")
	// Replace disallowed characters with an underscore
	result := reg.ReplaceAllString(input, "-")

	return result
}

// ClosestLevensteinString finds the closest string from a list of possible strings using Levenshtein distance.
func ClosestLevensteinString(wrongString string, possibleStrings []string) string {
	minDistance := math.MaxInt32
	closest := ""

	for _, str := range possibleStrings {
		distance := levenshtein.ComputeDistance(wrongString, str)
		if distance < minDistance {
			minDistance = distance
			closest = str
		}
	}

	return closest
}

// ClosestString returns the closest string from a list of possible strings.
func ClosestString(input string, candidates []string, sizes []int) string {
	// Create a closestmatch object with the possible strings and sizes
	cm := closestmatch.New(candidates, sizes)

	// Find and return the closest match
	return cm.Closest(input)
}
