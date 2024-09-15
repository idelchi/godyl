package env

import (
	"fmt"
	"os"
	"slices"
	"strings"
)

// AsSlice converts the environment to a string slice.
// Returns a sorted slice where each element is a "key=value" string
// representing an environment variable.
func (e *Env) AsSlice() []string {
	slice := make([]string, 0, len(*e))

	for k, v := range *e {
		slice = append(slice, strings.Join([]string{k, v}, "="))
	}

	slices.Sort(slice)

	return slice
}

// Keys returns the keys of the environment variables, sorted alphabetically.
func (e *Env) Keys() []string {
	keys := make([]string, 0, len(*e))

	for k := range *e {
		keys = append(keys, k)
	}

	slices.Sort(keys)

	return keys
}

// Export applies the environment to the current process.
// Expands any variable references in the values.
func (e *Env) Export() error {
	e.Expand()

	for k, v := range *e {
		if err := os.Setenv(k, v); err != nil {
			return fmt.Errorf("setting env var %q: %w", k, err)
		}
	}

	return nil
}
