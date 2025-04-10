package tools

import (
	"fmt"
	"strconv"

	"github.com/idelchi/godyl/pkg/unmarshal"

	"gopkg.in/yaml.v3"
)

// Skip represents a list of conditions under which certain operations (e.g., downloading, testing) should be skipped.
type Skip []Condition

// Condition defines a condition and an optional reason for skipping an operation.
type Condition struct {
	// Condition is a string that represents a boolean expression (e.g., "true" or "false") that determines whether the
	// operation should be skipped.
	Condition string
	// Reason provides an optional explanation for why the operation is being skipped.
	Reason string

	condition bool
}

func (c Condition) True() bool {
	return c.condition
}

func (c *Condition) Parse() (err error) {
	// Parse the condition string into a boolean value.
	c.condition, err = strconv.ParseBool(c.Condition)
	if err != nil {
		return fmt.Errorf("parsing condition %q: %w", c.Condition, err)
	}

	return
}

func (s Skip) Has() bool {
	return len(s) > 0
}

// Any checks if any condition in the Skip list evaluates to true.
func (s *Skip) Evaluate() (Skip, error) {
	skip := Skip{}

	for _, condition := range *s {
		if err := condition.Parse(); err != nil {
			return nil, err
		}

		if condition.True() {
			skip = append(skip, condition)
		}
	}

	return skip, nil
}

// UnmarshalYAML implements custom unmarshaling for Skip,
// allowing the Skip field to be either a single condition or a list of conditions.
func (s *Skip) UnmarshalYAML(value *yaml.Node) error {
	// If the YAML value is a scalar (e.g., just a single condition), handle it directly.
	if value.Kind == yaml.ScalarNode {
		*s = []Condition{{Condition: value.Value}}

		return nil
	}

	// Otherwise, treat it as a list of conditions and unmarshal accordingly.
	result, err := unmarshal.SingleOrSlice[Condition](value, true)
	if err != nil {
		return err
	}

	*s = result

	return nil
}
