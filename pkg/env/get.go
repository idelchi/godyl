package env

import (
	"fmt"
)

// MustGet retrieves an environment variable's value by key.
// Returns an error if the key doesn't exist in the environment.
func (e *Env) MustGet(key string) (string, error) {
	if !e.Exists(key) {
		return "", fmt.Errorf("%w: %q", ErrEnvVarNotFound, key)
	}

	return e.Get(key), nil
}

// Get retrieves an environment variable's value by key.
// Returns an empty string if the key doesn't exist.
func (e *Env) Get(key string) string {
	return (*e)[key]
}

// GetAsEnv retrieves an environment variable's value as a formatted string.
func (e *Env) GetAsEnv(key string) string {
	if !e.Exists(key) {
		return ""
	}

	return fmt.Sprintf("%s=%s", key, (*e)[key])
}

// Exists checks if a key exists in the environment.
// Returns true if the key exists, false otherwise.
func (e *Env) Exists(key string) bool {
	_, ok := (*e)[key]

	return ok
}

// GetOrDefault retrieves an environment value with fallback.
// Returns the environment value if the key exists,
// otherwise returns the provided default value.
func (e *Env) GetOrDefault(key, defaultValue string) string {
	if !e.Exists(key) {
		return defaultValue
	}

	return e.Get(key)
}

// GetAny retrieves the first available environment value.
// Tries each key in order until a value is found, returning
// an empty string if none of the keys exist.
func (e *Env) GetAny(keys ...string) string {
	for _, key := range keys {
		if value, err := e.MustGet(key); err == nil {
			return value
		}
	}

	return ""
}

// Predicate defines a function type that takes two strings
// and returns a boolean indicating whether the predicate is satisfied.
type Predicate func(string, string) bool

// GetWithPredicates retrieves all environment values that match any of the provided predicates.
func (e *Env) GetWithPredicates(predicates ...Predicate) Env {
	matches := make(Env)

	for key, value := range *e {
		for _, predicate := range predicates {
			if predicate(key, value) {
				matches[key] = value

				break
			}
		}
	}

	return matches
}
