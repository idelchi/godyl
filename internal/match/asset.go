package match

import (
	"strings"

	"github.com/idelchi/godyl/internal/detect"
	"github.com/idelchi/godyl/internal/tools/hints"
)

// Asset represents a downloadable file, typically for a specific platform.
// It includes metadata like the asset name and the platform it targets.
type Asset struct {
	Name     string          // Name is the file name or identifier of the asset.
	Platform detect.Platform // Platform describes the OS, architecture, and other relevant details for compatibility.
}

// Lower returns the asset's name in lowercase.
// This is useful for case-insensitive matching operations.
func (a *Asset) Lower() string {
	return strings.ToLower(a.Name)
}

// Parse invokes the platform's parsing logic on the asset's name.
// It populates or derives the platform information from the asset's name.
func (a *Asset) Parse() {
	a.Platform.ParseFrom(a.Name)
}

// PlatformMatch evaluates whether the asset's platform matches the required platform.
// It calculates a score based on the degree of compatibility and returns whether the asset is qualified.
func (a *Asset) PlatformMatch(req Requirements) (int, bool) {
	var score int

	qualified := true

	// Match operating system
	if req.Platform.OS.Is(a.Platform.OS) {
		score++
	}

	if req.Platform.OS.IsCompatibleWith(a.Platform.OS) {
		score++
	} else if !a.Platform.OS.IsUnset() && !req.Platform.OS.IsUnset() {
		qualified = false
	}

	if req.Platform.Architecture.Is(a.Platform.Architecture) {
		score++
	}

	switch {
	case req.Platform.Architecture.IsCompatibleWith(a.Platform.Architecture):
		score++
	case req.Platform.OS.Type() == "windows" && req.Platform.Architecture.Is64Bit() && a.Platform.Architecture.IsX86():
		// Special case: on Windows, 32bit binaries can run on 64-bit systems
		score--
	case a.Platform.Architecture.IsSet() && req.Platform.Architecture.IsSet():
		qualified = false
	default:
		score--
	}

	if req.Platform.Library.Is(a.Platform.Library) {
		score++
	}

	if req.Platform.Library.IsCompatibleWith(a.Platform.Library) {
		score++
	} else if a.Platform.Library.IsSet() && req.Platform.Library.IsSet() {
		qualified = false
	}

	return score, qualified
}

// Match evaluates if the asset satisfies the given requirements.
// It aggregates scores from both platform compatibility and matching hints.
func (a *Asset) Match(req Requirements) (int, bool, error) {
	// Check mandatory hints
	for _, hint := range req.Hints {
		match, err := hint.Matches(a.Lower())
		if err != nil {
			return 0, false, err
		}

		if hint.Match.Value == hints.Required && !match {
			return 0, false, nil
		}

		if hint.Match.Value == hints.Excluded && match {
			return 0, false, nil
		}
	}

	// Match platform requirements
	score, qualified := a.PlatformMatch(req)

	// Check non-mandatory hints and adjust the score
	for _, hint := range req.Hints {
		match, err := hint.Matches(a.Lower())
		if err != nil {
			return 0, false, err
		}

		if hint.Match.Value == hints.Weighted && match {
			score += hint.Weight.Value
		}
	}

	return score, qualified, nil
}
