// Package env provides utilities for working with environment variables,
// including methods to normalize, merge, retrieve, and manipulate them in a map-like structure.
package env

import (
	"fmt"
	"maps"
	"strings"
)

// Add splits a `key=value` string and adds it to the Env map.
// It returns an error if the input is not properly formatted.
func (e *Env) Add(kv string) error {
	parts := strings.SplitN(kv, "=", 2)
	if len(parts) != 2 {
		return fmt.Errorf("%w: %q", ErrEnvMalformed, kv)
	}

	(*e)[parts[0]] = parts[1]
	return nil
}

// Merge merges another Env into the current Env, without overwriting existing keys in the current Env.
func (e *Env) Merge(envs ...Env) {
	for _, env := range envs {
		maps.Copy(env, *e)

		*e = env
	}
}

// Merged returns a new Env by merging the given Env into the current Env,
// without overwriting existing keys in the original Env.
func (e Env) Merged(envs ...Env) Env {
	merged := maps.Clone(e)

	for _, env := range envs {
		merged.Merge(env)
	}

	return merged
}
