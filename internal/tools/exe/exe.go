// Package exe provides functionality for configuring tool executables.
package exe

import (
	"github.com/goccy/go-yaml/ast"

	"github.com/idelchi/godyl/pkg/unmarshal"
)

// Exe defines the configuration for locating and naming tool executables.
type Exe struct {
	Patterns *Patterns
	Name     string `single:"true"`
}

// Patterns represents executable pattern matching rules.
type Patterns = unmarshal.SingleOrSliceType[string]

// UnmarshalYAML implements custom YAML unmarshaling for Exe configuration.
// Supports both scalar values (treated as executable name) and map values.
func (e *Exe) UnmarshalYAML(node ast.Node) error {
	type raw Exe

	return unmarshal.SingleStringOrStruct(node, (*raw)(e))
}
