package unmarshal

import (
	"errors"

	"github.com/fatih/structs"
	"github.com/goccy/go-yaml/ast"
)

// SingleStringOrStruct handles either a string value or a full struct decoding.
func SingleStringOrStruct[T any](node ast.Node, out *T) error {
	// Handle primitive value types
	if isPrimitiveNode(node) {
		field, found := findTaggedSingleField(out)
		if !found {
			return errors.New("no field with `single:\"true\"` tag found")
		}

		// for scalars only
		if sn, ok := node.(*ast.StringNode); ok {
			return field.Set(sn.Value) // .Value is the de‑quoted scalar content
		}

		return field.Set(node.String())
	}

	// Not a primitive type, decode the whole structure
	return Decode(node, out)
}

// Helper to check if node is a primitive value type.
func isPrimitiveNode(node ast.Node) bool {
	switch node.(type) {
	case *ast.StringNode, *ast.IntegerNode, *ast.FloatNode, *ast.BoolNode:
		return true
	default:
		return false
	}
}

// findTaggedSingleField finds the field tagged with single:"true".
func findTaggedSingleField[T any](out *T) (*structs.Field, bool) {
	s := structs.New(out)

	for _, field := range s.Fields() {
		if field.Tag("single") == "true" {
			return field, true
		}
	}

	return nil, false
}
