// Package env provides utilities for working with environment variables,
// including methods to normalize, merge, retrieve, and manipulate them in a map-like structure.
package env

import (
	"os"
	"runtime"
	"strings"
)

// Normalized returns a copy of the Env with all keys normalized to uppercase on Windows.
func (e Env) Normalized() Env {
	if runtime.GOOS == "windows" {
		return e.Normalize()
	}

	return e
}

// Normalized returns a copy of the Env with all keys normalized to uppercase on Windows.
func (e Env) Normalize() Env {
	normalized := make(Env, len(e))

	for k, v := range e {
		normalized[strings.ToUpper(k)] = v
	}

	return normalized
}

func (e *Env) Expand() {
	for k, v := range *e {
		(*e)[k] = os.ExpandEnv(v)
	}
}
