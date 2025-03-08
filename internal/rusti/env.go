package rusti

import (
	"fmt"
	"path/filepath"
)

// Env represents the environment variables for a Rust binary.
type Env map[string]string

// ToSlice converts the environment variables to a slice.
func (e Env) ToSlice() []string {
	env := make([]string, 0, len(e))
	for k, v := range e {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}

	return env
}

// Append appends the given environment variables to the current environment.
func (e *Env) Append(env Env) {
	for k, v := range env {
		(*e)[k] = v
	}
}

// Default sets the default environment variables for a Rust binary.
func (e *Env) Default(dir string) {
	*e = Env{
		"CARGO_HOME":  filepath.Join(dir, ".cargo"),
		"RUSTUP_HOME": filepath.Join(dir, ".rustup"),
	}
}
