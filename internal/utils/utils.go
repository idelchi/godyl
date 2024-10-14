package utils

import (
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func SetIfEmpty[S comparable](input *S, value S) {
	var empty S

	if *input == empty {
		*input = value
	}
}

func SetSliceIfNil[S ~[]T, T any](input *S, values ...T) {
	if *input == nil {
		*input = append([]T(nil), values...)
	}
}

func IsEmpty[S comparable](input S) bool {
	var empty S

	return input == empty
}

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
