package tools

import (
	"github.com/fatih/structs"

	"github.com/idelchi/godyl/pkg/unmarshal"

	"gopkg.in/yaml.v3"
)

// Exe defines the configuration for locating and naming tool executables.
type Exe struct {
	// Name is the final name for the executable in the output directory.
	// This will be the name used to run the tool after installation.
	Name string

	// Patterns contains regex patterns for finding the executable.
	// Used to locate the correct binary file within downloaded content.
	// Can be a single pattern or multiple patterns in order of preference.
	Patterns unmarshal.SingleOrSliceType[string]
}

// UnmarshalYAML implements custom YAML unmarshaling for Exe configuration.
// Supports both scalar values (treated as executable name) and map values with
// additional configuration for executable discovery.
func (e *Exe) UnmarshalYAML(value *yaml.Node) error {
	// If the YAML value is a scalar (e.g., just the name), handle it directly by setting the Name field.
	if value.Kind == yaml.ScalarNode {
		e.Name = value.Value

		return nil
	}

	// Perform custom unmarshaling with field validation, allowing only known fields.
	type raw Exe

	return unmarshal.DecodeWithOptionalKnownFields(value, (*raw)(e), true, structs.New(e).Name())
}
