package env

import (
	"fmt"
	"regexp"
	"strings"
)

// MustGet retrieves an environment variable's value by key.
// Returns an error if the key doesn't exist in the environment.
func (e Env) MustGet(key string) (string, error) {
	if v, ok := e[key]; ok {
		return v, nil
	}

	return "", fmt.Errorf("%w: %q", ErrEnvVarNotFound, key)
}

// Get retrieves an environment variable's value by key.
// Returns an empty string if the key doesn't exist.
func (e Env) Get(key string) string {
	v, _ := e.MustGet(key)

	return v
}

// GetAny retrieves the first available environment value.
// Tries each key in order until a value is found, returning
// an empty string if none of the keys exist.
func (e Env) GetAny(keys ...string) string {
	for _, key := range keys {
		if value, err := e.MustGet(key); err == nil {
			return value
		}
	}

	return ""
}

// Has checks if a key exists in the environment.
// Returns true if the key exists, false otherwise.
func (e Env) Has(key string) bool {
	_, ok := e[key]
	return ok
}

// GetOrDefault retrieves an environment value with fallback.
// Returns the environment value if the key exists,
// otherwise returns the provided default value.
func (e Env) GetOrDefault(key, defaultValue string) string {
	if value, err := e.MustGet(key); err != nil {
		return defaultValue
	} else {
		return value
	}
}

// GetAll filters environment variables using a predicate.
// Returns a new Env containing only the key-value pairs
// for which the predicate function returns true.
func (e Env) GetAll(predicate func(key, value string) bool) Env {
	result := make(Env)

	for k, v := range e {
		if predicate(k, v) {
			result[k] = v
		}
	}

	return result
}

// GetAllWithPrefix returns environment variables by prefix.
// Returns a new Env containing only variables whose keys
// start with the specified prefix.
func (e Env) GetAllWithPrefix(prefix string) Env {
	return e.GetAll(func(key, _ string) bool {
		return strings.HasPrefix(key, prefix)
	})
}

// GetAllWithSuffix returns environment variables by suffix.
// Returns a new Env containing only variables whose keys
// end with the specified suffix.
func (e Env) GetAllWithSuffix(suffix string) Env {
	return e.GetAll(func(key, _ string) bool {
		return strings.HasSuffix(key, suffix)
	})
}

// GetAllMatching returns environment variables by regex pattern.
// Returns a new Env containing only variables whose keys match
// the specified regex pattern. Returns an error if the pattern
// is invalid.
func (e Env) GetAllMatching(pattern string) (Env, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("invalid regex pattern: %w", err)
	}

	return e.GetAll(func(key, _ string) bool {
		return re.MatchString(key)
	}), nil
}

// GetAllWithValue returns environment variables by value.
// Returns a new Env containing only variables that have
// the exact specified value.
func (e Env) GetAllWithValue(value string) Env {
	return e.GetAll(func(_, v string) bool {
		return v == value
	})
}
