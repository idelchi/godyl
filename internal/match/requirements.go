package match

import (
	"github.com/idelchi/godyl/internal/detect"
)

// Requirements represents the criteria an asset must meet.
// It includes platform compatibility and a list of hints for name matching.
type Requirements struct {
	Hints    []Hint
	Platform detect.Platform
}
