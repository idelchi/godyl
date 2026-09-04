package match

import (
	"context"

	"github.com/idelchi/godyl/internal/detect"
	"github.com/idelchi/godyl/internal/tools/hints"
)

// AssetSelector chooses one exact release asset from an ambiguous set of
// equally ranked deterministic matches.
type AssetSelector interface {
	Select(ctx context.Context, request SelectionRequest) (string, error)
}

// SelectionRequest describes the target and the exact candidates available to
// an optional fallback selector.
type SelectionRequest struct {
	// Target identifies the tool being resolved.
	Target string
	// Platform describes the requested runtime target.
	Platform detect.Platform
	// Candidates contains the permitted return values.
	Candidates []string
}

// Requirements represents the criteria an asset must meet.
// It includes platform compatibility and a list of hints for name matching.
type Requirements struct {
	// Hints are used to match asset names.
	Hints hints.Hints
	// Platform represents the target platform for the asset.
	Platform detect.Platform
	// Checksum is a pattern to match checksum assets.
	Checksum string
	// Target identifies the tool being resolved.
	Target string
	// Selector is an optional tie-breaker for ambiguous deterministic results.
	Selector AssetSelector
}
