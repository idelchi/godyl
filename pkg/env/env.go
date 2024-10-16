package env

import (
	"errors"
)

var (
	// ErrEnvVarNotFound is returned when an expected environment variable is not found.
	ErrEnvVarNotFound = errors.New("environment variable not found")

	// ErrEnvMalformed is returned when an environment variable is malformed or contains invalid data.
	ErrEnvMalformed = errors.New("environment variable is malformed")
)

// Env represents a map of environment variables, where the keys and values are strings.
// It can be used to access and manage environment settings.
type Env map[string]string
