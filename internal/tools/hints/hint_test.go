package hints_test

import (
	"testing"

	"github.com/idelchi/godyl/internal/tools/hints"
)

func TestHintMatches(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		hint      hints.Hint
		input     string
		wantMatch bool
		wantErr   bool
	}{
		{
			name:      "contains match",
			hint:      hints.Hint{Pattern: "linux", Type: hints.Contains},
			input:     "tool-linux-amd64",
			wantMatch: true,
		},
		{
			name:      "contains no match",
			hint:      hints.Hint{Pattern: "linux", Type: hints.Contains},
			input:     "tool-darwin",
			wantMatch: false,
		},
		{
			name:      "startswith match",
			hint:      hints.Hint{Pattern: "tool", Type: hints.StartsWith},
			input:     "tool-v1.0",
			wantMatch: true,
		},
		{
			name:      "startswith no match",
			hint:      hints.Hint{Pattern: "tool", Type: hints.StartsWith},
			input:     "my-tool",
			wantMatch: false,
		},
		{
			name:      "endswith match",
			hint:      hints.Hint{Pattern: ".tar.gz", Type: hints.EndsWith},
			input:     "file.tar.gz",
			wantMatch: true,
		},
		{
			name:      "endswith no match",
			hint:      hints.Hint{Pattern: ".tar.gz", Type: hints.EndsWith},
			input:     "file.zip",
			wantMatch: false,
		},
		{
			name:      "glob match",
			hint:      hints.Hint{Pattern: "*.tar.gz", Type: hints.Glob},
			input:     "file.tar.gz",
			wantMatch: true,
		},
		{
			name:      "glob no match",
			hint:      hints.Hint{Pattern: "*.tar.gz", Type: hints.Glob},
			input:     "file.zip",
			wantMatch: false,
		},
		{
			name:    "glob bad pattern returns error",
			hint:    hints.Hint{Pattern: "[", Type: hints.Glob},
			input:   "anything",
			wantErr: true,
		},
		{
			name:      "globstar match across path separator",
			hint:      hints.Hint{Pattern: "**/file.tar.gz", Type: hints.GlobStar},
			input:     "some/nested/dir/file.tar.gz",
			wantMatch: true,
		},
		{
			name:      "globstar no match",
			hint:      hints.Hint{Pattern: "**/file.tar.gz", Type: hints.GlobStar},
			input:     "file.zip",
			wantMatch: false,
		},
		{
			name:    "globstar bad pattern returns error",
			hint:    hints.Hint{Pattern: "[", Type: hints.GlobStar},
			input:   "anything",
			wantErr: true,
		},
		{
			name:      "regex match",
			hint:      hints.Hint{Pattern: `^[^.]+$`, Type: hints.Regex},
			input:     "noext",
			wantMatch: true,
		},
		{
			name:      "regex no match",
			hint:      hints.Hint{Pattern: `^[^.]+$`, Type: hints.Regex},
			input:     "has.ext",
			wantMatch: false,
		},
		{
			name:    "regex invalid pattern returns error",
			hint:    hints.Hint{Pattern: "[invalid", Type: hints.Regex},
			input:   "anything",
			wantErr: true,
		},
		{
			name:    "unknown type returns error",
			hint:    hints.Hint{Pattern: "anything", Type: hints.Type("unknown")},
			input:   "anything",
			wantErr: true,
		},
		{
			name:    "empty type returns error",
			hint:    hints.Hint{Pattern: "anything", Type: hints.Type("")},
			input:   "anything",
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, err := tc.hint.Matches(tc.input)

			if tc.wantErr {
				if err == nil {
					t.Fatalf("Matches(%q) expected error, got nil", tc.input)
				}

				return
			}

			if err != nil {
				t.Fatalf("Matches(%q) unexpected error: %v", tc.input, err)
			}

			if got != tc.wantMatch {
				t.Errorf("Matches(%q) = %v, want %v (pattern=%q, type=%q)",
					tc.input, got, tc.wantMatch, tc.hint.Pattern, tc.hint.Type)
			}
		})
	}
}

func TestHintParse(t *testing.T) {
	t.Parallel()

	t.Run("zero-value hint gets defaults", func(t *testing.T) {
		t.Parallel()

		h := hints.Hint{}
		if err := h.Parse(); err != nil {
			t.Fatalf("Parse() unexpected error on zero-value Hint: %v", err)
		}

		if h.Type != hints.Glob {
			t.Errorf("Parse() Type = %q, want %q", h.Type, hints.Glob)
		}

		weight, err := h.Weight.Get()
		if err != nil {
			t.Fatalf("Weight.Get() unexpected error: %v", err)
		}

		if weight != 1 {
			t.Errorf("Parse() Weight = %d, want 1", weight)
		}

		match, err := h.Match.Get()
		if err != nil {
			t.Fatalf("Match.Get() unexpected error: %v", err)
		}

		if match != hints.Weighted {
			t.Errorf("Parse() Match = %q, want %q", match, hints.Weighted)
		}
	})

	t.Run("invalid match value returns error", func(t *testing.T) {
		t.Parallel()

		h := hints.Hint{Pattern: "linux", Type: hints.Contains}
		h.Match.Set("bogus")

		if err := h.Parse(); err == nil {
			t.Fatal("Parse() expected error for invalid Match value, got nil")
		}
	})

	t.Run("invalid weight template returns error", func(t *testing.T) {
		t.Parallel()

		h := hints.Hint{Pattern: "linux", Type: hints.Contains}
		// yaml.Unmarshal of "not-a-number" into int will fail.
		h.Weight.Set("not-a-number")

		if err := h.Parse(); err == nil {
			t.Fatal("Parse() expected error for non-integer Weight template, got nil")
		}
	})

	t.Run("valid weight template stores parsed value", func(t *testing.T) {
		t.Parallel()

		h := hints.Hint{Pattern: "linux", Type: hints.Contains}
		h.Weight.Set("5")

		if err := h.Parse(); err != nil {
			t.Fatalf("Parse() unexpected error for valid Weight template: %v", err)
		}

		weight, err := h.Weight.Get()
		if err != nil {
			t.Fatalf("Weight.Get() unexpected error: %v", err)
		}

		if weight != 5 {
			t.Errorf("Weight.Get() = %d, want 5", weight)
		}
	})
}

// TestHintParseAndMatch exercises the full Parse→Matches path in one test,
// verifying that a Hint round-trips correctly from construction through matching.
func TestHintParseAndMatch(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		pattern   string
		typ       hints.Type
		matchSet  hints.Match // value used to call Match.Set
		input     string
		wantMatch bool
		wantErr   bool
	}{
		{
			name:      "contains type matches after parse",
			pattern:   "linux",
			typ:       hints.Contains,
			matchSet:  hints.Required,
			input:     "tool-linux-amd64",
			wantMatch: true,
		},
		{
			name:      "contains type does not match after parse",
			pattern:   "linux",
			typ:       hints.Contains,
			matchSet:  hints.Required,
			input:     "tool-darwin",
			wantMatch: false,
		},
		{
			name:      "glob type matches after parse",
			pattern:   "*.tar.gz",
			typ:       hints.Glob,
			matchSet:  hints.Weighted,
			input:     "tool.tar.gz",
			wantMatch: true,
		},
		{
			name:      "regex type matches after parse",
			pattern:   `^tool-`,
			typ:       hints.Regex,
			matchSet:  hints.Excluded,
			input:     "tool-linux-amd64",
			wantMatch: true,
		},
		{
			name:      "startswith type matches after parse",
			pattern:   "tool",
			typ:       hints.StartsWith,
			matchSet:  hints.Weighted,
			input:     "tool-v1.0",
			wantMatch: true,
		},
		{
			name:      "endswith type matches after parse",
			pattern:   ".zip",
			typ:       hints.EndsWith,
			matchSet:  hints.Weighted,
			input:     "archive.zip",
			wantMatch: true,
		},
		{
			name:      "globstar type matches nested path after parse",
			pattern:   "**/file.tar.gz",
			typ:       hints.GlobStar,
			matchSet:  hints.Weighted,
			input:     "some/nested/file.tar.gz",
			wantMatch: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			h := hints.Hint{Pattern: tc.pattern, Type: tc.typ}
			h.Match.Set(string(tc.matchSet))

			if err := h.Parse(); err != nil {
				t.Fatalf("Parse() unexpected error: %v", err)
			}

			got, err := h.Matches(tc.input)

			if tc.wantErr {
				if err == nil {
					t.Fatalf("Matches(%q) expected error, got nil", tc.input)
				}

				return
			}

			if err != nil {
				t.Fatalf("Matches(%q) unexpected error: %v", tc.input, err)
			}

			if got != tc.wantMatch {
				t.Errorf("Matches(%q) = %v, want %v (pattern=%q, type=%q)",
					tc.input, got, tc.wantMatch, tc.pattern, tc.typ)
			}
		})
	}
}

// TestHintsReducedIdentity verifies which specific hints remain after Reduced(),
// not just the count. An empty Pattern causes the hint to be removed; non-empty
// patterns are preserved in their original order.
func TestHintsReducedIdentity(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		hints        hints.Hints
		wantPatterns []string // patterns of the hints that should survive Reduced()
	}{
		{
			name: "empty pattern hint removed; others preserved in order",
			hints: hints.Hints{
				{Pattern: "linux", Type: hints.Contains},
				{Pattern: "", Type: hints.Glob},
				{Pattern: "amd64", Type: hints.Contains},
			},
			wantPatterns: []string{"linux", "amd64"},
		},
		{
			name: "all non-empty patterns survive",
			hints: hints.Hints{
				{Pattern: "linux", Type: hints.Contains},
				{Pattern: "*.tar.gz", Type: hints.Glob},
			},
			wantPatterns: []string{"linux", "*.tar.gz"},
		},
		{
			name: "all empty patterns results in empty slice",
			hints: hints.Hints{
				{Pattern: "", Type: hints.Glob},
				{Pattern: "", Type: hints.Contains},
			},
			wantPatterns: []string{},
		},
		{
			name:         "nil-equivalent empty hints produces empty",
			hints:        hints.Hints{},
			wantPatterns: []string{},
		},
		{
			name: "single non-empty pattern survives",
			hints: hints.Hints{
				{Pattern: "windows", Type: hints.Contains},
			},
			wantPatterns: []string{"windows"},
		},
		{
			name: "leading and trailing empty patterns removed",
			hints: hints.Hints{
				{Pattern: "", Type: hints.Glob},
				{Pattern: "darwin", Type: hints.Contains},
				{Pattern: "", Type: hints.Glob},
			},
			wantPatterns: []string{"darwin"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			reduced := tc.hints.Reduced()
			if reduced == nil {
				t.Fatal("Reduced() returned nil, want non-nil pointer")
			}

			if len(*reduced) != len(tc.wantPatterns) {
				t.Fatalf("Reduced() len = %d, want %d", len(*reduced), len(tc.wantPatterns))
			}

			for i, h := range *reduced {
				if h.Pattern != tc.wantPatterns[i] {
					t.Errorf("Reduced()[%d].Pattern = %q, want %q", i, h.Pattern, tc.wantPatterns[i])
				}
			}
		})
	}
}

// TestHintHas verifies the Has() method on Hints.
func TestHintHas(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		hints hints.Hints
		want  bool
	}{
		{
			name:  "nil hints returns false",
			hints: nil,
			want:  false,
		},
		{
			name:  "empty hints returns false",
			hints: hints.Hints{},
			want:  false,
		},
		{
			name:  "single hint returns true",
			hints: hints.Hints{{Pattern: "linux", Type: hints.Contains}},
			want:  true,
		},
		{
			name: "multiple hints returns true",
			hints: hints.Hints{
				{Pattern: "linux", Type: hints.Contains},
				{Pattern: "amd64", Type: hints.Contains},
			},
			want: true,
		},
		{
			// Has() checks length only, not whether the hint patterns are non-empty.
			name:  "hint with empty pattern still counts as present",
			hints: hints.Hints{{Pattern: "", Type: hints.Glob}},
			want:  true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := tc.hints.Has()
			if got != tc.want {
				t.Errorf("Has() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestHintsParseAndReduced(t *testing.T) {
	t.Parallel()

	t.Run("Parse sets defaults on each hint", func(t *testing.T) {
		t.Parallel()

		hs := hints.Hints{
			{Pattern: "linux", Type: hints.Contains},
			{Pattern: "amd64", Type: hints.Contains},
		}

		if err := hs.Parse(); err != nil {
			t.Fatalf("Parse() unexpected error: %v", err)
		}

		for i := range hs {
			w, err := hs[i].Weight.Get()
			if err != nil {
				t.Fatalf("hint[%d] Weight.Get() error: %v", i, err)
			}

			if w != 1 {
				t.Errorf("hint[%d] Weight = %d, want 1 after Parse()", i, w)
			}
		}
	})

	t.Run("Reduced removes hints with empty pattern", func(t *testing.T) {
		t.Parallel()

		hs := hints.Hints{
			{Pattern: "linux", Type: hints.Contains},
			{Pattern: "", Type: hints.Glob},
			{Pattern: "amd64", Type: hints.Contains},
		}

		reduced := hs.Reduced()
		if reduced == nil {
			t.Fatal("Reduced() returned nil, want non-nil")
		}

		if len(*reduced) != 2 {
			t.Errorf("Reduced() len = %d, want 2", len(*reduced))
		}
	})

	t.Run("Has returns false for empty collection", func(t *testing.T) {
		t.Parallel()

		hs := hints.Hints{}
		if hs.Has() {
			t.Error("Has() = true for empty Hints, want false")
		}
	})

	t.Run("Has returns true for non-empty collection", func(t *testing.T) {
		t.Parallel()

		hs := hints.Hints{{Pattern: "linux", Type: hints.Contains}}
		if !hs.Has() {
			t.Error("Has() = false for non-empty Hints, want true")
		}
	})
}
