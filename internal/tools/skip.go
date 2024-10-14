package tools

import (
	"fmt"
	"strconv"

	"github.com/idelchi/godyl/pkg/pretty"
	"github.com/idelchi/godyl/pkg/unmarshal"
	"gopkg.in/yaml.v3"
)

type Skip []Condition

type Condition struct {
	Condition string
	Reason    string
}

func (s *Skip) UnmarshalYAML(value *yaml.Node) error {
	// Manually set value to Skip[0].Condition if it's a scalar

	fmt.Printf("UnmarshalYAML: %v\n", pretty.YAML(value))

	result, err := unmarshal.UnmarshalSingleOrSlice[Condition](value, true)
	if err != nil {
		return err
	}
	*s = result
	return nil
}

func (s Skip) IsSkipped() (bool, string, error) {
	for _, condition := range s {
		if val, err := strconv.ParseBool(condition.Condition); err != nil {
			return false, condition.Reason, err
		} else {
			if val {
				return true, condition.Reason, nil
			}
		}
	}

	return false, "", nil
}

// func (s *Skip) UnmarshalYAML(value *yaml.Node) error {
// 	if value.Kind == yaml.ScalarNode {
// 		s.Conditions = []string{value.Value}

// 		return nil
// 	}

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

// 	type rawSkip Skip
// 	return unmarshal.DecodeWithOptionalKnownFields(value, (*rawSkip)(s), true, s)
// }
