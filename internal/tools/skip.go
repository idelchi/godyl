package tools

import (
	"strconv"

	"github.com/idelchi/godyl/pkg/unmarshal"

	"gopkg.in/yaml.v3"
)

type Skip []Condition

type Condition struct {
	Condition string
	Reason    string
}

func (s Skip) True() (bool, string, error) {
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

// UnmarshalYAML implements custom unmarshaling for `Exe`,
// allowing to set only the name directly or the full struct.
func (s *Skip) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind == yaml.ScalarNode {
		*s = []Condition{{Condition: value.Value}}

		return nil
	}

	result, err := unmarshal.UnmarshalSingleOrSlice[Condition](value, true)
	if err != nil {
		return err
	}
	*s = result
	return nil
}
