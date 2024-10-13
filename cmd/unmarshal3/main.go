package main

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

func UnmarshalSingleOrSlice[T any](value *yaml.Node) ([]T, error) {
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
	// Test cases
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

	fmt.Println("Testing StringList:")
	for _, test := range stringTests {
		var sl StringList
		err := yaml.Unmarshal([]byte(test), &sl)
		fmt.Printf("Input: %s\nResult: %v\nError: %v\n\n", test, sl, err)
	}

	fmt.Println("Testing IntList:")
	for _, test := range intTests {
		var il IntList
		err := yaml.Unmarshal([]byte(test), &il)
		fmt.Printf("Input: %s\nResult: %v\nError: %v\n\n", test, il, err)
	}

	fmt.Println("Testing PersonList:")
	for _, test := range personTests {
		var pl PersonList
		err := yaml.Unmarshal([]byte(test), &pl)
		fmt.Printf("Input: %s\nResult: %v\nError: %v\n\n", test, pl, err)
	}
}
