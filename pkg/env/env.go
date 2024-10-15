// Package env provides utilities for working with environment variables,
// including methods to normalize, merge, retrieve, and manipulate them in a map-like structure.
package env

import (
	"errors"
)

var (
	ErrEnvVarNotFound = errors.New("environment variable not found")
	ErrEnvMalformed   = errors.New("environment variable is malformed")
)

// Env represents a map of environment variables with string keys and values.
type Env map[string]string
