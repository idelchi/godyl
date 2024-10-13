package unmarshal

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

// StringOrSlice is a type constraint for types that are slices of strings
type StringOrSlice interface {
	~[]string
}

// UnmarshalStringOrSlice is a generic function to unmarshal YAML into either a single string or a slice of strings
func UnmarshalStringOrSlice[T StringOrSlice](value *yaml.Node) (T, error) {
	var result T

	// Try unmarshalling into a single string
	var single string
	if err := value.Decode(&single); err == nil {
		result = append(result, single)
		return result, nil
	}

	// Try unmarshalling into a slice of strings
	var multiple []string
	if err := value.Decode(&multiple); err == nil {
		result = T(multiple)
		return result, nil
	}

	return result, fmt.Errorf("failed to unmarshal: expected string or slice of strings")
}

// UnmarshalStringOrSlice is a generic function to unmarshal YAML into either a single string or a slice of strings
func UnmarshalSingleOrSlice[T []any](value *yaml.Node) (T, error) {
	var result T

	// Try unmarshalling into a single string
	var single T
	if err := value.Decode(&single); err == nil {
		result = append(result, single)
		return result, nil
	}

	// Try unmarshalling into a slice of strings
	var multiple []string
	if err := value.Decode(&multiple); err == nil {
		result = T(multiple)
		return result, nil
	}

	return result, fmt.Errorf("failed to unmarshal: expected string or slice of strings")
}
