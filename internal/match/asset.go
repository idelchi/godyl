package match

import (
	"regexp"
	"strings"

	"github.com/idelchi/godyl/internal/detect"
)

// Asset represents a downloadable file, typically for a specific platform.
// It includes metadata like the asset name and the platform it targets.
type Asset struct {
	Name     string          // Name is the file name or identifier of the asset.
	Platform detect.Platform // Platform describes the OS, architecture, and other relevant details for compatibility.
}

// NameLower returns the asset's name in lowercase.
// This is useful for case-insensitive matching operations.
func (a Asset) NameLower() string {
	return strings.ToLower(a.Name)
}

// Parse invokes the platform's parsing logic on the asset's name.
// It populates or derives the platform information from the asset's name.
func (a *Asset) Parse() {
	a.Platform.Parse(a.Name)
}

// MatchHint checks if the asset's name matches the provided hint.
// The hint can be a regular expression or a simple substring match.
func (a Asset) MatchHint(hint Hint) bool {
	regex, err := regexp.Compile(hint.Pattern)
	return err == nil && regex.MatchString(a.NameLower())
}

// PlatformMatch evaluates whether the asset's platform matches the required platform.
// It calculates a score based on the degree of compatibility and returns whether the asset is qualified.
func (a Asset) PlatformMatch(req Requirements) (int, bool) {
	var score int
	qualified := true

	// Match operating system
	if a.Platform.OS != "" {
		if a.Platform.OS == req.Platform.OS {
			score++
		}
		if req.Platform.OS.IsCompatibleWith(a.Platform.OS.Name()) {
			score++
		} else {
			qualified = false
		}
	}

	// Match architecture
	if a.Platform.Architecture.Name() != "" {
		if a.Platform.Architecture.Name() == req.Platform.Architecture.Name() {
			score++
		}
		if req.Platform.Architecture.IsCompatibleWith(a.Platform.Architecture.Name(), req.Platform.Distribution) {
			score++
		} else {
			qualified = false
		}
	}

	// Match library (e.g., runtime or linking library)
	if a.Platform.Library != "" {
		if a.Platform.Library == req.Platform.Library {
			score++
		}
		if req.Platform.Library.IsCompatibleWith(a.Platform.Library.Name()) {
			score++
		} else {
			score-- // Negative score for incompatible library
		}
	}

	return score, qualified
}

// Match evaluates if the asset satisfies the given requirements.
// It aggregates scores from both platform compatibility and matching hints.
func (a Asset) Match(req Requirements) (int, bool) {
	var score int
	qualified := true

	// Check mandatory hints
	for _, hint := range req.Hints {
		if hint.Must && !a.MatchHint(hint) {
			return 0, false
		}
	}

	// Match platform requirements
	score, qualified = a.PlatformMatch(req)

	// Check non-mandatory hints and adjust the score
	for _, hint := range req.Hints {
		if !hint.Must && a.MatchHint(hint) {
			score += hint.WeightInt
		}
	}

	return score, qualified
}
