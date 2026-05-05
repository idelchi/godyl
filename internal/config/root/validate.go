package root

import (
	"fmt"

	"github.com/idelchi/godyl/internal/ierrors"
)

// Validate checks the root configuration for correctness.
func (c *Config) Validate() error {
	if c.IsSet("defaults") && !c.Defaults.Exists() {
		return fmt.Errorf("%w: %q does not exist", ierrors.ErrUsage, c.Defaults)
	}

	if c.IsSet("go") && !c.Go.Expanded().Exists() {
		return fmt.Errorf("%w: %q does not exist", ierrors.ErrUsage, c.Go)
	}

	return nil
}
