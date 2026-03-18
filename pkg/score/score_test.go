package score_test

import (
	"slices"
	"testing"

	"github.com/idelchi/godyl/pkg/score"
)

// scoreExpect captures (Item, Score) pairs for comparison without
// depending on the anonymous struct layout of score.Scores.
type scoreExpect struct {
	item  int
	score int
}

func extractScores(s score.Scores[int]) []scoreExpect {
	out := make([]scoreExpect, len(s))
	for i, e := range s {
		out[i] = scoreExpect{item: e.Item, score: e.Score}
	}

	return out
}

func TestScore(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		items   []int
		scorers []func(int) int
		want    []scoreExpect
	}{
		{
			name:    "single identity scorer",
			items:   []int{1, 2, 3},
			scorers: []func(int) int{func(x int) int { return x }},
			want:    []scoreExpect{{1, 1}, {2, 2}, {3, 3}},
		},
		{
			name:    "constant scorer",
			items:   []int{1, 2, 3},
			scorers: []func(int) int{func(int) int { return 5 }},
			want:    []scoreExpect{{1, 5}, {2, 5}, {3, 5}},
		},
		{
			name:  "two scorers summed",
			items: []int{1, 2},
			scorers: []func(int) int{
				func(x int) int { return x },
				func(x int) int { return x * 10 },
			},
			want: []scoreExpect{{1, 11}, {2, 22}},
		},
		{
			name:    "empty items returns nil",
			items:   []int{},
			scorers: []func(int) int{func(x int) int { return x }},
			want:    nil,
		},
		{
			name:    "no scorers yields zero scores",
			items:   []int{1, 2, 3},
			scorers: nil,
			want:    []scoreExpect{{1, 0}, {2, 0}, {3, 0}},
		},
		{
			name:    "scorer returns negative values",
			items:   []int{1, 2, 3},
			scorers: []func(int) int{func(x int) int { return -x }},
			want:    []scoreExpect{{1, -1}, {2, -2}, {3, -3}},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := score.Score(tc.items, tc.scorers...)

			if tc.want == nil {
				if got != nil {
					t.Errorf("Score() = %v, want nil", got)
				}

				return
			}

			gotExpect := extractScores(got)

			if !slices.Equal(gotExpect, tc.want) {
				t.Errorf("Score() = %v, want %v", gotExpect, tc.want)
			}
		})
	}
}

func TestTop(t *testing.T) {
	t.Parallel()

	// makeScores constructs Scores[int] directly as struct literals,
	// without calling Score(), to avoid depending on scorer logic here.
	makeScores := func(pairs [][2]int) score.Scores[int] {
		s := make(score.Scores[int], len(pairs))
		for i, p := range pairs {
			s[i] = struct {
				Item  int
				Score int
			}{Item: p[0], Score: p[1]}
		}

		return s
	}

	tests := []struct {
		name      string
		input     [][2]int // [item, score] pairs
		wantLen   int
		wantScore int
	}{
		{
			name:      "single maximum",
			input:     [][2]int{{1, 1}, {2, 5}, {3, 3}},
			wantLen:   1,
			wantScore: 5,
		},
		{
			name:      "two tied maximums",
			input:     [][2]int{{1, 5}, {2, 1}, {3, 5}},
			wantLen:   2,
			wantScore: 5,
		},
		{
			name:      "all tied",
			input:     [][2]int{{1, 3}, {2, 3}, {3, 3}},
			wantLen:   3,
			wantScore: 3,
		},
		{
			name:    "empty input returns nil",
			input:   nil,
			wantLen: 0,
		},
		{
			name:      "all negative scores",
			input:     [][2]int{{1, -3}, {2, -1}, {3, -2}},
			wantLen:   1,
			wantScore: -1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			s := makeScores(tc.input)
			got := s.Top()

			if len(got) != tc.wantLen {
				t.Fatalf("Top() len = %d, want %d; entries: %v", len(got), tc.wantLen, extractScores(got))
			}

			if tc.wantLen == 0 {
				if got != nil {
					t.Errorf("Top() = %v, want nil", got)
				}

				return
			}

			for i, e := range got {
				if e.Score != tc.wantScore {
					t.Errorf("Top()[%d].Score = %d, want %d", i, e.Score, tc.wantScore)
				}
			}
		})
	}
}

func TestTopSingleElement(t *testing.T) {
	t.Parallel()

	// A single-element Scores slice must return that single element as Top.
	s := score.Score([]int{42}, func(x int) int { return x })

	top := s.Top()

	if len(top) != 1 {
		t.Fatalf("Top() len = %d on single-element Scores, want 1; got %v", len(top), extractScores(top))
	}

	if top[0].Item != 42 {
		t.Errorf("Top()[0].Item = %d, want 42", top[0].Item)
	}

	if top[0].Score != 42 {
		t.Errorf("Top()[0].Score = %d, want 42", top[0].Score)
	}
}

func TestScoreEmptyNilReturn(t *testing.T) {
	t.Parallel()

	// Score with an empty items slice must return nil, not an empty non-nil
	// slice. This is the documented contract (return nil if len(items)==0).
	got := score.Score([]string{}, func(s string) int { return len(s) })

	if got != nil {
		t.Errorf("Score(empty items) = %v (non-nil), want nil", got)
	}
}
