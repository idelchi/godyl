// Package inherit provides types for handling tool inheritance configurations.
package inherit

import "github.com/idelchi/godyl/pkg/unmarshal"

// Inherit represents inheritance specifications for tools.
type Inherit = unmarshal.SingleOrSliceType[string]
