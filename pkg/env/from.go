package env

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// FromEnv returns the current environment variables as an Env.
// It uses os.Environ to fetch all the environment variables and normalizes them before returning.
func FromEnv() Env {
	env, _ := FromSlice(os.Environ()...)

	return env.Normalized()
}

// FromSlice constructs an Env from a slice of `key=value` strings.
// It returns an error if any string in the slice is malformed.
func FromSlice(slice ...string) (Env, error) {
	e := make(Env, len(slice))

	for _, v := range slice {
		if err := e.Add(v); err != nil {
			return nil, err
		}
	}

	return e.Normalized(), nil
}

// FromDotEnv loads environment variables from a .env file specified by the path.
// It returns an error if there is an issue reading the file or processing its contents.
func FromDotEnv(path string) (Env, error) {
	env, err := godotenv.Read(path)
	if err != nil {
		return nil, fmt.Errorf("loading dotenv: %w", err)
	}

	return Env(env).Normalized(), nil
}
