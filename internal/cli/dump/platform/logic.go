package platform

import (
	"fmt"

	"github.com/idelchi/godyl/internal/detect"
	"github.com/idelchi/godyl/internal/iutils"
)

func run() error {
	platform := &detect.Platform{}
	if err := platform.Detect(); err != nil {
		return fmt.Errorf("detecting platform: %w", err)
	}

	iutils.Print(iutils.YAML, platform)

	return nil
}
