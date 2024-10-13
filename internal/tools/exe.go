package tools

import (
	"bytes"

	"gopkg.in/yaml.v3"
)

// Exe represents the configuration for the executable name
// and the patterns to search for.
type Exe struct {
	// Name to store the binary as in the output folder.
	Name string
	// Patterns to search for the binary in the downloaded folder.
	Patterns []string
}

// UnmarshalYAML implements custom unmarshaling for Exe,
// allowing to set only the name directly or the full struct.
// When setting the name on the format:
//
//	exe: name
//
// It will be unmarshaled as:
//
//	exe:
//	  name: name
//
// TODO(Idelchi): Add support for setting the patterns as:
//
//	exe:
//	  patterns: pattern
func (e *Exe) UnmarshalYAML(value *yaml.Node) error {
	// If it's a scalar (e.g., just the name), handle it directly
	if value.Kind == yaml.ScalarNode {
		e.Name = value.Value
		// e.Patterns = []string{value.Value}
		return nil
	}

	// Re-encode the yaml.Node to bytes
	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)
	if err := enc.Encode(value); err != nil {
		return err
	}
	enc.Close()

	// Decode from the buffer with KnownFields enabled
	decoder := yaml.NewDecoder(&buf)
	decoder.KnownFields(true)

	// Decode the Exe
	type rawExe Exe
	if err := decoder.Decode((*rawExe)(e)); err != nil {
		return err
	}

	return nil
}
