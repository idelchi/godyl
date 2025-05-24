// Package skip provides functionality for managing tool skip conditions.
package skip

import (
	"fmt"

	"github.com/goccy/go-yaml/ast"

	"github.com/idelchi/godyl/pkg/unmarshal"
)

// Skip represents a list of conditions under which certain operations (e.g., downloading, testing) should be skipped.
type Skip unmarshal.SingleOrSliceType[Condition]

// Condition defines a condition and an optional reason for skipping an operation.
type Condition struct {
	Reason    string
	Condition unmarshal.Templatable[bool]
}

func (c Condition) True() bool {
	return c.Condition.Value
}

func (c *Condition) Parse() error {
	// Parse the condition string into a boolean value.
	err := c.Condition.Parse()
	if err != nil {
		return fmt.Errorf("parsing condition: %w", err)
	}

	return nil
}

func (s Skip) Has() bool {
	return len(s) > 0
}

// Evaluate checks if any condition in the Skip list evaluates to true.
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
func (s *Skip) UnmarshalYAML(node ast.Node) (err error) {
	*s, err = unmarshal.SingleOrSlice[Condition](node)
	if err != nil {
		return fmt.Errorf("unmarshaling skip: %w", err)
	}

	return nil
}
