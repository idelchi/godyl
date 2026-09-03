package platform

import (
	"fmt"

	"github.com/idelchi/godyl/internal/cli/core"
	"github.com/idelchi/godyl/internal/detect"
	"github.com/idelchi/godyl/internal/iutils"
)

// run detects and renders the current operating system, architecture, library, and distribution.
func run(_ core.Input) error {
	platform := &detect.Platform{}
	if err := platform.Detect(); err != nil {
		return fmt.Errorf("detecting platform: %w", err)
	}

	iutils.Print(iutils.YAML, platform)

	return nil
}
