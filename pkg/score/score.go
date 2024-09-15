// Package score provides utilities for scoring and selecting items based on multiple criteria.
package score

import "math"

// Scores holds items with their calculated scores.
type Scores[T any] []struct {
	Item  T
	Score int
}

// Score calculates scores for all items using the provided scoring functions.
func Score[T any](items []T, scorers ...func(T) int) Scores[T] {
	if len(items) == 0 {
		return nil
	}

	scores := make(Scores[T], len(items))

	for i, item := range items {
		score := 0

		for _, scorer := range scorers {
			score += scorer(item)
		}

		scores[i] = struct {
			Item  T
			Score int
		}{item, score}
	}

	return scores
}

// Top returns all items with the highest score.
func (s Scores[T]) Top() Scores[T] {
	if len(s) == 0 {
		return nil
	}

	maxScore := math.MinInt

	for _, scored := range s {
		if scored.Score > maxScore {
			maxScore = scored.Score
		}
	}

	var result Scores[T]

	for _, scored := range s {
		if scored.Score == maxScore {
			result = append(result, scored)
		}
	}

	return result
}
