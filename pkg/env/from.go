package env

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// FromEnv creates an Env from the current process environment.
// Retrieves all environment variables using os.Environ, normalizes
// their format, and returns them as an Env map.
func FromEnv() Env {
	env, _ := FromSlice(os.Environ()...)

	return env.Normalized()
}

// FromSlice creates an Env from a list of key-value strings.
// Each string should be in the format "key=value". Returns an error
// if any string is malformed. Normalizes the resulting environment.
func FromSlice(slice ...string) (Env, error) {
	e := make(Env, len(slice))

	for _, v := range slice {
		if err := e.Add(v); err != nil {
			return nil, err
		}
	}

	return e.Normalized(), nil
}

// FromDotEnv loads environment variables from a .env file.
// Reads and parses the file at the given path using godotenv.
// Returns the variables as a normalized Env map or an error if
// the file cannot be read or parsed.
func FromDotEnv(path string) (Env, error) {
	env, err := godotenv.Read(path)
	if err != nil {
		return nil, fmt.Errorf("loading dotenv: %w", err)
	}

	return Env(env).Normalized(), nil
}
