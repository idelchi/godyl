package release_test

import (
	"context"
	"errors"
	"testing"

	"github.com/idelchi/godyl/internal/detect"
	"github.com/idelchi/godyl/internal/match"
	"github.com/idelchi/godyl/internal/release"
	"github.com/idelchi/godyl/internal/tools/hints"
)

type selectorFunc func(context.Context, match.SelectionRequest) (string, error)

func (f selectorFunc) Select(ctx context.Context, request match.SelectionRequest) (string, error) {
	return f(ctx, request)
}

func TestAssetsSelectSkipsFallbackOnDeterministicMatch(t *testing.T) {
	t.Parallel()

	called := false
	requirements := match.Requirements{
		Hints: hints.Hints{parsedHint(t, "only-one", hints.Weighted)},
		Selector: selectorFunc(func(context.Context, match.SelectionRequest) (string, error) {
			called = true

			return "", errors.New("unexpected fallback")
		}),
	}
	assets := release.Assets{{Name: "only-one"}}

	asset, err := assets.Select(t.Context(), requirements)
	if err != nil {
		t.Fatalf("Select() unexpected error: %v", err)
	}

	if called {
		t.Fatal("Select() called fallback after deterministic matching succeeded")
	}

	if asset.Name != "only-one" {
		t.Errorf("Select() = %q, want %q", asset.Name, "only-one")
	}
}

func TestAssetsSelectSkipsFallbackAfterNoQualifiedMatch(t *testing.T) {
	t.Parallel()

	called := false
	requirements := match.Requirements{
		Hints: hints.Hints{parsedHint(t, "missing", hints.Required)},
		Selector: selectorFunc(func(context.Context, match.SelectionRequest) (string, error) {
			called = true

			return "", errors.New("unexpected fallback")
		}),
	}
	assets := release.Assets{{Name: "tool-linux-amd64.tar.gz"}}

	_, err := assets.Select(t.Context(), requirements)
	if !errors.Is(err, match.ErrNoQualified) {
		t.Fatalf("Select() error = %v, want errors.Is ErrNoQualified", err)
	}

	var selectionErr *match.SelectionError

	if !errors.As(err, &selectionErr) {
		t.Fatalf("Select() error = %T, want *match.SelectionError", err)
	}

	if len(selectionErr.Candidates) != 1 || selectionErr.Candidates[0] != assets[0].Name {
		t.Errorf("Select() diagnostic candidates = %v, want %q", selectionErr.Candidates, assets[0].Name)
	}

	if called {
		t.Fatal("Select() called fallback after no qualified match")
	}
}

func TestAssetsSelectReturnsAmbiguousDiagnosticWithoutFallback(t *testing.T) {
	t.Parallel()

	var platform detect.Platform

	platform.ParseFrom("target-linux-amd64")

	assets := release.Assets{
		{Name: "tool-a-linux-amd64.tar.gz"},
		{Name: "tool-b-linux-amd64.zip"},
	}

	_, err := assets.Select(t.Context(), match.Requirements{Platform: platform})
	if !errors.Is(err, match.ErrAmbiguous) {
		t.Fatalf("Select() error = %v, want errors.Is ErrAmbiguous", err)
	}

	var selectionErr *match.SelectionError

	if !errors.As(err, &selectionErr) {
		t.Fatalf("Select() error = %T, want *match.SelectionError", err)
	}

	if len(selectionErr.Results) != len(assets) {
		t.Errorf("Select() diagnostic results = %v, want both tied candidates", selectionErr.Results)
	}
}

func TestAssetsSelectOffersOnlyBestMatchesToFallback(t *testing.T) {
	t.Parallel()

	var platform detect.Platform

	platform.ParseFrom("target-linux-amd64")

	requirements := match.Requirements{
		Platform: platform,
		Hints:    hints.Hints{parsedHint(t, "preferred", hints.Weighted)},
		Selector: selectorFunc(func(_ context.Context, request match.SelectionRequest) (string, error) {
			if len(request.Candidates) != 2 {
				t.Errorf("Select() candidates = %v, want the two tied best matches", request.Candidates)
			}

			for _, candidate := range request.Candidates {
				if candidate == "tool-other-linux-amd64.tar.gz" {
					t.Errorf("Select() candidates unexpectedly contains lower-scored asset %q", candidate)
				}
			}

			return "tool-preferred-a-linux-amd64.tar.gz", nil
		}),
	}
	assets := release.Assets{
		{Name: "tool-preferred-a-linux-amd64.tar.gz"},
		{Name: "tool-preferred-b-linux-amd64.zip"},
		{Name: "tool-other-linux-amd64.tar.gz"},
	}

	asset, err := assets.Select(t.Context(), requirements)
	if err != nil {
		t.Fatalf("Select() unexpected error: %v", err)
	}

	if asset.Name != "tool-preferred-a-linux-amd64.tar.gz" {
		t.Errorf("Select() = %q, want %q", asset.Name, "tool-preferred-a-linux-amd64.tar.gz")
	}
}

func TestAssetsSelectRejectsCandidateOutsideAmbiguousSet(t *testing.T) {
	t.Parallel()

	var platform detect.Platform

	platform.ParseFrom("target-linux-amd64")

	requirements := match.Requirements{
		Platform: platform,
		Hints:    hints.Hints{parsedHint(t, "preferred", hints.Weighted)},
		Selector: selectorFunc(func(context.Context, match.SelectionRequest) (string, error) {
			return "tool-other-linux-amd64.tar.gz", nil
		}),
	}
	assets := release.Assets{
		{Name: "tool-preferred-a-linux-amd64.tar.gz"},
		{Name: "tool-preferred-b-linux-amd64.zip"},
		{Name: "tool-other-linux-amd64.tar.gz"},
	}

	_, err := assets.Select(t.Context(), requirements)
	if err == nil {
		t.Fatal("Select() expected an error for an asset outside the ambiguous set")
	}

	if !errors.Is(err, match.ErrAmbiguous) {
		t.Fatalf("Select() error = %v, want errors.Is ErrAmbiguous", err)
	}
}

func parsedHint(t *testing.T, pattern string, matchType hints.Match) hints.Hint {
	t.Helper()

	hint := hints.Hint{Pattern: pattern, Type: hints.Contains}
	hint.Match.Set(string(matchType))

	if err := hint.Parse(); err != nil {
		t.Fatalf("Parse() unexpected error: %v", err)
	}

	return hint
}
