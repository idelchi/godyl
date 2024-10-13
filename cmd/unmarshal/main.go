package main

import (
	"bytes"
	"fmt"

	"gopkg.in/yaml.v3"
)

func UnmarshalSingleOrSlice[T any](node *yaml.Node) ([]T, error) {
	var result []T

	switch node.Kind {
	case yaml.ScalarNode:
		var single T
		if err := node.Decode(&single); err != nil {
			return nil, fmt.Errorf("failed to decode scalar: %w", err)
		}
		result = append(result, single)

	case yaml.SequenceNode:
		if err := node.Decode(&result); err != nil {
			return nil, fmt.Errorf("failed to decode sequence: %w", err)
		}

	case yaml.MappingNode:
		var single T
		// Re-encode the yaml.Node to bytes to leverage yaml.NewDecoder
		var buf bytes.Buffer
		enc := yaml.NewEncoder(&buf)
		if err := enc.Encode(node); err != nil {
			return nil, fmt.Errorf("failed to re-encode mapping: %w", err)
		}
		enc.Close()

		// Now decode from the buffer with KnownFields enabled
		decoder := yaml.NewDecoder(&buf)
		decoder.KnownFields(true)

		if err := decoder.Decode(&single); err != nil {
			return nil, fmt.Errorf("failed to decode mapping with KnownFields: %w", err)
		}
		result = append(result, single)

	default:
		return nil, fmt.Errorf("unsupported YAML node kind: %v", node.Kind)
	}

	return result, nil
}

func UnmarshalSingleOrSlice4[T any](node *yaml.Node) ([]T, error) {
	var value any

	switch node.Kind {
	case yaml.ScalarNode:
		value = new(T)
	case yaml.SequenceNode:
		value = new([]T)
	case yaml.MappingNode:
		value = new(T)
	default:
		return nil, fmt.Errorf("unsupported YAML node kind: %v", node.Kind)
	}

	// Create a decoder with KnownFields enabled
	decoder := yaml.NewDecoder(nil)
	decoder.KnownFields(true)

	// Decode the node
	if err := decoder.Decode(value); err != nil {
		return nil, fmt.Errorf("failed to decode: %w", err)
	}

	// Convert the result to []T
	switch v := value.(type) {
	case *T:
		return []T{*v}, nil
	case *[]T:
		return *v, nil
	default:
		return nil, fmt.Errorf("unexpected type after decoding")
	}
}

func UnmarshalSingleOrSlice3[T any](value *yaml.Node) ([]T, error) {
	var result []T

	switch value.Kind {
	case yaml.ScalarNode:
		var single T
		if err := value.Decode(&single); err != nil {
			return nil, fmt.Errorf("failed to unmarshal scalar: %w", err)
		}
		result = append(result, single)
	case yaml.SequenceNode:
		var multiple []T
		if err := value.Decode(&multiple); err != nil {
			return nil, fmt.Errorf("failed to unmarshal sequence: %w", err)
		}
		result = multiple
	case yaml.MappingNode:
		var single T
		if err := value.Decode(&single); err != nil {
			return nil, fmt.Errorf("failed to unmarshal mapping: %w", err)
		}
		result = append(result, single)
	default:
		return nil, fmt.Errorf("unsupported YAML node kind: %v", value.Kind)
	}

	return result, nil
}

// GenericList is a generic wrapper type that implements UnmarshalYAML
type GenericList[T any] []T

func (g *GenericList[T]) UnmarshalYAML(value *yaml.Node) error {
	result, err := UnmarshalSingleOrSlice[T](value)
	if err != nil {
		return err
	}
	*g = result
	return nil
}

type Person struct {
	Name string `yaml:"name"`
	Age  int    `yaml:"age"`
}

// Type aliases using the generic wrapper
type (
	StringList = GenericList[string]
	IntList    = GenericList[int]
	PersonList = GenericList[Person]
)

func main() {
	// Success test cases (as before)
	stringTests := []string{
		"hello",
		"[hello, world]",
	}
	intTests := []string{
		"42",
		"[1, 2, 3]",
	}
	personTests := []string{
		"{name: Alice, age: 30}",
		"[{name: Alice, age: 30}, {name: Bob, age: 25}]",
	}

	// Failure test cases
	failureTests := []struct {
		name     string
		yaml     string
		testType string
	}{
		{"Invalid YAML", "[:invalid", "StringList"},
		{"Wrong type for StringList", "42", "StringList"},
		{"Wrong type for IntList", "not a number", "IntList"},
		{"Invalid Person structure", "Alice", "PersonList"},
		{"Mixed types in IntList", "[1, two, 3]", "IntList"},
		{"Invalid nested structure", "{persons: [{name: Alice}]}", "PersonList"},
	}

	// Run success tests
	runTests := func(tests []string, list interface{}) {
		for _, test := range tests {
			err := yaml.Unmarshal([]byte(test), list)
			fmt.Printf("Input: %s\nResult: %v\nError: %v\n\n", test, list, err)
		}
	}

	fmt.Println("Testing StringList (Success cases):")
	runTests(stringTests, &StringList{})

	fmt.Println("Testing IntList (Success cases):")
	runTests(intTests, &IntList{})

	fmt.Println("Testing PersonList (Success cases):")
	runTests(personTests, &PersonList{})

	// Run failure tests
	fmt.Println("Testing Failure Cases:")
	for _, test := range failureTests {
		fmt.Printf("Test: %s\nInput: %s\n", test.name, test.yaml)
		var err error
		switch test.testType {
		case "StringList":
			var sl StringList
			err = yaml.Unmarshal([]byte(test.yaml), &sl)
		case "IntList":
			var il IntList
			err = yaml.Unmarshal([]byte(test.yaml), &il)
		case "PersonList":
			var pl PersonList
			err = yaml.Unmarshal([]byte(test.yaml), &pl)

			fmt.Println(pl)
		}
		if err != nil {
			fmt.Printf("Error (as expected): %v\n\n", err)
		} else {
			fmt.Printf("Unexpected success: no error occurred\n\n")
		}
	}
}
