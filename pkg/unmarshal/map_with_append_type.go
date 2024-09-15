package unmarshal

import (
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/ast"
)

type MapWithAppendType[T comparable, V any] map[T]V

func (m *MapWithAppendType[T, V]) UnmarshalYAML(node ast.Node) error {
	return MapWithAppend(m, node)
}

func MapWithAppend[M ~map[T]V, T comparable, V any](m *M, node ast.Node) error {
	// first assignment
	if *m == nil {
		var tmp map[T]V
		if err := Decode(node, &tmp); err != nil {
			return err
		}

		*m = M(tmp)

		return nil
	}

	// merge into existing value
	existing, err := yaml.ValueToNode(*m)
	if err != nil {
		return fmt.Errorf("value→node: %w", err)
	}

	if err := ast.Merge(existing, node); err != nil {
		return fmt.Errorf("merge: %w", err)
	}

	var merged map[T]V
	if err := Decode(existing, &merged); err != nil {
		return fmt.Errorf("decode merged: %w", err)
	}
	*m = M(merged)

	return nil
}
