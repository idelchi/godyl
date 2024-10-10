package tools

import (
	"bytes"

	"gopkg.in/yaml.v3"
)

// Aliases is a custom type for handling alias names.
type Exe struct {
	Name     string
	Patterns []string
}

func (e *Exe) UnmarshalYAML(value *yaml.Node) error {
	// If it's a scalar (e.g., just the name), handle it directly
	if value.Kind == yaml.ScalarNode {
		e.Name = value.Value
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
