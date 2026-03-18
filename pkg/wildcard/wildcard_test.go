// Package wildcard provides simple, safe string matching where only '*' is
// treated as a wildcard matching any sequence of characters (excluding newlines).
//
// Example:
//
//	ok := wildcard.Match("foo*bar", "fooXYZbar")                // single
//	okAny := wildcard.Match("foo*bar", "nope", "fooZZZbar")     // any-of
//	okEmpty := wildcard.Match("foo*bar")                        // false (no inputs)
package wildcard_test

import (
	"testing"

	"github.com/idelchi/godyl/pkg/wildcard"
)

func TestMatch(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		pattern string
		input   []string
		want    bool
	}{
		{name: "exact match", pattern: "foo", input: []string{"foo"}, want: true},
		{name: "no match", pattern: "foo", input: []string{"bar"}, want: false},
		{name: "wildcard in middle", pattern: "foo*bar", input: []string{"fooXYZbar"}, want: true},
		{name: "wildcard in middle no match", pattern: "foo*bar", input: []string{"foobaz"}, want: false},
		{name: "leading wildcard", pattern: "*bar", input: []string{"foobar"}, want: true},
		{name: "trailing wildcard", pattern: "foo*", input: []string{"foobar"}, want: true},
		{name: "wildcard only", pattern: "*", input: []string{"anything"}, want: true},
		{name: "any of inputs matches", pattern: "foo", input: []string{"bar", "foo"}, want: true},
		{name: "dot is literal", pattern: "foo.bar", input: []string{"foo.bar"}, want: true},
		{name: "dot does not match arbitrary char", pattern: "foo.bar", input: []string{"fooXbar"}, want: false},
		{name: "parens are literal", pattern: "foo(bar)", input: []string{"foo(bar)"}, want: true},
		{name: "empty pattern matches empty string", pattern: "", input: []string{""}, want: true},
		{name: "empty pattern does not match non-empty", pattern: "", input: []string{"x"}, want: false},
		{name: "multiple wildcards", pattern: "a*b*c", input: []string{"aXbYc"}, want: true},
		{name: "wildcard matches empty span", pattern: "foo*bar", input: []string{"foobar"}, want: true},
		{name: "zero variadic args returns false", pattern: "foo", input: nil, want: false},
		{
			name:    "regex metachar caret-dollar treated as literal match",
			pattern: "^foo$",
			input:   []string{"^foo$"},
			want:    true,
		},
		{
			name:    "regex metachar caret-dollar does not match bare foo",
			pattern: "^foo$",
			input:   []string{"foo"},
			want:    false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := wildcard.Match(tc.pattern, tc.input...)
			if got != tc.want {
				t.Errorf("Match(%q, %v) = %v, want %v", tc.pattern, tc.input, got, tc.want)
			}
		})
	}

	t.Run("newline does not match wildcard", func(t *testing.T) {
		t.Parallel()

		if wildcard.Match("a*b", "a\nb") {
			t.Error(`Match("a*b", "a\nb") = true, want false`)
		}
	})

	t.Run("case sensitive: FOO does not match foo", func(t *testing.T) {
		t.Parallel()

		if wildcard.Match("FOO", "foo") {
			t.Error(`Match("FOO", "foo") = true, want false (matching is case-sensitive)`)
		}
	})

	t.Run("double wildcard matches slash-separated path", func(t *testing.T) {
		t.Parallel()

		// ** is treated as two consecutive * wildcards, each matching any
		// run of non-newline characters. Together they behave identically to
		// a single *, so "**" matches "foo/bar" the same way "*" does.
		if !wildcard.Match("**", "foo/bar") {
			t.Error(`Match("**", "foo/bar") = false, want true (** collapses to single wildcard span)`)
		}
	})
}

func TestMatchMultipleWildcards(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		pattern string
		input   string
		want    bool
	}{
		{
			// Three consecutive '*' collapse to a single span — same as "*".
			name:    "triple star matches any string",
			pattern: "***",
			input:   "anything/at/all",
			want:    true,
		},
		{
			// "a**b" is equivalent to "a.*b" in regex terms — the two wildcards
			// together match any sequence of characters between 'a' and 'b'.
			name:    "a**b matches string with content between a and b",
			pattern: "a**b",
			input:   "aXXb",
			want:    true,
		},
		{
			// No content can turn "a**b" into a match for a string that ends before 'b'.
			name:    "a**b does not match string without trailing b",
			pattern: "a**b",
			input:   "aXX",
			want:    false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := wildcard.Match(tc.pattern, tc.input)
			if got != tc.want {
				t.Errorf("Match(%q, %q) = %v, want %v", tc.pattern, tc.input, got, tc.want)
			}
		})
	}
}

func TestMatchEmptyPatternEmptyString(t *testing.T) {
	t.Parallel()

	// An empty pattern anchors to the full string: "^$" — it matches only
	// an empty string.
	if !wildcard.Match("", "") {
		t.Error(`Match("", "") = false, want true`)
	}

	// An empty pattern must not match a non-empty string.
	if wildcard.Match("", "x") {
		t.Error(`Match("", "x") = true, want false`)
	}
}

func TestMatchMultipleInputsShortCircuit(t *testing.T) {
	t.Parallel()

	// The first input "" matches the empty pattern, so Match must return true
	// regardless of subsequent inputs.  This verifies the any-of semantics.
	if !wildcard.Match("", "", "x") {
		t.Error(`Match("", "", "x") = false, want true (first input matches)`)
	}

	// When no input matches, the result is false even with multiple candidates.
	if wildcard.Match("exact", "nope", "also-nope") {
		t.Error(`Match("exact", "nope", "also-nope") = true, want false`)
	}
}
