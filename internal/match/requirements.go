package match

import (
	"github.com/idelchi/godyl/internal/detect"
)

type Requirements struct {
	Platform detect.Platform
	Hints    []Hint
}


