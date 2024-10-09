package tools

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

// Aliases is a custom type for handling alias names.
type Aliases []string

// UnmarshalYAML handles custom unmarshalling for Aliases.
// Aliases can hold either a single string or a slice of strings.
func (a *Aliases) UnmarshalYAML(value *yaml.Node) error {
	var single string
	var multiple []string

	// Try unmarshalling into a single string
	if err := value.Decode(&single); err == nil {
		*a = Aliases{single}
		return nil
	}

	// Try unmarshalling into a slice of strings
	if err := value.Decode(&multiple); err == nil {
		*a = Aliases(multiple)
		return nil
	}

	return fmt.Errorf("failed to unmarshal Aliases")
}
