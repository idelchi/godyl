package root

import (
	"fmt"

	"github.com/idelchi/godyl/internal/ierrors"
)

// Validate checks the configuration for errors.
// At the moment, it only checks if the defaults file exists.
func (r *Root) Validate() error {
	if r.IsSet("defaults") && !r.Defaults.Exists() {
		return fmt.Errorf("%w: %q does not exist", ierrors.ErrUsage, r.Defaults)
	}

	return nil
}
