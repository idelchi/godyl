// Package aliases provides functionality for managing tool alias names.
package aliases

import (
	"github.com/idelchi/godyl/pkg/unmarshal"
)

// Aliases represents a tool's alias names.
type Aliases = unmarshal.SingleOrSliceType[string]
