// Package exe provides functionality for configuring tool executables.
package exe

import (
	"github.com/goccy/go-yaml/ast"
	"github.com/idelchi/godyl/pkg/unmarshal"
)

// Exe defines the configuration for locating and naming tool executables.
type Exe struct {
	// Name is the final name for the executable in the output directory.
	// This will be the name used to run the tool after installation.
	Name string `single:"true"`

	// Patterns contains regex patterns for finding the executable.
	// Used to locate the correct binary file within downloaded content.
	// Can be a single pattern or multiple patterns in order of preference.
	Patterns *Patterns
}

type Patterns = unmarshal.SingleOrSliceType[string]

// UnmarshalYAML implements custom YAML unmarshaling for Exe configuration.
// Supports both scalar values (treated as executable name) and map values.
func (e *Exe) UnmarshalYAML(node ast.Node) error {
	type raw Exe

	return unmarshal.SingleStringOrStruct(node, (*raw)(e))
}
