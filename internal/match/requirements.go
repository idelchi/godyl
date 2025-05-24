package match

import (
	"github.com/idelchi/godyl/internal/detect"
	"github.com/idelchi/godyl/internal/tools/hints"
)

// Requirements represents the criteria an asset must meet.
// It includes platform compatibility and a list of hints for name matching.
type Requirements struct {
	Hints    hints.Hints
	Platform detect.Platform
}
