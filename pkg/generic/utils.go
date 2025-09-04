// Package generic provides lightweight, generic helpers for slices, zero checks,
// path expansion, deep copying, and basic filtering.
package generic

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/jinzhu/copier"
)

// AnyNil checks if any pointer in the provided slice is nil.
func AnyNil[T any](ptrs ...*T) bool {
	// Check if any pointer in the slice is nil.
	for _, ptr := range ptrs {
		if ptr == nil {
			return true
		}
	}

	return false
}

// IsSliceNilOrEmpty checks if a slice pointer is nil or points to an empty slice.
func IsSliceNilOrEmpty[T ~[]E, E any](ptr *T) bool {
	return IsSliceNil(ptr) || IsSliceEmpty(*ptr)
}

// IsSliceNil checks if a slice pointer is nil.
func IsSliceNil[T ~[]E, E any](ptr *T) bool {
	return ptr == nil
}

// IsSliceEmpty checks if a slice is empty.
func IsSliceEmpty[T ~[]E, E any](ptr T) bool {
	return len(ptr) == 0
}

// IsURL validates a string as a properly formatted URL.
// Returns true if the string contains both a scheme and host.
func IsURL(str string) bool {
	u, err := url.Parse(str)

	return err == nil && u.Scheme != "" && u.Host != ""
}

// IsZero checks if a value equals its type's zero value.
// Works with any comparable type (numbers, strings, etc.).
func IsZero[S comparable](input S) bool {
	var zero S

	return input == zero
}

// SetIfZero conditionally updates a pointer's value.
// Sets the pointed-to value to a new value only if the current
// value equals the type's zero value.
func SetIfZero[S comparable](input *S, value S) {
	if IsZero(*input) {
		*input = value
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

// DeepCopy copies the contents of the object and returns it.
func DeepCopy[T any](src T) (dst T, err error) {
	if err := copier.CopyWithOption(&dst, &src, copier.Option{DeepCopy: true, CaseSensitive: true}); err != nil {
		return dst, fmt.Errorf("copying object: %w", err)
	}

	return dst, nil
}

// DeepCopyPtr copies a pointer type object and returns a new pointer to the copied object.
func DeepCopyPtr[T any](src *T) (*T, error) {
	if src == nil {
		return nil, nil //nolint:nilnil 	// Return nil if source is nil
	}

	dst := new(T) // Create a new non-nil destination

	if err := copier.CopyWithOption(dst, src, copier.Option{DeepCopy: true, CaseSensitive: true}); err != nil {
		return nil, fmt.Errorf("copying object: %w", err)
	}

	return dst, nil
}

// PickByIndices returns a slice of elements from the input slice
// based on the provided indices.
func PickByIndices[T any](s []T, indices []int) []T {
	out := make([]T, 0, len(indices))

	for _, i := range indices {
		out = append(out, s[i])
	}

	return out
}
