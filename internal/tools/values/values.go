// Package values provides functionality for managing tool configuration values.
package values

import "maps"

// Values represents a map of string keys to any values.
type Values map[string]any

// Merge combines value maps while preserving receiver values on conflicts.
// It mutates each supplied map in sequence and aliases the receiver to the final map.
func (v *Values) Merge(values ...Values) {
	for _, value := range values {
		if value == nil {
			continue
		}

		maps.Copy(value, *v)

		*v = value
	}
}
