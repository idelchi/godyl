// Package env provides utilities for working with environment variables,
// including methods to normalize, merge, retrieve, and manipulate them in a map-like structure.
package env

import (
	"errors"
	"fmt"
	"maps"
	"os"
	"regexp"
	"runtime"
	"strings"

	"github.com/joho/godotenv"
)

var (
	ErrEnvVarNotFound = errors.New("environment variable not found")
	ErrEnvMalformed   = errors.New("environment variable is malformed")
)

// Env represents a map of environment variables with string keys and values.
type Env map[string]string

// FromEnv returns the current environment variables as an Env.
func FromEnv() Env {
	env, _ := FromSlice(os.Environ()...)

	return env.Normalized()
}

// Normalized returns a copy of the Env with all keys normalized to uppercase on Windows.
func (e Env) Normalized() Env {
	if runtime.GOOS == "windows" {
		normalized := make(Env, len(e))

		for k, v := range e {
			normalized[strings.ToUpper(k)] = v
		}

		return normalized
	}

	return e
}

// Get retrieves the value associated with the given key or an error if the key is not found.
func (e Env) Get(key string) (string, error) {
	if v, ok := e[key]; ok {
		return v, nil
	}

	return "", fmt.Errorf("%w: %q", ErrEnvVarNotFound, key)
}

func (e Env) GetOrDefault(key, defaultValue string) string {
	if value, err := e.Get(key); err != nil {
		return defaultValue
	} else {
		return value
	}
}

// ToSlice converts the Env to a slice of `key=value“ strings.
func (e Env) ToSlice() []string {
	slice := make([]string, 0, len(e))

	for k, v := range e {
		slice = append(slice, k+"="+v)
	}

	return slice
}

// FromSlice constructs an Env from a slice of `key=value` strings.
func FromSlice(slice ...string) (Env, error) {
	e := make(Env, len(slice))

	for _, v := range slice {
		if err := e.Add(v); err != nil {
			return nil, err
		}
	}

	return e.Normalized(), nil
}

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

func FromDotEnv(path string) (Env, error) {
	env, err := godotenv.Read(path)
	if err != nil {
		return nil, fmt.Errorf("loading dotenv: %w", err)
	}

	return Env(env).Normalized(), nil
}

// GetAll returns a new Env containing all key-value pairs that satisfy the given predicate function.
func (e Env) GetAll(predicate func(key, value string) bool) Env {
	result := make(Env)
	for k, v := range e {
		if predicate(k, v) {
			result[k] = v
		}
	}
	return result
}

// Examples of using GetAll with different predicates:

// GetAllWithPrefix returns all environment variables with keys starting with the given prefix.
func (e Env) GetAllWithPrefix(prefix string) Env {
	return e.GetAll(func(key, _ string) bool {
		return strings.HasPrefix(key, prefix)
	})
}

// GetAllWithSuffix returns all environment variables with keys ending with the given suffix.
func (e Env) GetAllWithSuffix(suffix string) Env {
	return e.GetAll(func(key, _ string) bool {
		return strings.HasSuffix(key, suffix)
	})
}

// GetAllMatching returns all environment variables with keys matching the given regex pattern.
func (e Env) GetAllMatching(pattern string) (Env, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("invalid regex pattern: %w", err)
	}
	return e.GetAll(func(key, _ string) bool {
		return re.MatchString(key)
	}), nil
}

// GetAllWithValue returns all environment variables with the exact given value.
func (e Env) GetAllWithValue(value string) Env {
	return e.GetAll(func(_, v string) bool {
		return v == value
	})
}
