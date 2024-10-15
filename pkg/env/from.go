// Package env provides utilities for working with environment variables,
// including methods to normalize, merge, retrieve, and manipulate them in a map-like structure.
package env

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// FromEnv returns the current environment variables as an Env.
func FromEnv() Env {
	env, _ := FromSlice(os.Environ()...)

	return env.Normalized()
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

func FromDotEnv(path string) (Env, error) {
	env, err := godotenv.Read(path)
	if err != nil {
		return nil, fmt.Errorf("loading dotenv: %w", err)
	}

	return Env(env).Normalized(), nil
}
