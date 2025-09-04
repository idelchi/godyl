// Package values provides functionality for managing tool configuration values.
package values

import "maps"

// Values represents a map of string keys to any values.
type Values map[string]any

// Merge combines multiple environments into the current one.
// Copies values from the provided environments into this one,
// preserving existing values in case of key conflicts.
func (v *Values) Merge(values ...Values) {
	for _, value := range values {
		if value == nil {
			continue
		}

		maps.Copy(value, *v)

		*v = value
	}
}
