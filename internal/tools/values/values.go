// Package values provides functionality for managing tool configuration values.
package values

import (
	"github.com/idelchi/godyl/pkg/unmarshal"
)

// Values represents a map of string keys to any values.
type Values = unmarshal.MapWithAppendType[string, any]
