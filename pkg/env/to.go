// Package env provides utilities for working with environment variables,
// including methods to normalize, merge, retrieve, and manipulate them in a map-like structure.
package env

import "os"

// ToSlice converts the Env to a slice of `key=value“ strings.
func (e Env) ToSlice() []string {
	slice := make([]string, 0, len(e))

	for k, v := range e {
		slice = append(slice, k+"="+v)
	}

	return slice
}

func (e Env) ToEnv() error {
	for k, v := range e {
		if err := os.Setenv(k, v); err != nil {
			return err
		}
	}

	return nil
}
