package match

// Hints represents a collection of Hint objects used to evaluate asset matches.
type Hints []Hint

// Add appends a set of hints to the current collection.
func (h *Hints) Add(hints Hints) {
	for _, hint := range hints {
		*h = append(*h, hint)
	}
}

// Hint represents a pattern used to match asset names.
// It can be a regular expression or a simple string pattern.
type Hint struct {
	Pattern        string // Pattern to match against the asset's name.
	WeightTemplate string `json:"-"         mapstructure:"weight" yaml:"weight"` // Template for calculating the weight (not used in matching).
	Weight         int    `mapstructure:"-" yaml:"-"`      // Weight used to adjust the score for non-mandatory hints.
	Regex          bool   // Whether the pattern is a regular expression.
	Must           bool   // Indicates if the hint is mandatory for a match.
}

// NewDefaultHint creates a new Hint with the given pattern.
// The hint is mandatory by default and has a weight of 1.
func NewDefaultHint(pattern string) Hint {
	return Hint{
		Pattern: pattern,
		Weight:  1,
		Must:    true,
	}
}
