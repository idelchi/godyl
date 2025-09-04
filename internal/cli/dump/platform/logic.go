package platform

import (
	"fmt"

	"github.com/idelchi/godyl/internal/cli/core"
	"github.com/idelchi/godyl/internal/detect"
	"github.com/idelchi/godyl/internal/iutils"
)

// run executes the `dump platform` command.
func run(_ core.Input) error {
	platform := &detect.Platform{}
	if err := platform.Detect(); err != nil {
		return fmt.Errorf("detecting platform: %w", err)
	}

	iutils.Print(iutils.YAML, platform)

	return nil
}
