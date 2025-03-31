package env

import (
	"os"
)

// ToSlice converts the Env to a slice of `key=value` strings.
// Each key-value pair in the Env is formatted as "key=value" and added to the slice.
func (e Env) ToSlice() []string {
	slice := make([]string, 0, len(e))

	for k, v := range e {
		slice = append(slice, k+"="+v)
	}

	return slice
}

// ToEnv sets all the key-value pairs in the Env as environment variables using os.Setenv.
// It returns an error if setting any environment variable fails.
func (e Env) ToEnv() error {
	for k, v := range e {
		if err := os.Setenv(k, v); err != nil {
			return err
		}
	}

	return nil
}
