package match

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/idelchi/godyl/internal/detect"
)

type Asset struct {
	Name     string
	Platfrom detect.Platform
}

func (a Asset) NameLower() string {
	return strings.ToLower(a.Name)
}

func (a *Asset) Parse() {
	a.Platfrom.Parse(a.Name)
}

type Hints []Hint

func (h *Hints) Add(hints Hints) {
	for _, hint := range hints {
		*h = append(*h, hint)
	}
}

type Hint struct {
	Pattern        string
	WeightTemplate string `yaml:"weight" mapstructure:"weight"`
	Weight         int    `yaml:"-" mapstructure:"-"`
	Regex          bool
	Must           bool
}

func NewDefaultHint(pattern string) Hint {
	return Hint{
		Pattern: pattern,
		Weight:  1,
		Must:    true,
	}
}

type Requirements struct {
	Platform detect.Platform
	Hints    []Hint
}

func (a Asset) MatchHint(hint Hint) bool {
	if hint.Regex {
		regex, err := regexp.Compile(hint.Pattern)
		return err == nil && regex.MatchString(a.NameLower())
	}
	return strings.Contains(a.NameLower(), hint.Pattern)
}

type Assets []Asset

func (as Assets) FromNames(names ...string) Assets {
	assets := make(Assets, len(names))

	for i, name := range names {
		assets[i] = Asset{Name: name}
	}

	return assets
}

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

var (
	ErrNoMatch     = fmt.Errorf("no matches found")
	ErrAmbiguous   = fmt.Errorf("ambiguous matches found")
	ErrNoQualified = fmt.Errorf("no qualified matches found")
)

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

func (as Assets) Select(req Requirements) Results {
	results := as.Match(req)

	if !results.HasQualified() {
		return results
	}

	return results.Best().Sorted()
}

// Results evaluates all assets against requirements and returns Results.
func (as Assets) Match(req Requirements) Results {
	var results Results
	for _, a := range as {
		score, qualified := a.Match(req)
		results = append(results, Result{Name: a.Name, Score: score, Qualified: qualified})
	}
	return results
}

func (a Asset) PlatformMatch(req Requirements) (int, bool) {
	var score int
	qualified := true

	// fmt.Printf("Asset: %s\n", a.Name)

	if a.Platfrom.OS != "" {
		if a.Platfrom.OS == req.Platform.OS {
			// fmt.Printf("OS %s == %s\n", a.Platfrom.OS, req.Platform.OS)
			score++
		}
		if req.Platform.OS.IsCompatibleWith(a.Platfrom.OS.Name()) {
			// fmt.Printf("OS %s compatible with %s\n", a.Platfrom.OS, req.Platform.OS)
			score++
		} else {
			// fmt.Printf("OS %s not compatible with %s\n", a.Platfrom.OS, req.Platform.OS)
			qualified = false
		}
	}

	if a.Platfrom.Architecture.Name() != "" {
		if a.Platfrom.Architecture.Name() == req.Platform.Architecture.Name() {
			// fmt.Printf("Arch %s == %s\n", a.Platfrom.Architecture, req.Platform.Architecture)
			score++
		}
		if req.Platform.Architecture.IsCompatibleWith(a.Platfrom.Architecture.Name(), req.Platform.Distribution) {
			// fmt.Printf("Arch %s compatible with %s\n", a.Platfrom.Architecture, req.Platform.Architecture)
			score++
		} else {
			// fmt.Printf("Arch %s not compatible with %s\n", a.Platfrom.Architecture, req.Platform.Architecture)
			qualified = false
		}
	}

	if a.Platfrom.Library != "" {
		if a.Platfrom.Library == req.Platform.Library {
			// fmt.Printf("Library %s == %s\n", a.Platfrom.Library, req.Platform.Library)
			score++
		}
		if req.Platform.Library.IsCompatibleWith(a.Platfrom.Library.Name()) {
			// fmt.Printf("Library %s compatible with %s\n", a.Platfrom.Library, req.Platform.Library)
			score++
		} else {
			// fmt.Printf("Library %s not compatible with %s\n", a.Platfrom.Library, req.Platform.Library)
			// qualified = false
			score--
		}
	}

	return score, qualified
}

func (a Asset) Match(req Requirements) (int, bool) {
	var score int
	qualified := true

	for _, hint := range req.Hints {
		if hint.Must && !a.MatchHint(hint) {
			return 0, false
		}
	}

	score, qualified = a.PlatformMatch(req)

	for _, hint := range req.Hints {
		if !hint.Must && a.MatchHint(hint) {
			score += hint.Weight
		}
	}

	return score, qualified
}
