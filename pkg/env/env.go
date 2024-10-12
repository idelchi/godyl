// Package env provides utilities for working with environment variables,
// including methods to normalize, merge, retrieve, and manipulate them in a map-like structure.
package env

import (
	"errors"
	"fmt"
	"maps"
	"os"
	"runtime"
	"strings"

	"github.com/joho/godotenv"
)

var (
	ErrEnvVarNotFound = errors.New("environment variable not found")
	ErrDotEnvLoad     = errors.New("failed to load dotenv")
)

// Env represents a map of environment variables with string keys and values.
type Env map[string]string

// FromEnv returns the current environment variables as an Env.
func FromEnv() Env {
	return FromSlice(os.Environ()...)
}

// FromSlice constructs an Env from a slice of key=value strings.
func FromSlice(slice ...string) Env {
	e := make(Env, len(slice))

	for _, v := range slice {
		kv := strings.SplitN(v, "=", 2)
		e[kv[0]] = kv[1]
	}

	return e
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

// ToSlice converts the Env to a slice of `key=value“ strings.
func (e Env) ToSlice() []string {
	slice := make([]string, 0, len(e))

	for k, v := range e {
		slice = append(slice, k+"="+v)
	}

	return slice
}

// Add adds `key=value“ pairs from a slice of strings to the Env.
func (e *Env) Add(slice ...string) {
	e.Merge(FromSlice(slice...))
}

// Merge merges another Env into the current Env, without overwriting existing keys.
func (e *Env) Merge(env Env) {
	maps.Copy(env, *e)

	*e = env
}

// Merged returns a new Env by merging the given Env into the current Env,
// without overwriting existing keys in the original Env.
func (e Env) Merged(env Env) Env {
	merged := maps.Clone(e)

	merged.Merge(env)

	return merged
}

func FromDotEnv(path string) (Env, error) {
	env, err := godotenv.Read(path)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrDotEnvLoad, err)
	}

	return Env(env), nil
}
