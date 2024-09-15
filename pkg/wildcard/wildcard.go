// Package wildcard provides simple, safe string matching where only '*' is
// treated as a wildcard matching any sequence of characters (excluding newlines).
//
// Example:
//
//	ok := wildcard.Match("foo*bar", "fooXYZbar")                // single
//	okAny := wildcard.Match("foo*bar", "nope", "fooZZZbar")     // any-of
//	okEmpty := wildcard.Match("foo*bar")                        // false (no inputs)
package wildcard

import (
	"regexp"
	"slices"
	"strings"
)

// Match reports whether ANY of the provided strings match pattern.
//   - Only '*' is special; it matches any run of characters (excluding newlines).
//   - All other characters are literal. Matching is anchored to the full string.
//   - With zero inputs, it returns false.
func Match(pattern string, s ...string) bool {
	quoted := regexp.QuoteMeta(pattern)
	glob := strings.ReplaceAll(quoted, `\*`, `.*`)

	re := regexp.MustCompile("^" + glob + "$")

	return slices.ContainsFunc(s, re.MatchString)
}
