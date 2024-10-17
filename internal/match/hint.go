package match

// Hint represents a pattern used to match asset names.
// It can be a regular expression or a simple string pattern.
type Hint struct {
	Pattern string // Pattern to match against the asset's name.
	Weight  string // Weight used to adjust the score for non-mandatory hints.
	Must    bool   // Indicates if the hint is mandatory for a match.
}

// NewDefaultHint creates a new Hint with the given pattern.
// The hint is mandatory by default and has a weight of 1.
func NewDefaultHint(pattern string) Hint {
	return Hint{
		Pattern: pattern,
		Weight:  "1",
		Must:    true,
	}
}
