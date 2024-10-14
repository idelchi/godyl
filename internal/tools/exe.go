package tools

import (
	"github.com/idelchi/godyl/pkg/unmarshal"
	"gopkg.in/yaml.v3"
)

// Exe represents the configuration for the executable name
// and the patterns to search for.
type Exe struct {
	// Name to store the binary as in the output folder.
	Name string
	// Patterns to search for the binary in the downloaded folder.
	Patterns unmarshal.SingleOrSlice[string]
}

// UnmarshalYAML implements custom unmarshaling for `Exe`,
// allowing to set only the name directly or the full struct.
func (e *Exe) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind == yaml.ScalarNode {
		e.Name = value.Value

		return nil
	}

	if value.Kind == yaml.ScalarNode {
		name := value.Value
		value.Kind = yaml.MappingNode
		value.Content = []*yaml.Node{
			{Kind: yaml.ScalarNode, Value: "name"},
			{Kind: yaml.ScalarNode, Value: name},
		}
	}

	type raw Exe
	return unmarshal.DecodeWithOptionalKnownFields(value, (*raw)(e), true, e)
}
