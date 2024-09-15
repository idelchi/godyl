package unmarshal

import (
	"fmt"

	"github.com/goccy/go-yaml/ast"
)

// SingleOrSliceType represents a flexible YAML unmarshaling target.
// Implements custom unmarshaling to handle both single values and
// slices of values in YAML input. Works with any comparable type T.
type SingleOrSliceType[T any] []T

// UnmarshalYAML implements yaml.NodeUnmarshaler for SingleOrSliceType.
// Calls SingleOrSlice to do the actual decoding.
func (ss *SingleOrSliceType[T]) UnmarshalYAML(node ast.Node) (err error) {
	*ss, err = SingleOrSlice[T](node)
	if err != nil {
		return err
	}

	return nil
}

// SingleOrSlice unmarshals YAML data into a slice of type T.
// Handles both single values and sequences by wrapping scalars/mappings.
func SingleOrSlice[T any](node ast.Node) ([]T, error) {
	var out []T

	switch n := node.(type) {
	case *ast.SequenceNode:
		for _, elem := range n.Values {
			var v T

			if err := Decode(elem, &v); err != nil {
				return nil, fmt.Errorf("unmarshal sequence element: %w", err)
			}

			out = append(out, v)
		}
	default:
		var v T

		if err := Decode(node, &v); err != nil {
			return nil, fmt.Errorf("unmarshal single value: %w", err)
		}

		out = append(out, v)
	}

	return out, nil
}
