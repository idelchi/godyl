// Package fallbacks provides functionality for managing tool fallback sources.
package fallbacks

import (
	"fmt"

	"github.com/goccy/go-yaml/ast"

	"github.com/idelchi/godyl/internal/tools/sources"
	"github.com/idelchi/godyl/pkg/unmarshal"
)

// Fallbacks represents a collection of fallback sources for the tool.
// It can either be a single source type or a slice of source types, allowing flexibility
// in specifying multiple fallback methods if the primary source fails.
type Fallbacks []sources.Type

// UnmarshalYAML implements custom unmarshaling for Tags,
// allowing the field to be either a single string or a list of strings.
func (f *Fallbacks) UnmarshalYAML(node ast.Node) (err error) {
	*f, err = unmarshal.SingleOrSlice[sources.Type](node)
	if err != nil {
		return fmt.Errorf("unmarshaling fallbacks: %w", err)
	}

	return nil
}

// Compact removes duplicate elements from a slice while preserving order.
func Compact[T comparable](s []T) []T {
	seen := make(map[T]bool)

	result := make([]T, 0, len(s))
	for _, v := range s {
		if !seen[v] {
			seen[v] = true
			result = append(result, v)
		}
	}

	return result
}

// Compacted returns a new Fallbacks slice with duplicates removed.
func (f Fallbacks) Compacted() Fallbacks {
	return Compact(f)
}

// Build creates a source type list by prepending the given type to fallbacks.
func (f Fallbacks) Build(sourceType sources.Type) []sources.Type {
	// Prepend sourceType to the existing fallbacks and remove duplicates
	return append(Fallbacks{sourceType}, f...).Compacted()
}
