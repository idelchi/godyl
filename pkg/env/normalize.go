package env

import (
	"os"
	"runtime"
	"strings"
)

// Normalized returns a platform-appropriate normalized environment.
// On Windows, converts all keys to uppercase for case-insensitive
// compatibility. On other systems, returns the environment unchanged.
func (e *Env) Normalized() Env {
	if runtime.GOOS == "windows" {
		return e.toUpper()
	}

	return *e
}

// toUpper creates a new environment with uppercase keys.
// Returns a copy of the environment with all keys converted to
// uppercase, which ensures consistent behavior on case-insensitive systems like Windows.
func (e *Env) toUpper() Env {
	normalized := make(Env, len(*e))

	for k, v := range *e {
		normalized[strings.ToUpper(k)] = v
	}

	return normalized
}

// Expand processes environment variable references in values.
// Replaces ${var} or $var patterns in each value with the
// corresponding environment variable value. Modifies the
// environment in place.
func (e *Env) Expand() {
	for k, v := range *e {
		(*e)[k] = os.ExpandEnv(v)
	}
}
