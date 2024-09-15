package match

import (
	"github.com/idelchi/godyl/internal/detect"
	"github.com/idelchi/godyl/internal/tools/hints"
)

// Requirements represents the criteria an asset must meet.
// It includes platform compatibility and a list of hints for name matching.
type Requirements struct {
	// Hints are used to match asset names.
	Hints hints.Hints
	// Platform represents the target platform for the asset.
	Platform detect.Platform
	// Checksum is a pattern to match checksum assets.
	Checksum string
}
