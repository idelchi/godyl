package match

import (
	"regexp"
	"strings"

	"github.com/idelchi/godyl/internal/detect"
)

type Asset struct {
	Name     string
	Platform detect.Platform
}

func (a Asset) NameLower() string {
	return strings.ToLower(a.Name)
}

func (a *Asset) Parse() {
	a.Platform.Parse(a.Name)
}

func (a Asset) MatchHint(hint Hint) bool {
	if hint.Regex {
		regex, err := regexp.Compile(hint.Pattern)
		return err == nil && regex.MatchString(a.NameLower())
	}
	return strings.Contains(a.NameLower(), hint.Pattern)
}

func (a Asset) PlatformMatch(req Requirements) (int, bool) {
	var score int
	qualified := true

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

	if a.Platform.Library != "" {
		if a.Platform.Library == req.Platform.Library {
			score++
		}
		if req.Platform.Library.IsCompatibleWith(a.Platform.Library.Name()) {
			score++
		} else {
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
