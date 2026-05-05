package tags_test

import (
	"slices"
	"testing"

	"github.com/goccy/go-yaml"

	"github.com/idelchi/godyl/internal/tools/tags"
)

func TestAppend(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		initial  tags.Tags
		toAppend []string
		wantTags tags.Tags
	}{
		{
			name:     "append new tag to empty list",
			initial:  tags.Tags{},
			toAppend: []string{"linux"},
			wantTags: tags.Tags{"linux"},
		},
		{
			name:     "append duplicate tag is not added",
			initial:  tags.Tags{"linux"},
			toAppend: []string{"linux"},
			wantTags: tags.Tags{"linux"},
		},
		{
			name:     "append multiple new tags",
			initial:  tags.Tags{},
			toAppend: []string{"linux", "amd64", "gnu"},
			wantTags: tags.Tags{"linux", "amd64", "gnu"},
		},
		{
			name:     "append multiple tags some duplicate",
			initial:  tags.Tags{"linux"},
			toAppend: []string{"linux", "amd64"},
			wantTags: tags.Tags{"linux", "amd64"},
		},
		{
			name:     "append nothing to non-empty list",
			initial:  tags.Tags{"linux"},
			toAppend: []string{},
			wantTags: tags.Tags{"linux"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Clone to prevent shared backing array from leaking mutations between subtests.
			got := slices.Clone(tc.initial)
			got.Append(tc.toAppend...)

			if !slices.Equal([]string(got), []string(tc.wantTags)) {
				t.Errorf("Append(%v): got %v, want %v", tc.toAppend, got, tc.wantTags)
			}
		})
	}
}

func TestInclude(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		tool   tags.Tags
		filter tags.Tags
		want   bool
	}{
		{
			name:   "tag present in tool tags returns true",
			tool:   tags.Tags{"linux", "amd64"},
			filter: tags.Tags{"linux"},
			want:   true,
		},
		{
			name:   "tag absent from tool tags returns false",
			tool:   tags.Tags{"linux", "amd64"},
			filter: tags.Tags{"windows"},
			want:   false,
		},
		{
			name:   "wildcard star matches any tool tag",
			tool:   tags.Tags{"linux", "amd64"},
			filter: tags.Tags{"*"},
			want:   true,
		},
		{
			name:   "wildcard prefix pattern matches",
			tool:   tags.Tags{"linux-gnu", "amd64"},
			filter: tags.Tags{"linux-*"},
			want:   true,
		},
		{
			name:   "wildcard suffix pattern matches",
			tool:   tags.Tags{"amd64", "linux-gnu"},
			filter: tags.Tags{"*-gnu"},
			want:   true,
		},
		{
			name:   "wildcard mid-string pattern matches",
			tool:   tags.Tags{"linux-gnu-amd64"},
			filter: tags.Tags{"linux-*-amd64"},
			want:   true,
		},
		{
			name:   "wildcard pattern does not match when no tool tag fits",
			tool:   tags.Tags{"windows", "amd64"},
			filter: tags.Tags{"linux-*"},
			want:   false,
		},
		{
			name:   "empty filter returns true regardless of tool tags",
			tool:   tags.Tags{"linux"},
			filter: tags.Tags{},
			want:   true,
		},
		{
			name:   "empty filter with empty tool tags returns true",
			tool:   tags.Tags{},
			filter: tags.Tags{},
			want:   true,
		},
		{
			name:   "non-empty filter against empty tool tags returns false",
			tool:   tags.Tags{},
			filter: tags.Tags{"linux"},
			want:   false,
		},
		{
			name:   "multiple filter tags one match returns true",
			tool:   tags.Tags{"linux"},
			filter: tags.Tags{"windows", "linux"},
			want:   true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := tc.tool.Include(tc.filter)
			if got != tc.want {
				t.Errorf("Tags(%v).Include(%v): got %v, want %v", tc.tool, tc.filter, got, tc.want)
			}
		})
	}
}

// TestExclude verifies that Exclude is the logical inverse of Include.
// Exhaustive pattern coverage lives in TestInclude; only essential cases are here.
func TestExclude(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		tool   tags.Tags
		filter tags.Tags
		want   bool
	}{
		{
			name:   "filter present returns false",
			tool:   tags.Tags{"linux", "amd64"},
			filter: tags.Tags{"linux"},
			want:   false,
		},
		{
			name:   "filter absent returns true",
			tool:   tags.Tags{"linux", "amd64"},
			filter: tags.Tags{"windows"},
			want:   true,
		},
		{
			// Empty filter: Exclude delegates to Include, which returns true for empty filter.
			// Exclude inverts that, so !true == false... but looking at the source:
			// Exclude returns true when len(tags)==0, matching Include's short-circuit.
			name:   "empty filter returns true",
			tool:   tags.Tags{"linux"},
			filter: tags.Tags{},
			want:   true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := tc.tool.Exclude(tc.filter)
			if got != tc.want {
				t.Errorf("Tags(%v).Exclude(%v): got %v, want %v", tc.tool, tc.filter, got, tc.want)
			}
		})
	}
}

// TestTagsUnmarshalYAMLEdgeCases exercises unusual YAML inputs for Tags:
// null, empty string, integer, and boolean scalars. The goccy/go-yaml library
// coerces non-string scalars to their string representation, so integers and
// booleans are converted rather than rejected. YAML null and empty string both
// produce an empty (or single empty-string) Tags slice.
func TestTagsUnmarshalYAMLEdgeCases(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		wantTags tags.Tags
		wantErr  bool
	}{
		{
			// YAML null (~) is decoded via SingleOrSlice as a single null/empty-string
			// value, yielding a one-element Tags with an empty string. The goccy/go-yaml
			// library decodes the null scalar node into a string as "".
			name:     "null YAML yields single empty-string tag",
			input:    `~`,
			wantTags: tags.Tags{""},
		},
		{
			// An empty quoted string becomes a Tags slice with a single empty string.
			name:     "empty quoted string yields single empty-string tag",
			input:    `""`,
			wantTags: tags.Tags{""},
		},
		{
			// An integer scalar is coerced to its string representation "42".
			name:     "integer scalar is coerced to string tag",
			input:    `42`,
			wantTags: tags.Tags{"42"},
		},
		{
			// A boolean scalar is coerced to its string representation "true".
			name:     "boolean scalar is coerced to string tag",
			input:    `true`,
			wantTags: tags.Tags{"true"},
		},
		{
			// A sequence mixing strings and integers coerces all to strings.
			name:     "mixed sequence of string and integer coerces all to string",
			input:    `["linux", 42]`,
			wantTags: tags.Tags{"linux", "42"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var got tags.Tags

			if err := yaml.Unmarshal([]byte(tc.input), &got); err != nil {
				if tc.wantErr {
					return
				}

				t.Fatalf("yaml.Unmarshal(%q) unexpected error: %v", tc.input, err)
			}

			if tc.wantErr {
				t.Fatalf("yaml.Unmarshal(%q) expected error, got nil", tc.input)
			}

			if !slices.Equal([]string(got), []string(tc.wantTags)) {
				t.Errorf("UnmarshalYAML(%q) = %v, want %v", tc.input, got, tc.wantTags)
			}
		})
	}
}

func TestTagsUnmarshalYAML(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  tags.Tags
	}{
		{
			name:  "scalar string form",
			input: `"linux"`,
			want:  tags.Tags{"linux"},
		},
		{
			name:  "slice form",
			input: `["linux", "amd64"]`,
			want:  tags.Tags{"linux", "amd64"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var got tags.Tags

			if err := yaml.Unmarshal([]byte(tc.input), &got); err != nil {
				t.Fatalf("yaml.Unmarshal(%q) unexpected error: %v", tc.input, err)
			}

			if !slices.Equal([]string(got), []string(tc.want)) {
				t.Errorf("UnmarshalYAML(%q) = %v, want %v", tc.input, got, tc.want)
			}
		})
	}
}
