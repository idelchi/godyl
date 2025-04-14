package env

import (
	"fmt"
	"maps"
	"strings"
)

// Add parses and adds a key-value pair to the environment.
// Takes a string in the format "key=value" and adds it to the Env map.
// Returns an error if the string doesn't contain exactly one '=' separator.
func (e *Env) Add(kv string) error {
	const expectedParts = 2

	parts := strings.SplitN(kv, "=", expectedParts)
	if len(parts) != expectedParts {
		return fmt.Errorf("%w: %q", ErrEnvMalformed, kv)
	}

	(*e)[parts[0]] = parts[1]

	return nil
}

// Merge combines multiple environments into the current one.
// Copies values from the provided environments into this one,
// preserving existing values in case of key conflicts.
func (e *Env) Merge(envs ...Env) {
	for _, env := range envs {
		maps.Copy(env, *e)
		*e = env
	}
}

// Merged creates a new environment by combining multiple environments.
// Returns a new Env containing all values from this environment plus
// any non-conflicting values from the provided environments.
// Does not modify the original environment.
func (e Env) Merged(envs ...Env) Env {
	merged := maps.Clone(e)

	for _, env := range envs {
		merged.Merge(env)
	}

	return merged
}
