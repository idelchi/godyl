package presentation

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/idelchi/godyl/internal/aimatch"
	"github.com/idelchi/godyl/internal/match"
	"github.com/idelchi/godyl/internal/processor"
	"github.com/idelchi/godyl/pkg/pretty"
)

type suggestedConfig struct {
	Hints []aimatch.SuggestedHint `yaml:"hints"`
}

// RenderSuggestions formats verified recommendations and suggestion failures.
func RenderSuggestions(results []processor.Result) string {
	results = slices.Clone(results)
	slices.SortFunc(results, func(a, b processor.Result) int {
		return strings.Compare(a.Tool.Name, b.Tool.Name)
	})

	var output strings.Builder

	for _, result := range results {
		if result.Suggestion == nil && result.SuggestionError == nil {
			continue
		}

		if output.Len() > 0 {
			output.WriteString("\n\n")
		}

		fmt.Fprintf(&output, "%s — %s\n", result.Tool.Name, matchingFailureName(result.Error))

		if result.SuggestionError != nil {
			fmt.Fprintf(&output, "  Suggestion unavailable: %v", result.SuggestionError)

			continue
		}

		suggestion := result.Suggestion

		output.WriteString("  Candidates considered:\n")

		for _, candidate := range suggestion.Candidates {
			fmt.Fprintf(&output, "    - %s\n", candidate)
		}

		fmt.Fprintf(&output, "  Likely asset: %s\n", suggestion.SelectedAsset)
		fmt.Fprintf(&output, "  Confidence: %s\n", suggestion.Confidence)
		fmt.Fprintf(&output, "  Reason: %s\n", suggestion.Reason)
		output.WriteString("  Verified replacement hints:\n")

		yaml := strings.TrimSpace(pretty.YAML(suggestedConfig{Hints: suggestion.ProposedHints}))
		fmt.Fprintf(&output, "%s\n", indent(yaml, "    "))
		output.WriteString("  Validation: resolves uniquely to the likely asset")
	}

	return output.String()
}

func matchingFailureName(err error) string {
	switch {
	case errors.Is(err, match.ErrAmbiguous):
		return "ambiguous asset match"
	case errors.Is(err, match.ErrNoQualified):
		return "no qualified asset match"
	default:
		return "asset matching failure"
	}
}

func indent(value, prefix string) string {
	return prefix + strings.ReplaceAll(value, "\n", "\n"+prefix)
}
