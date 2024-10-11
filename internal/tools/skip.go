package tools

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type Skip struct {
	Template     string `json:"-"`
	Message      string
	SkipTemplate string `json:"-" yaml:"skip" mapstructure:"skip"`
	Skip         bool   `yaml:"-" mapstructure:"-"`
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

	// If it's a mapping node, handle it as normal unmarshalling
	if value.Kind == yaml.MappingNode {
		type rawSkip Skip
		return value.Decode((*rawSkip)(s))
	}

	return fmt.Errorf("unexpected node kind for Skip: %v", value.Kind)
}
