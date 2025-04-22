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

// func (e *Env) UnmarshalYAML(value *yaml.Node) error {
// 	// Log the values before unmarshalling
// 	fmt.Printf("Before unmarshalling: %+v\n", *e)

// 	// Create a temporary map to unmarshal into
// 	var temp map[string]string

// 	// Unmarshal the YAML node into the temporary map
// 	if err := value.Decode(&temp); err != nil {
// 		return err
// 	}

// 	// If e is nil, initialize it
// 	if *e == nil {
// 		*e = make(Env)
// 	}

// 	// Copy values from temp to e
// 	for k, v := range temp {
// 		(*e)[k] = v
// 	}

// 	// Log the values after unmarshalling
// 	fmt.Printf("After unmarshalling: %+v\n\n\n", *e)

// 	return nil
// }
