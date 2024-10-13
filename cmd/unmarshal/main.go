package main

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

// UnmarshalSingleOrSlice is a generic function to unmarshal YAML into either a single value or a slice of values
func UnmarshalSingleOrSlice[T any](value *yaml.Node) ([]T, error) {
	var result []T

	// Try unmarshalling into a single value
	var single T
	if err := value.Decode(&single); err == nil {
		result = append(result, single)
		return result, nil
	}

	// Try unmarshalling into a slice of T
	var multiple []T
	if err := value.Decode(&multiple); err == nil {
		result = multiple
		return result, nil
	}

	return result, fmt.Errorf("failed to unmarshal: expected single value or slice")
}

// Custom types
type StringList []string

func (s *StringList) UnmarshalYAML(value *yaml.Node) error {
	result, err := UnmarshalSingleOrSlice[string](value)
	if err != nil {
		return err
	}
	*s = result
	return nil
}

type IntList []int

func (i *IntList) UnmarshalYAML(value *yaml.Node) error {
	result, err := UnmarshalSingleOrSlice[int](value)
	if err != nil {
		return err
	}
	*i = result
	return nil
}

type Person struct {
	Name string `yaml:"name"`
	Age  int    `yaml:"age"`
}

type PersonList []Person

func (p *PersonList) UnmarshalYAML(value *yaml.Node) error {
	result, err := UnmarshalSingleOrSlice[Person](value)
	if err != nil {
		return err
	}
	*p = result
	return nil
}

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
