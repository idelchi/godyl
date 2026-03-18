package match_test

import (
	"errors"
	"strconv"
	"testing"

	"github.com/idelchi/godyl/internal/match"
	"github.com/idelchi/godyl/internal/tools/hints"
)

// zeroReq returns a Requirements with no platform constraints and no hints.
// With a zero platform, PlatformMatch sets qualified=true and score=-1 for any
// asset whose platform is also unset (the default: score-- hits the default branch
// in the architecture switch).
func zeroReq() match.Requirements {
	return match.Requirements{}
}

// mustParseHint is a test helper that creates a Hint that has been through
// Parse() so Match.Value and Weight.Value are populated and ready for use in Match().
func mustParseHint(t *testing.T, pattern string, typ hints.Type, matchKind hints.Match, weight int) hints.Hint {
	t.Helper()

	h := hints.Hint{Pattern: pattern, Type: typ}
	h.Match.Set(string(matchKind))

	if weight != 1 {
		h.Weight.Set(strconv.Itoa(weight))
	}

	if err := h.Parse(); err != nil {
		t.Fatalf("mustParseHint: Parse() failed: %v", err)
	}

	return h
}

// TestAssetsFromNames verifies that FromNames produces an Assets slice whose
// elements carry exactly the provided names and nothing else.
func TestAssetsFromNames(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		input     []string
		wantNames []string
	}{
		{
			name:      "empty names returns empty Assets",
			input:     []string{},
			wantNames: []string{},
		},
		{
			name:      "single name produces one Asset",
			input:     []string{"tool-linux-amd64.tar.gz"},
			wantNames: []string{"tool-linux-amd64.tar.gz"},
		},
		{
			name:      "multiple names produce multiple Assets in order",
			input:     []string{"tool-linux-amd64.tar.gz", "tool-darwin-arm64.tar.gz", "tool-windows-amd64.zip"},
			wantNames: []string{"tool-linux-amd64.tar.gz", "tool-darwin-arm64.tar.gz", "tool-windows-amd64.zip"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var as match.Assets

			got := as.FromNames(tc.input...)

			if len(got) != len(tc.wantNames) {
				t.Fatalf("FromNames() len = %d, want %d", len(got), len(tc.wantNames))
			}

			for i, a := range got {
				if a.Name != tc.wantNames[i] {
					t.Errorf("FromNames()[%d].Name = %q, want %q", i, a.Name, tc.wantNames[i])
				}
			}
		})
	}
}

// hintSpec captures hint parameters for deferred construction inside t.Run.
type hintSpec struct {
	pattern   string
	typ       hints.Type
	matchKind hints.Match
	weight    int
}

// TestAssetsMatch verifies that Match returns one Result per asset and that
// each Result references the correct asset name.
func TestAssetsMatch(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		assets    match.Assets
		hints     []hintSpec
		wantNames []string
	}{
		{
			name:      "empty Assets returns empty Results",
			assets:    match.Assets{},
			wantNames: []string{},
		},
		{
			name: "single asset produces one Result with matching name",
			assets: match.Assets{
				{Name: "tool-linux-amd64.tar.gz"},
			},
			wantNames: []string{"tool-linux-amd64.tar.gz"},
		},
		{
			name: "multiple assets produce Results in the same order",
			assets: match.Assets{
				{Name: "alpha"},
				{Name: "beta"},
				{Name: "gamma"},
			},
			wantNames: []string{"alpha", "beta", "gamma"},
		},
		{
			name: "required hint that matches keeps the asset qualified",
			assets: match.Assets{
				{Name: "tool-linux-amd64.tar.gz"},
			},
			hints:     []hintSpec{{"linux", hints.Contains, hints.Required, 1}},
			wantNames: []string{"tool-linux-amd64.tar.gz"},
		},
		{
			name: "required hint that does not match marks asset as not qualified",
			assets: match.Assets{
				{Name: "tool-darwin-arm64.tar.gz"},
			},
			hints:     []hintSpec{{"linux", hints.Contains, hints.Required, 1}},
			wantNames: []string{"tool-darwin-arm64.tar.gz"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			req := zeroReq()

			for _, hs := range tc.hints {
				req.Hints = append(req.Hints, mustParseHint(t, hs.pattern, hs.typ, hs.matchKind, hs.weight))
			}

			results := tc.assets.Match(req)

			if len(results) != len(tc.wantNames) {
				t.Fatalf("Match() len = %d, want %d", len(results), len(tc.wantNames))
			}

			for i, r := range results {
				if r.Asset.Name != tc.wantNames[i] {
					t.Errorf("Match()[%d].Asset.Name = %q, want %q", i, r.Asset.Name, tc.wantNames[i])
				}
			}
		})
	}
}

// TestAssetsMatchQualification verifies the Qualified field in results when
// hints control eligibility.
func TestAssetsMatchQualification(t *testing.T) {
	t.Parallel()

	t.Run("excluded hint that matches marks asset as not qualified", func(t *testing.T) {
		t.Parallel()

		as := match.Assets{{Name: "tool-linux-amd64.tar.gz"}}
		req := match.Requirements{
			Hints: hints.Hints{
				mustParseHint(t, "linux", hints.Contains, hints.Excluded, 1),
			},
		}

		results := as.Match(req)
		if len(results) != 1 {
			t.Fatalf("Match() len = %d, want 1", len(results))
		}

		if results[0].Qualified {
			t.Errorf("Match()[0].Qualified = true, want false for excluded hint that matches")
		}
	})

	t.Run("weighted hint that matches increases score", func(t *testing.T) {
		t.Parallel()

		as := match.Assets{{Name: "tool-linux-amd64.tar.gz"}}

		baseResults := as.Match(zeroReq())
		if len(baseResults) != 1 {
			t.Fatalf("base Match() len = %d, want 1", len(baseResults))
		}

		baseScore := baseResults[0].Score

		req := match.Requirements{
			Hints: hints.Hints{
				mustParseHint(t, "linux", hints.Contains, hints.Weighted, 1),
			},
		}

		hintResults := as.Match(req)
		if len(hintResults) != 1 {
			t.Fatalf("hint Match() len = %d, want 1", len(hintResults))
		}

		if hintResults[0].Score <= baseScore {
			t.Errorf("Match() score with matching weighted hint = %d, want > base score %d",
				hintResults[0].Score, baseScore)
		}
	})
}

// TestAssetsSelect verifies that Select returns the best matching Results or
// propagates an appropriate error state.
func TestAssetsSelect(t *testing.T) {
	t.Parallel()

	t.Run("empty Assets returns empty Results with ErrNoQualified via Status", func(t *testing.T) {
		t.Parallel()

		as := match.Assets{}
		results := as.Select(zeroReq())

		if err := results.Status(); !errors.Is(err, match.ErrNoQualified) {
			t.Errorf("Select() Status() = %v, want errors.Is ErrNoQualified", err)
		}
	})

	// With a zero platform, PlatformMatch gives qualified=true but score=-1 (the
	// architecture default branch decrements). Best() initialises bestScore at 0
	// and only appends results whose score >= 0 (via > 0 or == 0). A single
	// weighted hint that matches the asset name lifts score to 0, landing in the
	// == bestScore branch so Best() returns it.
	t.Run("single asset with weighted hint returns one Result", func(t *testing.T) {
		t.Parallel()

		req := match.Requirements{
			Hints: hints.Hints{
				mustParseHint(t, "only-one", hints.Contains, hints.Weighted, 1),
			},
		}
		as := match.Assets{{Name: "only-one"}}
		results := as.Select(req)

		if len(results) != 1 {
			t.Fatalf("Select() len = %d, want 1", len(results))
		}

		if results[0].Asset.Name != "only-one" {
			t.Errorf("Select()[0].Asset.Name = %q, want %q", results[0].Asset.Name, "only-one")
		}
	})

	// Two weighted hints that both match the asset push score to +1 (−1 + 1 + 1),
	// which is unambiguous since only one asset qualifies.
	t.Run("required hint selects only matching asset", func(t *testing.T) {
		t.Parallel()

		req := match.Requirements{
			Hints: hints.Hints{
				// Required: only the linux asset passes this gate.
				mustParseHint(t, "linux", hints.Contains, hints.Required, 1),
				// Weighted: adds +1 so the passing asset reaches score 0 in Best().
				mustParseHint(t, "linux", hints.Contains, hints.Weighted, 1),
			},
		}
		as := match.Assets{
			{Name: "tool-linux-amd64.tar.gz"},
			{Name: "tool-darwin-arm64.tar.gz"},
		}

		results := as.Select(req)

		// Only the linux asset satisfies the required hint.
		if len(results) != 1 {
			t.Fatalf("Select() len = %d, want 1", len(results))
		}

		if results[0].Asset.Name != "tool-linux-amd64.tar.gz" {
			t.Errorf("Select()[0].Asset.Name = %q, want %q", results[0].Asset.Name, "tool-linux-amd64.tar.gz")
		}
	})

	t.Run("ambiguous results when two assets score equally", func(t *testing.T) {
		t.Parallel()

		// Both assets contain "tool", so the weighted hint gives each the same
		// score delta. Both are qualified and equal-scoring → IsAmbiguous.
		req := match.Requirements{
			Hints: hints.Hints{
				mustParseHint(t, "tool", hints.Contains, hints.Weighted, 1),
			},
		}
		as := match.Assets{
			{Name: "tool-linux-amd64.tar.gz"},
			{Name: "tool-darwin-arm64.tar.gz"},
		}

		results := as.Select(req)

		if !results.IsAmbiguous() {
			t.Errorf("Select() IsAmbiguous() = false, want true for equally-scored assets")
		}
	})
}

// TestAssetsMatchParseAndMatch exercises the full Parse→Match path: a Hint is
// constructed, Parse() is called, and then Match() is used through Assets.Match().
func TestAssetsMatchParseAndMatch(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		assetName     string
		pattern       string
		typ           hints.Type
		matchKind     hints.Match
		weight        int
		wantQualified bool
	}{
		{
			name:          "required contains hint matches asset",
			assetName:     "tool-linux-amd64.tar.gz",
			pattern:       "linux",
			typ:           hints.Contains,
			matchKind:     hints.Required,
			weight:        1,
			wantQualified: true,
		},
		{
			name:          "required contains hint does not match asset",
			assetName:     "tool-darwin-arm64.tar.gz",
			pattern:       "linux",
			typ:           hints.Contains,
			matchKind:     hints.Required,
			weight:        1,
			wantQualified: false,
		},
		{
			name:          "excluded contains hint matches so asset is disqualified",
			assetName:     "tool-windows-amd64.zip",
			pattern:       "windows",
			typ:           hints.Contains,
			matchKind:     hints.Excluded,
			weight:        1,
			wantQualified: false,
		},
		{
			name:          "weighted glob hint matching asset qualifies it",
			assetName:     "tool-linux-amd64.tar.gz",
			pattern:       "*.tar.gz",
			typ:           hints.Glob,
			matchKind:     hints.Weighted,
			weight:        1,
			wantQualified: true,
		},
		{
			name:          "required regex hint matching asset qualifies it",
			assetName:     "tool-v1.2.3-linux-amd64",
			pattern:       `v\d+\.\d+\.\d+`,
			typ:           hints.Regex,
			matchKind:     hints.Required,
			weight:        1,
			wantQualified: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			h := mustParseHint(t, tc.pattern, tc.typ, tc.matchKind, tc.weight)
			as := match.Assets{{Name: tc.assetName}}
			req := match.Requirements{
				Hints: hints.Hints{h},
			}

			results := as.Match(req)
			if len(results) != 1 {
				t.Fatalf("Match() len = %d, want 1", len(results))
			}

			if results[0].Qualified != tc.wantQualified {
				t.Errorf("Match()[0].Qualified = %v, want %v (asset=%q, pattern=%q, type=%q, match=%q)",
					results[0].Qualified, tc.wantQualified, tc.assetName, tc.pattern, tc.typ, tc.matchKind)
			}
		})
	}
}

// TestResultsErrors verifies that Errors() and HasErrors() correctly reflect
// the presence or absence of per-result errors.
func TestResultsErrors(t *testing.T) {
	t.Parallel()

	t.Run("no errors returns empty slice and HasErrors false", func(t *testing.T) {
		t.Parallel()

		results := match.Results{
			{Score: 5, Qualified: true},
			{Score: 3, Qualified: false},
		}

		errs := results.Errors()
		if len(errs) != 0 {
			t.Errorf("Errors() len = %d, want 0", len(errs))
		}

		if results.HasErrors() {
			t.Error("HasErrors() = true, want false")
		}
	})

	t.Run("one error returns that error and HasErrors true", func(t *testing.T) {
		t.Parallel()

		sentinel := errors.New("something went wrong")
		results := match.Results{
			{Score: 5, Qualified: true, Error: sentinel},
			{Score: 3, Qualified: false},
		}

		errs := results.Errors()
		if len(errs) != 1 {
			t.Fatalf("Errors() len = %d, want 1", len(errs))
		}

		if !errors.Is(errs[0], sentinel) {
			t.Errorf("Errors()[0] = %v, want sentinel error", errs[0])
		}

		if !results.HasErrors() {
			t.Error("HasErrors() = false, want true")
		}
	})

	t.Run("multiple errors all returned", func(t *testing.T) {
		t.Parallel()

		err1 := errors.New("first error")
		err2 := errors.New("second error")

		results := match.Results{
			{Error: err1},
			{Error: nil},
			{Error: err2},
		}

		errs := results.Errors()
		if len(errs) != 2 {
			t.Fatalf("Errors() len = %d, want 2", len(errs))
		}

		if !errors.Is(errs[0], err1) {
			t.Errorf("Errors()[0] = %v, want err1", errs[0])
		}

		if !errors.Is(errs[1], err2) {
			t.Errorf("Errors()[1] = %v, want err2", errs[1])
		}
	})

	t.Run("empty Results has no errors", func(t *testing.T) {
		t.Parallel()

		results := match.Results{}

		if results.HasErrors() {
			t.Error("HasErrors() = true for empty Results, want false")
		}

		if len(results.Errors()) != 0 {
			t.Errorf("Errors() len = %d, want 0 for empty Results", len(results.Errors()))
		}
	})
}
