//nolint:testpackage // Tests exercise the unexported suggestion dispatch boundary.
package processor

import (
	"context"
	"errors"
	"testing"

	"github.com/idelchi/godyl/internal/aimatch"
	"github.com/idelchi/godyl/internal/match"
	"github.com/idelchi/godyl/internal/tools/tool"
)

type suggesterFunc func(context.Context, aimatch.SuggestionRequest) (aimatch.Suggestion, error)

func (f suggesterFunc) Suggest(
	ctx context.Context,
	request aimatch.SuggestionRequest,
) (aimatch.Suggestion, error) {
	return f(ctx, request)
}

func TestAddSuggestionHandlesOnlyStructuredMatchingFailures(t *testing.T) {
	t.Parallel()

	called := false
	processor := Processor{
		assetSuggester: suggesterFunc(func(
			_ context.Context,
			request aimatch.SuggestionRequest,
		) (aimatch.Suggestion, error) {
			called = true

			if request.Failure == nil {
				t.Fatal("Suggest() matching failure is nil")
			}

			return aimatch.Suggestion{SelectedAsset: "tool.tar.gz"}, nil
		}),
	}
	target := &tool.Tool{Name: "tool"}

	converted := Result{}
	processor.addSuggestion(t.Context(), &converted, target, errors.New("network unavailable"))

	if called {
		t.Fatal("addSuggestion() called the model for a non-matching error")
	}

	matchingErr := &match.SelectionError{Cause: match.ErrNoQualified}
	processor.addSuggestion(t.Context(), &converted, target, matchingErr)

	if !called {
		t.Fatal("addSuggestion() did not call the model for a structured matching error")
	}

	if converted.Suggestion == nil {
		t.Fatal("addSuggestion() did not retain the suggestion")
	}
}
