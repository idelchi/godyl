package goi

import (
	"fmt"
	"maps"
	"path/filepath"
)

// Env represents a map of environment variables, where each key is an environment variable name,
// and each value is the corresponding value of that variable.
type Env map[string]string

// ToSlice converts the Env map into a slice of strings in the format "KEY=VALUE" suitable for
// passing to external processes.
func (e *Env) ToSlice() []string {
	env := make([]string, 0, len(*e))
	for k, v := range *e {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}

	return env
}

// Append merges the given Env into the current Env, overwriting any existing keys with the same names.
func (e *Env) Append(env Env) {
	maps.Copy(*e, env)
}

// Default sets up default environment variables typically used in Go projects.
// It configures GOMODCACHE, GOCACHE, and GOPATH based on the provided directory.
func (e *Env) Default(dir string) {
	*e = Env{
		"GOMODCACHE": filepath.Join(dir, ".cache"),
		"GOCACHE":    filepath.Join(dir, ".cache"),
		"GOPATH":     filepath.Join(dir, ".path"),
	}
}
