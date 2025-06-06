package root

import (
	"fmt"

	"github.com/idelchi/godyl/internal/ierrors"
)

// Validate checks the root configuration for correctness.
func (r *Root) Validate() error {
	if r.IsSet("defaults") && !r.Defaults.Exists() {
		return fmt.Errorf("%w: %q does not exist", ierrors.ErrUsage, r.Defaults)
	}

	return nil
}
