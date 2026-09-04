//nolint:testpackage // Tests exercise suggestion prompt construction and deterministic verification.
package aimatch

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/idelchi/godyl/internal/detect"
	"github.com/idelchi/godyl/internal/match"
	"github.com/idelchi/godyl/internal/release"
	"github.com/idelchi/godyl/internal/tools/hints"
)

func TestBuildSuggestionPromptUsesAmbiguousShortlist(t *testing.T) {
	t.Parallel()

	failure := testSelectionFailure(t, release.Assets{
		{Name: "tool-preferred-a-linux-amd64.tar.gz"},
		{Name: "tool-preferred-b-linux-amd64.zip"},
		{Name: "tool-other-linux-amd64.tar.gz"},
	}, match.Requirements{
		Platform: testPlatform(),
		Hints:    hints.Hints{testHint(t, "preferred", hints.Weighted)},
	})

	prompt, candidates, err := buildSuggestionPrompt(SuggestionRequest{Target: "tool", Failure: failure})
	if err != nil {
		t.Fatalf("buildSuggestionPrompt() unexpected error: %v", err)
	}

	if len(candidates) != 2 {
		t.Fatalf("buildSuggestionPrompt() candidates = %v, want two tied best candidates", candidates)
	}

	var request suggestionPrompt

	if err := json.Unmarshal([]byte(prompt), &request); err != nil {
		t.Fatalf("unmarshalling prompt: %v", err)
	}

	if request.Failure != "ambiguous" {
		t.Errorf("buildSuggestionPrompt() failure = %q, want %q", request.Failure, "ambiguous")
	}

	if len(request.Candidates) != 2 {
		t.Errorf("buildSuggestionPrompt() prompt candidates = %v, want two candidates", request.Candidates)
	}
}

func TestBuildSuggestionPromptUsesFullListForNoQualifiedMatch(t *testing.T) {
	t.Parallel()

	assets := release.Assets{{Name: "tool-linux-amd64.tar.gz"}, {Name: "tool-linux-amd64.zip"}}
	failure := testSelectionFailure(t, assets, match.Requirements{
		Platform: testPlatform(),
		Hints:    hints.Hints{testHint(t, "missing", hints.Required)},
	})

	_, candidates, err := buildSuggestionPrompt(SuggestionRequest{Target: "tool", Failure: failure})
	if err != nil {
		t.Fatalf("buildSuggestionPrompt() unexpected error: %v", err)
	}

	if len(candidates) != len(assets) {
		t.Errorf("buildSuggestionPrompt() candidates = %v, want full asset list", candidates)
	}
}

func TestBuildSuggestionPromptRejectsEmptyAssetList(t *testing.T) {
	t.Parallel()

	failure := testSelectionFailure(t, nil, match.Requirements{})

	if _, _, err := buildSuggestionPrompt(SuggestionRequest{Target: "tool", Failure: failure}); err == nil {
		t.Fatal("buildSuggestionPrompt() expected an error for an empty release")
	}
}

func TestVerifySuggestionRerunsDeterministicMatcher(t *testing.T) {
	t.Parallel()

	assets := release.Assets{{Name: "tool-linux-amd64.tar.gz"}, {Name: "tool-linux-amd64.zip"}}
	failure := testSelectionFailure(t, assets, match.Requirements{
		Platform: testPlatform(),
		Hints:    hints.Hints{testHint(t, "missing", hints.Required)},
	})

	suggestion := Suggestion{
		SelectedAsset: assets[0].Name,
		ProposedHints: []SuggestedHint{{
			Pattern: ".tar.gz",
			Type:    string(hints.EndsWith),
			Match:   string(hints.Weighted),
			Weight:  1,
		}},
	}

	if err := verifySuggestion(t.Context(), suggestion, failure); err != nil {
		t.Fatalf("verifySuggestion() unexpected error: %v", err)
	}
}

func TestVerifySuggestionRejectsUnverifiableHints(t *testing.T) {
	t.Parallel()

	assets := release.Assets{{Name: "tool-linux-amd64.tar.gz"}, {Name: "tool-linux-amd64.zip"}}
	failure := testSelectionFailure(t, assets, match.Requirements{
		Platform: testPlatform(),
		Hints:    hints.Hints{testHint(t, "missing", hints.Required)},
	})

	suggestion := Suggestion{
		SelectedAsset: assets[0].Name,
		ProposedHints: []SuggestedHint{{
			Pattern: "linux-amd64",
			Type:    string(hints.Contains),
			Match:   string(hints.Weighted),
			Weight:  1,
		}},
	}

	if err := verifySuggestion(t.Context(), suggestion, failure); err == nil {
		t.Fatal("verifySuggestion() expected an error for hints that remain ambiguous")
	}
}

func testSelectionFailure(
	t *testing.T,
	assets release.Assets,
	requirements match.Requirements,
) *match.SelectionError {
	t.Helper()

	_, err := assets.Select(context.Background(), requirements)
	if err == nil {
		t.Fatal("Select() expected a matching failure")
	}

	var failure *match.SelectionError

	if !errors.As(err, &failure) {
		t.Fatalf("Select() error = %T, want *match.SelectionError", err)
	}

	return failure
}

func testPlatform() detect.Platform {
	var platform detect.Platform

	platform.ParseFrom("target-linux-amd64")

	return platform
}

func testHint(t *testing.T, pattern string, matchType hints.Match) hints.Hint {
	t.Helper()

	hint := hints.Hint{Pattern: pattern, Type: hints.Contains}
	hint.Match.Set(string(matchType))

	if err := hint.Parse(); err != nil {
		t.Fatalf("Parse() unexpected error: %v", err)
	}

	return hint
}
