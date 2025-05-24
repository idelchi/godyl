package env

import (
	"errors"
	"fmt"
	"maps"
	"strings"
)

// Add parses and adds a key-value pairs to the environment.
// Takes strings in the format "key=value" and adds it to the Env map.
// Returns errors for strings that do not contain an '=' or have an empty key.
func (e *Env) Add(keyValues ...string) error {
	var errs []error

	for _, keyValue := range keyValues {
		key, value, found := strings.Cut(keyValue, "=")

		if !found {
			errs = append(errs, fmt.Errorf("%w: %q", ErrEnvMalformed, keyValue))

			continue
		}

		if key == "" {
			errs = append(errs, fmt.Errorf("%w: empty key in pair %q", ErrEnvMalformed, keyValue))

			continue
		}

		(*e)[key] = value
	}

	return errors.Join(errs...)
}

// AddPair adds a key-value pair to the environment.
// Takes a key and value string and adds it to the Env map.
// Returns an error if the key is empty.
func (e *Env) AddPair(key, value string) error {
	if key == "" {
		return fmt.Errorf("%w: empty key", ErrEnvMalformed)
	}

	(*e)[key] = value

	return nil
}

// Delete removes an environment variable by key.
func (e *Env) Delete(key string) {
	delete(*e, key)
}

// Merge combines multiple environments into the current one.
// Copies values from the provided environments,
// preserving existing values in case of key conflicts.
func (e *Env) Merge(envs ...Env) {
	nEnv := Env{}

	for _, env := range append(envs, *e) {
		if env == nil {
			continue
		}

		maps.Copy(nEnv, env)
	}

	*e = nEnv
}

// MergedWith creates a new environment by combining multiple environments.
// Returns a new Env containing all values from this environment plus
// any non-conflicting values from the provided environments.
// Does not modify the original environment.
func (e *Env) MergedWith(envs ...Env) Env {
	merged := maps.Clone(*e)

	for _, env := range envs {
		merged.Merge(env)
	}

	return merged
}
