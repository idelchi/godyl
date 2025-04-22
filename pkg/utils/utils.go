// Package utils provides common utility functions for general purpose tasks.
// It includes helpers for URL validation, zero value checking, slice manipulation,
// map normalization, path handling, and deep copying objects.
package utils

import (
	"encoding/json"
	"fmt"
	"maps"
	"net/url"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/jinzhu/copier"
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
// Creates a new slice with the provided values only if the pointer is nil. Safe for any slice type.
func SetSliceIfNil[S ~[]T, T any](input *S, values ...T) {
	if *input == nil {
		*input = slices.Clone(values)
	}
}

// SetSliceIfZero initializes an empty or nil slice pointer.
// Creates a new slice with the provided values if the current
// slice is nil or has zero length.
func SetSliceIfZero[S ~[]T, T any](input *S, values ...T) {
	if *input == nil || len(*input) == 0 {
		*input = slices.Clone(values)
	}
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
		return dst, fmt.Errorf("marshaling object: %w", err)
	}

	err = json.Unmarshal(bytes, &dst)
	if err != nil {
		return dst, fmt.Errorf("unmarshaling object: %w", err)
	}

	return dst, nil
}

// DeepCopyInto copies the contents of the object and returns it.
func DeepCopy[T any](src T) (dst T, err error) {
	if err := copier.CopyWithOption(&dst, &src, copier.Option{DeepCopy: true, CaseSensitive: true}); err != nil {
		return dst, fmt.Errorf("copying object: %w", err)
	}

	return dst, nil
}

// MergeMaps constructs a merge between primary and secondary where primary values have priority.
// The primary map is modified in place.
func MergeMapsInPlace[T ~map[string]any](primary, secondary T) error {
	copied, err := DeepCopy(secondary)
	if err != nil {
		return fmt.Errorf("copying secondary map: %w", err)
	}

	copyMapIfNotExist(primary, copied)

	return nil
}

// MergeMaps constructs a merge between primary and secondary where primary values have priority.
func MergeMaps[T ~map[string]any](primary, secondary T) (T, error) {
	copied, err := DeepCopy(secondary)
	if err != nil {
		return nil, fmt.Errorf("copying secondary map: %w", err)
	}

	maps.Copy(copied, primary)

	return copied, nil
}

// copyMapIfNotExist copies all key/value pairs in src adding them to dst,
// but only if the key doesn't already exist in dst.
func copyMapIfNotExist[M1 ~map[K]V, M2 ~map[K]V, K comparable, V any](dst M1, src M2) {
	for k, v := range src {
		if _, exists := dst[k]; !exists {
			dst[k] = v
		}
	}
}

func PickByIndices[T any](s []T, indices []int) []T {
	out := make([]T, 0, len(indices))
	for _, i := range indices {
		out = append(out, s[i])
	}
	return out
}
