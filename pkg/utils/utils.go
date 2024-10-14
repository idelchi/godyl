// Package utils provides utility functions for handling common tasks
// such as setting default values, checking emptiness, and normalizing map keys.
package utils

import (
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// SetIfEmpty sets the value of input to the specified value if it is empty.
// S must be a comparable type.
func SetIfEmpty[S comparable](input *S, value S) {
	var zero S

	if *input == zero {
		*input = value
	}
}

// SetSliceIfNil sets the value of input to the provided values slice if input is nil.
// S is constrained to slices of any type T.
func SetSliceIfNil[S ~[]T, T any](input *S, values ...T) {
	if *input == nil {
		*input = append([]T(nil), values...)
	}
}

// IsEmpty checks if the input value is empty.
// S must be a comparable type.
func IsEmpty[S comparable](input S) bool {
	var zero S

	return input == zero
}

// NormalizeMap normalizes the keys of a map to title case recursively.
// If the value is another map, it will recursively normalize its keys as well.
func NormalizeMap(m map[string]any) map[string]any {
	normalizedMap := make(map[string]any)
	c := cases.Title(language.English)

	for key, value := range m {
		upperKey := c.String(key)

		switch v := value.(type) {
		case map[string]any:
			normalizedMap[upperKey] = NormalizeMap(v)
		default:
			normalizedMap[upperKey] = v
		}
	}

	return normalizedMap
}
