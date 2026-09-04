package presentation_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/idelchi/godyl/internal/aimatch"
	"github.com/idelchi/godyl/internal/match"
	"github.com/idelchi/godyl/internal/presentation"
	"github.com/idelchi/godyl/internal/processor"
	"github.com/idelchi/godyl/internal/tools/tool"
)

func TestRenderSuggestions(t *testing.T) {
	t.Parallel()

	output := presentation.RenderSuggestions([]processor.Result{{
		Tool:  &tool.Tool{Name: "tool"},
		Error: &match.SelectionError{Cause: match.ErrAmbiguous},
		Suggestion: &aimatch.Suggestion{
			SelectedAsset: "tool-linux-amd64.tar.gz",
			Confidence:    "high",
			Reason:        "portable archive",
			Candidates:    []string{"tool-linux-amd64.tar.gz", "tool-linux-amd64.zip"},
			ProposedHints: []aimatch.SuggestedHint{{
				Pattern: ".tar.gz",
				Type:    "endswith",
				Match:   "weighted",
				Weight:  1,
			}},
		},
	}})

	for _, expected := range []string{
		"tool — ambiguous asset match",
		"tool-linux-amd64.tar.gz",
		"Confidence: high",
		"pattern: .tar.gz",
		"Validation: resolves uniquely",
	} {
		if !strings.Contains(output, expected) {
			t.Errorf("RenderSuggestions() output missing %q:\n%s", expected, output)
		}
	}
}

func TestRenderSuggestionFailure(t *testing.T) {
	t.Parallel()

	output := presentation.RenderSuggestions([]processor.Result{{
		Tool:            &tool.Tool{Name: "tool"},
		Error:           &match.SelectionError{Cause: match.ErrNoQualified},
		SuggestionError: errors.New("model unavailable"),
	}})

	if !strings.Contains(output, "Suggestion unavailable: model unavailable") {
		t.Errorf("RenderSuggestions() output = %q", output)
	}
}
