package utils

import (
	"bytes"
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/ast"
)

// NodeToBytes converts an ast.Node to bytes.
func NodeToBytes(node ast.Node) ([]byte, error) {
	var buf bytes.Buffer

	// Encode the node to the buffer
	enc := yaml.NewEncoder(&buf)
	if err := enc.Encode(node); err != nil {
		return nil, fmt.Errorf("encoding YAML node: %w", err)
	}

	if err := enc.Close(); err != nil {
		return nil, fmt.Errorf("closing YAML encoder: %w", err)
	}

	return buf.Bytes(), nil
}

// BytesToScalar attempts to unmarshal YAML bytes as a scalar value.
// Returns an error if the bytes cannot be unmarshaled as a scalar.
func BytesToScalar(data []byte, out *string) error {
	return yaml.Unmarshal(data, out)
}
