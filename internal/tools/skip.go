package tools

import (
	"strconv"

	"github.com/idelchi/godyl/pkg/unmarshal"
	"gopkg.in/yaml.v3"
)

type Conditions []string

func (c *Conditions) UnmarshalYAML(value *yaml.Node) error {
	result, err := unmarshal.UnmarshalSingleOrSlice[string](value, true)
	if err != nil {
		return err
	}
	*c = result
	return nil
}

func (c Conditions) IsSkipped() (bool, error) {
	for _, condition := range c {
		if val, err := strconv.ParseBool(condition); err != nil {
			return false, err
		} else {
			if val {
				return true, nil
			}
		}
	}

	return false, nil
}

type Skip struct {
	Conditions Conditions
	Message    string

	skip bool
}

func (s *Skip) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind == yaml.ScalarNode {
		s.Conditions = []string{value.Value}

		return nil
	}

	// if value.Kind == yaml.ScalarNode {
	// 	switch value.Tag {
	// 	case "!!bool":
	// 		// If it's a boolean, assign to Skip.Skip
	// 		return value.Decode(&s.Skip)
	// 	case "!!str":
	// 		// If it's a string, assign to Skip.Template
	// 		return value.Decode(&s.Template)
	// 	default:
	// 		return fmt.Errorf("unexpected scalar type for Skip: %s", value.Tag)
	// 	}
	// }

	type rawSkip Skip
	return unmarshal.DecodeWithOptionalKnownFields(value, (*rawSkip)(s), true, s)
}
