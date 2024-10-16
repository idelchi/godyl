package match

import (
	"fmt"
	"sort"
)

var (
	// ErrNoMatch is returned when no matches are found.
	ErrNoMatch = fmt.Errorf("no matches found")
	// ErrAmbiguous is returned when multiple equally good matches are found.
	ErrAmbiguous = fmt.Errorf("ambiguous matches found")
	// ErrNoQualified is returned when no qualified matches are found.
	ErrNoQualified = fmt.Errorf("no qualified matches found")
)

// Result represents the outcome of matching an asset.
// It contains the asset's name, its score, and whether it is qualified.
type Result struct {
	Name      string // Name of the asset.
	Score     int    // Score representing how well the asset matches the requirements.
	Qualified bool   // Qualified indicates whether the asset meets the necessary criteria.
}

// Results is a collection of Result objects.
type Results []Result

// ToString converts the results into a formatted string for output.
func (m Results) ToString() string {
	var result string
	for _, r := range m {
		result += fmt.Sprintf("	- %s: %d\n", r.Name, r.Score)
	}
	return result
}

// Best returns the best qualified results based on the highest score.
// If multiple results have the same best score, they are all returned.
func (m Results) Best() Results {
	var best Results
	var bestScore int
	for _, result := range m {
		if result.Qualified {
			if result.Score > bestScore {
				best = Results{result}
				bestScore = result.Score
			} else if result.Score == bestScore {
				best = append(best, result)
			}
		}
	}
	return best
}

// Status evaluates the state of the results and returns an appropriate error if needed.
func (m Results) Status() (err error) {
	if !m.HasQualified() {
		err = ErrNoQualified
		return fmt.Errorf("%w: \n%s%s", err, m.ToString(), "  ** check settings **")
	} else if m.IsAmbigious() {
		err = ErrAmbiguous
		return fmt.Errorf("%w: \n%s%s", err, m.ToString(), "  ** try to tune weights **")
	} else if !m.Success() {
		err = ErrNoMatch
		return fmt.Errorf("%w: \n%s%s", err, m.ToString(), "  ** check settings **")
	}
	return nil
}

// Success returns true if there is exactly one result.
func (m Results) Success() bool {
	return len(m) == 1
}

// IsAmbigious returns true if there are multiple qualified results.
func (m Results) IsAmbigious() bool {
	return len(m) > 1
}

// HasQualified returns true if there's at least one qualified result in the set.
func (m Results) HasQualified() bool {
	for _, result := range m {
		if result.Qualified {
			return true
		}
	}
	return false
}

// Sorted returns a new sorted instance of Results.
// It sorts first by qualification status and then by score in descending order.
func (m Results) Sorted() Results {
	sortedResults := append(Results{}, m...) // Create a copy of the original slice
	sort.Slice(sortedResults, func(i, j int) bool {
		if sortedResults[i].Qualified != sortedResults[j].Qualified {
			return sortedResults[i].Qualified
		}
		return sortedResults[i].Score > sortedResults[j].Score
	})
	return sortedResults
}
