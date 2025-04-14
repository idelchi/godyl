// Package unmarshal provides utilities for unmarshalling YAML data that can
// represent either a single item or a slice of items. It includes a generic
// type `SingleOrSlice` to handle this pattern, allowing flexible unmarshalling
// from YAML input.
//
// The package also provides functions to decode YAML nodes while optionally
// enforcing that only known fields are present, improving error handling
// in case of unexpected fields in the YAML input.
package unmarshal

import (
	"bytes"
	"errors"
	"fmt"
	"slices"
	"strings"

	"gopkg.in/yaml.v3"
)

// SingleOrSliceType represents a flexible YAML unmarshaling target.
// Implements custom unmarshaling to handle both single values and
// slices of values in YAML input. Works with any comparable type T.
type SingleOrSliceType[T comparable] []T

// UnmarshalYAML implements yaml.Unmarshaler for SingleOrSliceType.
// Handles both scalar values and sequences in YAML, converting them
// to a slice of type T. Single values are wrapped in a slice.
func (ss *SingleOrSliceType[T]) UnmarshalYAML(value *yaml.Node) error {
	result, err := SingleOrSlice[T](value, true)
	if err != nil {
		return err
	}

	*ss = result

	return nil
}

// Compacted returns a new slice with duplicate elements removed.
// Preserves order while eliminating repeated values.
func (ss *SingleOrSliceType[T]) Compacted() SingleOrSliceType[T] {
	return SingleOrSliceType[T](slices.Compact([]T(*ss)))
}

// SingleOrSlice unmarshals YAML data into a slice of type T.
// Handles both single values and sequences by automatically wrapping
// scalar or mapping nodes in a sequence. The useKnownFields parameter
// controls validation of unknown YAML fields.
func SingleOrSlice[T any](node *yaml.Node, useKnownFields bool) ([]T, error) {
	// If it's a scalar or mapping node, wrap it in a sequence node
	if node.Kind == yaml.ScalarNode || node.Kind == yaml.MappingNode {
		node = &yaml.Node{
			Kind:    yaml.SequenceNode,
			Content: []*yaml.Node{node},
		}
	}

	var result []T

	// Use DecodeWithOptionalKnownFields to decode the node into the result slice.
	if err := DecodeWithOptionalKnownFields(node, &result, useKnownFields, "any"); err != nil {
		return nil, err
	}

	return result, nil
}

// DecodeWithOptionalKnownFields decodes YAML with field validation.
// Decodes a YAML node into the provided output interface. When
// useKnownFields is true, enforces that all YAML fields exist in
// the target type. The input parameter provides context for errors.
func DecodeWithOptionalKnownFields(value *yaml.Node, out any, useKnownFields bool, input string) error {
	// Re-encode the yaml.Node to bytes
	var buf bytes.Buffer

	enc := yaml.NewEncoder(&buf)
	if err := enc.Encode(value); err != nil {
		return fmt.Errorf("encoding YAML node: %w", err)
	}

	if err := enc.Close(); err != nil {
		return fmt.Errorf("closing YAML encoder: %w", err)
	}

	// Decode from the buffer
	decoder := yaml.NewDecoder(&buf)
	if useKnownFields {
		decoder.KnownFields(true)
	}

	// Decode into the provided interface
	err := decoder.Decode(out)

	return yamlTypeErrorConversion(err, input)
}

// yamlTypeErrorConversion enhances YAML type error messages.
// Improves yaml.TypeError messages by adding input type information
// to "not found in type" errors. Makes validation errors more
// descriptive and easier to understand.
func yamlTypeErrorConversion(err error, input string) error {
	if err == nil {
		return nil
	}

	var typeErr *yaml.TypeError
	if !errors.As(err, &typeErr) {
		return err
	}

	if !strings.Contains(err.Error(), " not found in type") {
		return err
	}

	const expectedParts = 2

	for idx, errMsg := range typeErr.Errors {
		parts := strings.SplitN(errMsg, " not found in type", expectedParts)
		if len(parts) != expectedParts {
			continue
		}

		typeErr.Errors[idx] = fmt.Sprintf("%s not found in type %q", parts[0], input)
	}

	return err
}
