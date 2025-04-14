// Package stringman provides string manipulation and comparison utilities.
// Includes functions for string normalization, fuzzy matching using
// Levenshtein distance, and finding closest matches in a set of strings.
package stringman

import (
	"math"
	"regexp"

	"github.com/agnivade/levenshtein"
	"github.com/schollz/closestmatch"
)

// TransformString normalizes strings to a safe format.
// Replaces all non-alphanumeric characters with dashes,
// making strings suitable for use in URLs, filenames, etc.
func TransformString(input string) string {
	// Regular expression to match allowed characters
	reg := regexp.MustCompile("[^a-zA-Z0-9-]+")
	// Replace disallowed characters with an underscore
	result := reg.ReplaceAllString(input, "-")

	return result
}

// ClosestLevensteinString finds the best match using edit distance.
// Uses Levenshtein distance to find the most similar string from
// a list of candidates. Returns the string with minimum edit
// distance from the input.
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

// ClosestString finds the best match using bag of words.
// Takes a list of candidate strings and substring sizes to
// consider when comparing. Uses bag of words algorithm for
// fuzzy matching, which can be more accurate than pure
// edit distance for longer strings.
func ClosestString(input string, candidates []string, sizes []int) string {
	// Create a closestmatch object with the possible strings and sizes
	cm := closestmatch.New(candidates, sizes)

	// Find and return the closest match
	return cm.Closest(input)
}
