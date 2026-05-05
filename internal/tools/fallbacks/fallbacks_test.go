package fallbacks_test

import (
	"slices"
	"testing"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/goccy/go-yaml"

	"github.com/idelchi/godyl/internal/tools/fallbacks"
	"github.com/idelchi/godyl/internal/tools/sources"
)

func TestCompact(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input []string
		want  []string
	}{
		{
			name:  "removes duplicates preserving first occurrence order",
			input: []string{"a", "b", "a", "c"},
			want:  []string{"a", "b", "c"},
		},
		{
			name:  "all duplicates collapses to single element",
			input: []string{"a", "a", "a"},
			want:  []string{"a"},
		},
		{
			name:  "no duplicates returns same elements in order",
			input: []string{"x", "y", "z"},
			want:  []string{"x", "y", "z"},
		},
		{
			name:  "single element unchanged",
			input: []string{"only"},
			want:  []string{"only"},
		},
		{
			name:  "nil input returns empty slice",
			input: nil,
			want:  []string{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := fallbacks.Compact(tc.input)

			if !slices.Equal(got, tc.want) {
				t.Errorf("Compact(%v) = %v, want %v", tc.input, got, tc.want)
			}
		})
	}
}

func TestCompacted(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input fallbacks.Fallbacks
		want  fallbacks.Fallbacks
	}{
		{
			name:  "removes duplicates preserving first occurrence order",
			input: fallbacks.Fallbacks{sources.GITHUB, sources.GITLAB, sources.GITHUB},
			want:  fallbacks.Fallbacks{sources.GITHUB, sources.GITLAB},
		},
		{
			name:  "all duplicates collapses to single element",
			input: fallbacks.Fallbacks{sources.GITHUB, sources.GITHUB, sources.GITHUB},
			want:  fallbacks.Fallbacks{sources.GITHUB},
		},
		{
			name:  "no duplicates returns same elements in order",
			input: fallbacks.Fallbacks{sources.GITHUB, sources.GITLAB, sources.URL},
			want:  fallbacks.Fallbacks{sources.GITHUB, sources.GITLAB, sources.URL},
		},
		{
			name:  "nil input returns empty slice",
			input: nil,
			want:  fallbacks.Fallbacks{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := tc.input.Compacted()

			if !slices.Equal(got, tc.want) {
				t.Errorf("Compacted(%v) = %v, want %v", tc.input, got, tc.want)
			}
		})
	}
}

// TestFallbacksUnmarshalYAML verifies that Fallbacks can be deserialized from
// both scalar string and sequence YAML forms, as required by UnmarshalYAML.
func TestFallbacksUnmarshalYAML(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  fallbacks.Fallbacks
	}{
		{
			// Single scalar string is wrapped into a one-element slice.
			name:  "scalar string form yields single-element Fallbacks",
			input: `github`,
			want:  fallbacks.Fallbacks{sources.GITHUB},
		},
		{
			// Sequence form is decoded element by element.
			name: "sequence form yields multi-element Fallbacks",
			input: heredoc.Doc(`
				- github
				- gitlab
			`),
			want: fallbacks.Fallbacks{sources.GITHUB, sources.GITLAB},
		},
		{
			// A sequence with a single element is also valid.
			name:  "sequence with single element",
			input: `[url]`,
			want:  fallbacks.Fallbacks{sources.URL},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var got fallbacks.Fallbacks

			if err := yaml.Unmarshal([]byte(tc.input), &got); err != nil {
				t.Fatalf("yaml.Unmarshal(%q) unexpected error: %v", tc.input, err)
			}

			if !slices.Equal([]sources.Type(got), []sources.Type(tc.want)) {
				t.Errorf("UnmarshalYAML(%q) = %v, want %v", tc.input, got, tc.want)
			}
		})
	}
}

func TestBuild(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		fallbacks  fallbacks.Fallbacks
		sourceType sources.Type
		want       []sources.Type
	}{
		{
			name:       "prepends sourceType and deduplicates",
			fallbacks:  fallbacks.Fallbacks{sources.GITLAB, sources.GITHUB},
			sourceType: sources.GITHUB,
			// Build: prepend GITHUB → [GITHUB, GITLAB, GITHUB] → compact → [GITHUB, GITLAB]
			want: []sources.Type{sources.GITHUB, sources.GITLAB},
		},
		{
			name:       "no duplicates when sourceType not in fallbacks",
			fallbacks:  fallbacks.Fallbacks{sources.GITLAB, sources.URL},
			sourceType: sources.GITHUB,
			// Build: [GITHUB, GITLAB, URL] → compact → [GITHUB, GITLAB, URL]
			want: []sources.Type{sources.GITHUB, sources.GITLAB, sources.URL},
		},
		{
			name:       "nil fallbacks returns only sourceType",
			fallbacks:  nil,
			sourceType: sources.GITLAB,
			want:       []sources.Type{sources.GITLAB},
		},
		{
			name:       "all fallbacks same as sourceType deduplicates to one",
			fallbacks:  fallbacks.Fallbacks{sources.GITHUB, sources.GITHUB},
			sourceType: sources.GITHUB,
			// Build: [GITHUB, GITHUB, GITHUB] → compact → [GITHUB]
			want: []sources.Type{sources.GITHUB},
		},
		{
			name:       "sourceType appears at end of fallbacks is moved to front",
			fallbacks:  fallbacks.Fallbacks{sources.GITLAB, sources.URL, sources.GITHUB},
			sourceType: sources.GITHUB,
			// Build: prepend GITHUB → [GITHUB, GITLAB, URL, GITHUB] → compact → [GITHUB, GITLAB, URL]
			want: []sources.Type{sources.GITHUB, sources.GITLAB, sources.URL},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := tc.fallbacks.Build(tc.sourceType)

			if !slices.Equal(got, tc.want) {
				t.Errorf("Build(%q) = %v, want %v", tc.sourceType, got, tc.want)
			}
		})
	}
}
