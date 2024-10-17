package utils

import "strings"

// Equal compares two strings case-insensitively and returns true if they are equal.
func EqualLower(a, b string) bool {
	return strings.ToLower(a) == strings.ToLower(b)
}

// ContainsLower checks if string 'a' contains string 'b' case-insensitively.
func ContainsLower(a, b string) bool {
	return strings.Contains(strings.ToLower(a), strings.ToLower(b))
}
