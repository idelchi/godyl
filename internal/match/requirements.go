package match

import (
	"github.com/idelchi/godyl/internal/detect"
)

// Requirements represents the criteria an asset must meet.
// It includes platform compatibility and a list of hints for name matching.
type Requirements struct {
	Platform detect.Platform // Platform specifies the required OS, architecture, and library compatibility.
	Hints    []Hint          // Hints contains patterns used to match the asset's name.
}
