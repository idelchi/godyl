package match_test

import (
	"errors"
	"slices"
	"strings"
	"testing"

	"github.com/idelchi/godyl/internal/match"
)

func makeResult(score int, qualified bool) match.Result {
	return match.Result{Score: score, Qualified: qualified}
}

func TestBest(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		input     match.Results
		wantLen   int
		wantScore int
	}{
		{
			name: "unique best score among qualified",
			input: match.Results{
				makeResult(5, true),
				makeResult(3, true),
				makeResult(1, true),
			},
			wantLen:   1,
			wantScore: 5,
		},
		{
			name: "two results share the best score",
			input: match.Results{
				makeResult(5, true),
				makeResult(5, true),
				makeResult(3, true),
			},
			wantLen:   2,
			wantScore: 5,
		},
		{
			name: "no qualified results",
			input: match.Results{
				makeResult(5, false),
				makeResult(3, false),
			},
			wantLen: 0,
		},
		{
			name:    "empty input",
			input:   match.Results{},
			wantLen: 0,
		},
		{
			name: "all-zero-score qualified returns all qualified",
			input: match.Results{
				makeResult(0, true),
				makeResult(0, true),
			},
			wantLen:   2,
			wantScore: 0,
		},
		{
			// Mixed scores: Best() picks the single highest-scoring qualified result.
			name: "mixed scores returns only highest scoring qualified",
			input: match.Results{
				makeResult(0, true),
				makeResult(1, true),
				makeResult(0, true),
			},
			wantLen:   1,
			wantScore: 1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := tc.input.Best()

			if len(got) != tc.wantLen {
				t.Fatalf("Best() returned %d results, want %d", len(got), tc.wantLen)
			}

			for i, r := range got {
				if r.Score != tc.wantScore {
					t.Errorf("Best()[%d].Score = %d, want %d", i, r.Score, tc.wantScore)
				}
			}
		})
	}
}

func TestStatus(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		input     match.Results
		wantErrIs error
	}{
		{
			name: "single qualified result returns nil",
			input: match.Results{
				makeResult(5, true),
			},
			wantErrIs: nil,
		},
		{
			name: "no qualified results returns ErrNoQualified",
			input: match.Results{
				makeResult(5, false),
				makeResult(3, false),
			},
			wantErrIs: match.ErrNoQualified,
		},
		{
			name:      "empty results returns ErrNoQualified",
			input:     match.Results{},
			wantErrIs: match.ErrNoQualified,
		},
		{
			name: "two qualified results with equal score returns ErrAmbiguous",
			input: match.Results{
				makeResult(5, true),
				makeResult(5, true),
			},
			wantErrIs: match.ErrAmbiguous,
		},
		{
			name: "two qualified results with different scores returns ErrAmbiguous",
			input: match.Results{
				makeResult(5, true),
				makeResult(3, true),
			},
			wantErrIs: match.ErrAmbiguous,
		},
		{
			name: "one qualified one unqualified is ambiguous due to len > 1",
			input: match.Results{
				makeResult(5, true),
				makeResult(3, false),
			},
			wantErrIs: match.ErrAmbiguous,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := tc.input.Status()

			if tc.wantErrIs == nil {
				if err != nil {
					t.Fatalf("Status() = %v, want nil", err)
				}

				return
			}

			if err == nil {
				t.Fatalf("Status() = nil, want errors.Is %v", tc.wantErrIs)
			}

			if !errors.Is(err, tc.wantErrIs) {
				t.Errorf("Status() error = %v, want errors.Is %v", err, tc.wantErrIs)
			}
		})
	}
}

func TestWithoutZero(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		input      match.Results
		wantScores []int
	}{
		{
			name: "filters out zero-score results",
			input: match.Results{
				makeResult(0, true),
				makeResult(1, true),
				makeResult(0, false),
				makeResult(2, true),
				makeResult(7, false),
			},
			wantScores: []int{1, 2, 7},
		},
		{
			name: "all zero scores returns nil",
			input: match.Results{
				makeResult(0, true),
				makeResult(0, false),
			},
			wantScores: nil,
		},
		{
			name:       "empty input returns nil",
			input:      match.Results{},
			wantScores: nil,
		},
		{
			name: "no zero scores returns all",
			input: match.Results{
				makeResult(3, true),
				makeResult(7, false),
			},
			wantScores: []int{3, 7},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := tc.input.WithoutZero()

			gotScores := make([]int, len(got))
			for i, r := range got {
				gotScores[i] = r.Score
			}

			if !slices.Equal(gotScores, tc.wantScores) {
				t.Errorf("WithoutZero() scores = %v, want %v", gotScores, tc.wantScores)
			}
		})
	}
}

func TestSorted(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		input         match.Results
		wantQualified []bool
		wantScores    []int
	}{
		{
			name: "qualified first then score descending",
			input: match.Results{
				makeResult(1, false),
				makeResult(5, true),
				makeResult(10, false),
				makeResult(3, true),
			},
			wantQualified: []bool{true, true, false, false},
			wantScores:    []int{5, 3, 10, 1},
		},
		{
			name: "all qualified sorted by score descending",
			input: match.Results{
				makeResult(2, true),
				makeResult(8, true),
				makeResult(5, true),
			},
			wantQualified: []bool{true, true, true},
			wantScores:    []int{8, 5, 2},
		},
		{
			name:          "empty input returns empty",
			input:         match.Results{},
			wantQualified: []bool{},
			wantScores:    []int{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			original := slices.Clone(tc.input)

			got := tc.input.Sorted()

			// Verify the original slice was not mutated.
			if !slices.Equal(tc.input, original) {
				t.Errorf("Sorted() mutated original slice: got %v, want %v", tc.input, original)
			}

			gotScores := make([]int, len(got))
			gotQualified := make([]bool, len(got))

			for i, r := range got {
				gotScores[i] = r.Score
				gotQualified[i] = r.Qualified
			}

			if !slices.Equal(gotScores, tc.wantScores) {
				t.Errorf("Sorted() scores = %v, want %v", gotScores, tc.wantScores)
			}

			if !slices.Equal(gotQualified, tc.wantQualified) {
				t.Errorf("Sorted() qualified = %v, want %v", gotQualified, tc.wantQualified)
			}
		})
	}
}

func TestSuccess(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input match.Results
		want  bool
	}{
		{
			name:  "single result returns true",
			input: match.Results{makeResult(5, true)},
			want:  true,
		},
		{
			name: "two results returns false",
			input: match.Results{
				makeResult(5, true),
				makeResult(3, true),
			},
			want: false,
		},
		{
			name:  "zero results returns false",
			input: match.Results{},
			want:  false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := tc.input.Success()
			if got != tc.want {
				t.Errorf("Success() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestIsAmbiguous(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input match.Results
		want  bool
	}{
		{
			name:  "0 results returns false",
			input: match.Results{},
			want:  false,
		},
		{
			name:  "1 result returns false",
			input: match.Results{makeResult(5, true)},
			want:  false,
		},
		{
			name: "2 results returns true",
			input: match.Results{
				makeResult(5, true),
				makeResult(3, true),
			},
			want: true,
		},
		{
			name: "3 results returns true",
			input: match.Results{
				makeResult(5, true),
				makeResult(3, true),
				makeResult(1, true),
			},
			want: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := tc.input.IsAmbiguous()
			if got != tc.want {
				t.Errorf("IsAmbiguous() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestHasQualified(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input match.Results
		want  bool
	}{
		{
			name: "all unqualified returns false",
			input: match.Results{
				makeResult(5, false),
				makeResult(3, false),
			},
			want: false,
		},
		{
			name: "one qualified returns true",
			input: match.Results{
				makeResult(5, false),
				makeResult(3, true),
			},
			want: true,
		},
		{
			name:  "empty returns false",
			input: match.Results{},
			want:  false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := tc.input.HasQualified()
			if got != tc.want {
				t.Errorf("HasQualified() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestToString(t *testing.T) {
	t.Parallel()

	t.Run("non-zero asset fields do not panic", func(t *testing.T) {
		t.Parallel()

		r := match.Results{
			{
				Score:     7,
				Qualified: true,
				Asset: match.Asset{
					Name: "tool-linux-amd64.tar.gz",
				},
			},
			{
				Score:     3,
				Qualified: false,
				Asset: match.Asset{
					Name: "tool-darwin-arm64.tar.gz",
				},
			},
		}

		got := r.ToString()

		if got == "" {
			t.Error("ToString() returned empty string for non-empty results")
		}

		const wantName = "tool-linux-amd64.tar.gz"
		if !strings.Contains(got, wantName) {
			t.Errorf("ToString() = %q, want it to contain asset name %q", got, wantName)
		}
	})
}
