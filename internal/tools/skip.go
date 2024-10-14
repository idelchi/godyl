package tools

import (
	"fmt"

	"github.com/idelchi/godyl/pkg/unmarshal"
	"gopkg.in/yaml.v3"
)

type Skip struct {
	Template string
	Message  string
	Skip     bool
}

func (s *Skip) UnmarshalYAML(value *yaml.Node) error {
	// If it's a scalar node, handle it directly
	if value.Kind == yaml.ScalarNode {
		switch value.Tag {
		case "!!bool":
			// If it's a boolean, assign to Skip.Skip
			return value.Decode(&s.Skip)
		case "!!str":
			// If it's a string, assign to Skip.Template
			return value.Decode(&s.Template)
		default:
			return fmt.Errorf("unexpected scalar type for Skip: %s", value.Tag)
		}
	}

	type rawSkip Skip
	return unmarshal.DecodeWithOptionalKnownFields(value, (*rawSkip)(s), true, s)
}
