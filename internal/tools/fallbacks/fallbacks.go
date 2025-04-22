// Package fallbacks provides functionality for managing tool fallback sources.
package fallbacks

import (
	"slices"

	"github.com/idelchi/godyl/internal/tools/sources"
	"github.com/idelchi/godyl/pkg/unmarshal"

	"gopkg.in/yaml.v3"
)

// Fallbacks represents a collection of fallback sources for the tool.
// It can either be a single source type or a slice of source types, allowing flexibility
// in specifying multiple fallback methods if the primary source fails.
type Fallbacks []sources.Type

// UnmarshalYAML implements custom unmarshaling for Tags,
// allowing the field to be either a single string or a list of strings.
func (f *Fallbacks) UnmarshalYAML(value *yaml.Node) error {
	result, err := unmarshal.SingleOrSlice[sources.Type](value, false)
	if err != nil {
		return err
	}

	*f = result

	return nil
}

func (f Fallbacks) Compacted() Fallbacks {
	return slices.Compact(f)
}

func (f Fallbacks) Build(sourceType sources.Type) []sources.Type {
	// Prepend sourceType to the existing fallbacks and remove duplicates
	return append(Fallbacks{sourceType}, f...).Compacted()
}
