package utils

import (
	"encoding/json"
	"fmt"
	"maps"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/jinzhu/copier"

	"github.com/mitchellh/copystructure"
)

// IsURL validates a string as a properly formatted URL.
// Returns true if the string contains both a scheme and host.
func IsURL(str string) bool {
	u, err := url.Parse(str)

	return err == nil && u.Scheme != "" && u.Host != ""
}

// IsZeroValue checks if a value equals its type's zero value.
// Works with any comparable type (numbers, strings, etc.).
func IsZeroValue[S comparable](input S) bool {
	var zero S

	return input == zero
}

// SetIfZeroValue conditionally updates a pointer's value.
// Sets the pointed-to value to a new value only if the current
// value equals the type's zero value.
func SetIfZeroValue[S comparable](input *S, value S) {
	if IsZeroValue(*input) {
		*input = value
	}
}

// SetSliceIfNil initializes a nil slice pointer.
// Creates a new slice with the provided values only if the
// pointer is nil. Safe for any slice type.
func SetSliceIfNil[S ~[]T, T any](input *S, values ...T) {
	if *input == nil {
		*input = append([]T(nil), values...)
	}
}

// SetSliceIfZero initializes an empty or nil slice pointer.
// Creates a new slice with the provided values if the current
// slice is nil or has zero length.
func SetSliceIfZero[S ~[]T, T any](input *S, values ...T) {
	if *input == nil || len(*input) == 0 {
		*input = append([]T(nil), values...)
	}
}

// SetMapIfNil initializes a nil map pointer.
// Creates a new map with the provided key-value pairs only if
// the pointer is nil. Safe for any map type.
func SetMapIfNil[M ~map[K]V, K comparable, V any](input *M, values M) {
	if *input == nil {
		*input = make(M, len(values))
		for k, v := range values {
			(*input)[k] = v
		}
	}
}

// CopySlice creates an independent copy of a slice.
// Returns a new slice with the same elements. Returns nil
// if the input slice is nil.
func CopySlice[S ~[]T, T any](input S) S {
	if input == nil {
		return nil
	}

	result := make(S, len(input))
	copy(result, input)
	return result
}

// CopyMap creates an independent copy of a map.
// Returns a new map with the same key-value pairs. Returns
// nil if the input map is nil.
func CopyMap[M ~map[K]V, K comparable, V any](input M) M {
	if input == nil {
		return nil
	}

	result := make(M, len(input))
	for k, v := range input {
		result[k] = v
	}
	return result
}

// NormalizeMap converts map keys to title case recursively.
// Creates a new map with all string keys converted to title case.
// Handles nested maps by recursively normalizing their keys as well.
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

// DeepMergeMapsWithoutOverwrite combines two maps preserving existing values.
// Merges the second map into the first, keeping original values in case
// of key conflicts. Handles nested maps by recursively merging them.
// Creates new map instances for nested maps to avoid shared references.
func DeepMergeMapsWithoutOverwrite(first, second map[string]any) {
	for key, secondVal := range second {
		if firstVal, exists := first[key]; exists {
			// If both values are maps, recursively merge them
			if firstMap, ok1 := firstVal.(map[string]any); ok1 {
				if secondMap, ok2 := secondVal.(map[string]any); ok2 {
					// If we already have an existing map at this key, make a clone
					// to avoid modifying the original reference, then merge into it
					clonedMap := maps.Clone(firstMap)
					DeepMergeMapsWithoutOverwrite(clonedMap, secondMap)
					first[key] = clonedMap
				}
			} // If the key exists but isn't a map, do nothing (keep the original value)
		} else {
			// If the key doesn't exist in first, add it from second
			// For maps, clone to avoid shared references
			if secondMap, ok := secondVal.(map[string]any); ok {
				first[key] = maps.Clone(secondMap)
			} else {
				first[key] = secondVal
			}
		}
	}
}

// ExpandHome resolves home directory references in paths.
// Replaces leading ~ with the user's home directory path.
// Returns the original path if expansion fails or isn't needed.
func ExpandHome(path string) string {
	if !strings.HasPrefix(path, "~") {
		return path
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return path
	}

	return filepath.Join(home, path[1:])
}

// DeepCopyWithMarshaling creates a deep copy of a source object using JSON marshaling.
func DeepCopyWithMarshaling[T any](src T) (dst T, err error) {
	bytes, err := json.Marshal(src)
	if err != nil {
		return dst, fmt.Errorf("Unable to marshal src: %s", err)
	}

	err = json.Unmarshal(bytes, &dst)
	if err != nil {
		return dst, fmt.Errorf("Unable to unmarshal into dst: %s", err)
	}
	return dst, nil
}

// DeepCopyInto copies the contents of one struct into another.
func DeepCopy[T any](src T) (dst T, err error) {
	if err := copier.CopyWithOption(&dst, &src, copier.Option{DeepCopy: true}); err != nil {
		return dst, fmt.Errorf("Unable to copy src to dst: %s", err)
	}

	return dst, nil
}

func DeepCopyInto[T any, V any](src T) (dst V, err error) {
	if err := copier.CopyWithOption(&dst, &src, copier.Option{DeepCopy: true}); err != nil {
		return dst, fmt.Errorf("Unable to copy src to dst: %s", err)
	}

	return dst, nil
}

// DeepCopyWithCs creates a deep copy of a source object using the copystructure package.
func DeepCopyWithCs[T any](src T) (dst T, err error) {
	if copy, err := copystructure.Copy(src); err != nil {
		return dst, fmt.Errorf("Unable to copy src to dst: %s", err)
	} else {
		dst, ok := copy.(T)
		if !ok {
			return dst, fmt.Errorf("Unable to cast copied structure to destination type")
		}

		return dst, nil
	}
}
