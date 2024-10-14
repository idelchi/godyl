package unmarshal

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/idelchi/godyl/pkg/pretty"
	"gopkg.in/yaml.v3"
)

type SingleOrSlice[T any] []T

func (ss *SingleOrSlice[T]) UnmarshalYAML(value *yaml.Node) error {
	fmt.Println("UnmarshalYAML")
	pretty.PrintJSON(value)

	result, err := UnmarshalSingleOrSlice[T](value, true)
	if err != nil {
		return err
	}
	*ss = result
	return nil
}

func UnmarshalSingleOrSlice[T any](node *yaml.Node, useKnownFields bool) ([]T, error) {
	// If it's a scalar or mapping node, wrap it in a sequence node
	if node.Kind == yaml.ScalarNode || node.Kind == yaml.MappingNode {
		node = &yaml.Node{
			Kind:    yaml.SequenceNode,
			Content: []*yaml.Node{node},
		}
	}

	var result []T

	// Use UnmarshalWithKnownFields instead of node.Decode
	if err := DecodeWithOptionalKnownFields(node, &result, useKnownFields, result); err != nil {
		return nil, err
	}

	return result, nil
}

// unmarshalYAMLWithOptions is a helper function that unmarshals YAML with the option to use KnownFields.
func DecodeWithOptionalKnownFields(value *yaml.Node, out any, useKnownFields bool, input any) error {
	// Re-encode the yaml.Node to bytes
	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)
	if err := enc.Encode(value); err != nil {
		return err
	}
	enc.Close()

	// Decode from the buffer
	decoder := yaml.NewDecoder(&buf)
	if useKnownFields {
		decoder.KnownFields(true)
	}

	// Decode into the provided interface
	err := decoder.Decode(out)

	return yamlTypeErrorConversion(err, input)
}

func yamlTypeErrorConversion(err error, input any) error {
	if err == nil {
		return nil
	}

	if _, ok := err.(*yaml.TypeError); !ok {
		return err
	}

	if !strings.Contains(err.Error(), " not found in type") {
		return err
	}

	typeErr := err.(*yaml.TypeError)

	for i, err := range typeErr.Errors {
		parts := strings.SplitN(err, " not found in type", 2)
		if len(parts) != 2 {
			continue
		}

		typeErr.Errors[i] = fmt.Sprintf("%s not found in type %T", parts[0], input)
	}

	return err
}
