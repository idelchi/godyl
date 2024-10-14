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

	// fmt.Printf("Asset: %s\n", a.Name)

	if a.Platform.OS != "" {
		if a.Platform.OS == req.Platform.OS {
			// fmt.Printf("OS %s == %s\n", a.Platform.OS, req.Platform.OS)
			score++
		}
		if req.Platform.OS.IsCompatibleWith(a.Platform.OS.Name()) {
			// fmt.Printf("OS %s compatible with %s\n", a.Platform.OS, req.Platform.OS)
			score++
		} else {
			// fmt.Printf("OS %s not compatible with %s\n", a.Platform.OS, req.Platform.OS)
			qualified = false
		}
	}

	if a.Platform.Architecture.Name() != "" {
		if a.Platform.Architecture.Name() == req.Platform.Architecture.Name() {
			// fmt.Printf("Arch %s == %s\n", a.Platform.Architecture, req.Platform.Architecture)
			score++
		}
		if req.Platform.Architecture.IsCompatibleWith(a.Platform.Architecture.Name(), req.Platform.Distribution) {
			// fmt.Printf("Arch %s compatible with %s\n", a.Platform.Architecture, req.Platform.Architecture)
			score++
		} else {
			// fmt.Printf("Arch %s not compatible with %s\n", a.Platform.Architecture, req.Platform.Architecture)
			qualified = false
		}
	}

	if a.Platform.Library != "" {
		if a.Platform.Library == req.Platform.Library {
			// fmt.Printf("Library %s == %s\n", a.Platform.Library, req.Platform.Library)
			score++
		}
		if req.Platform.Library.IsCompatibleWith(a.Platform.Library.Name()) {
			// fmt.Printf("Library %s compatible with %s\n", a.Platform.Library, req.Platform.Library)
			score++
		} else {
			// fmt.Printf("Library %s not compatible with %s\n", a.Platform.Library, req.Platform.Library)
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
