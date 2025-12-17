package match

import (
	"errors"
	"fmt"
	"slices"
	"sort"
	"strings"
)

var (
	// ErrNoMatch is returned when no matches are found.
	ErrNoMatch = errors.New("no matches found")
	// ErrAmbiguous is returned when multiple equally good matches are found.
	ErrAmbiguous = errors.New("ambiguous matches found")
	// ErrNoQualified is returned when no qualified matches are found.
	ErrNoQualified = errors.New("no qualified matches found")
)

// Result represents the outcome of matching an asset.
// It contains the asset, its score, and whether it is qualified.
type Result struct {
	Error     error
	Asset     Asset
	Score     int
	Qualified bool
}

// Results is a collection of Result objects.
type Results []Result

// ToString converts the results into a formatted string for output.
func (m Results) ToString() string {
	var (
		result     string
		resultSb35 strings.Builder
	)

	for _, res := range m {
		resultSb35.WriteString(fmt.Sprintf("	- %s\n", res.Asset.Name))
		resultSb35.WriteString(fmt.Sprintf("		score: %d\n", res.Score))
		resultSb35.WriteString(fmt.Sprintf("		qualified: %t\n", res.Qualified))
		resultSb35.WriteString("		detected as:\n")
		resultSb35.WriteString(fmt.Sprintf("		  os: %v\n", res.Asset.Platform.OS))
		resultSb35.WriteString(fmt.Sprintf("		  arch: %v\n", res.Asset.Platform.Architecture))
		resultSb35.WriteString(fmt.Sprintf("		  library: %s\n", res.Asset.Platform.Library))
		resultSb35.WriteString(fmt.Sprintf("		  extension: %s\n", res.Asset.Platform.Extension))

		if res.Error != nil {
			resultSb35.WriteString(fmt.Sprintf("		error: %s\n", res.Error))
		}
	}

	result += resultSb35.String()

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
	m = m.Sorted()
	if !m.HasQualified() { //nolint:gocritic // If-else chain is more readable than switch for this logic
		err = ErrNoQualified

		return fmt.Errorf("%w: \n%s%s", err, m.ToString(), "  ** check settings **")
	} else if m.IsAmbiguous() {
		err = ErrAmbiguous

		return fmt.Errorf("%w: \n%s%s", err, m.Best().ToString(), "  ** try to tune weights **")
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

// IsAmbiguous returns true if there are multiple qualified results.
func (m Results) IsAmbiguous() bool {
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

// Errors returns a combined error from all results.
func (m Results) Errors() []error {
	var errs []error //nolint:prealloc // Size unknown as it depends on error count

	for _, result := range m {
		if result.Error == nil {
			continue
		}

		errs = append(errs, result.Error)
	}

	return errs
}

// HasErrors returns true if there are any errors in the results.
func (m Results) HasErrors() bool {
	return len(m.Errors()) > 0
}

// WithoutZero returns a new instance of Results without zero scores.
func (m Results) WithoutZero() Results {
	var qualified Results

	for _, result := range m {
		if result.Score > 0 {
			qualified = append(qualified, result)
		}
	}

	return qualified
}

// Sorted returns a new sorted instance of Results.
// It sorts first by qualification status and then by score in descending order.
func (m Results) Sorted() Results {
	sortedResults := slices.Clone(m)
	sort.Slice(sortedResults, func(i, j int) bool {
		if sortedResults[i].Qualified != sortedResults[j].Qualified {
			return sortedResults[i].Qualified
		}

		return sortedResults[i].Score > sortedResults[j].Score
	})

	return sortedResults
}
