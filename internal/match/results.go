package match

import (
	"fmt"
	"sort"
)

var (
	ErrNoMatch     = fmt.Errorf("no matches found")
	ErrAmbiguous   = fmt.Errorf("ambiguous matches found")
	ErrNoQualified = fmt.Errorf("no qualified matches found")
)

type Result struct {
	Name      string
	Score     int
	Qualified bool
}

type Results []Result

func (m Results) ToString() string {
	var result string
	for _, r := range m {
		result += fmt.Sprintf("	- %s: %d\n", r.Name, r.Score)
	}
	return result
}

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
	} else {
	}

	return nil
}

func (m Results) Success() bool {
	return len(m) == 1
}

func (m Results) IsAmbigious() bool {
	return len(m) > 1
}

// HasQualified returns true if there's any qualified result.
func (m Results) HasQualified() bool {
	for _, result := range m {
		if result.Qualified {
			return true
		}
	}
	return false
}

// Sorted returns a new sorted instance of Results by Qualification and Score.
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
