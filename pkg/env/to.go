package env

import (
	"os"
)

// ToSlice converts the environment to a string slice.
// Returns a slice where each element is a "key=value" string
// representing an environment variable.
func (e Env) ToSlice() []string {
	slice := make([]string, 0, len(e))

	for k, v := range e {
		slice = append(slice, k+"="+v)
	}

	return slice
}

// ToEnv applies the environment to the current process.
// Sets each key-value pair as a process environment variable
// using os.Setenv. Returns an error if any variable cannot
// be set.
func (e Env) ToEnv() error {
	for k, v := range e {
		if err := os.Setenv(k, v); err != nil {
			return err
		}
	}

	return nil
}
