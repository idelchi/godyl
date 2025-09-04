// Package tokenstore provides secure token storage and retrieval functionality.
package tokenstore

import (
	"errors"
	"fmt"
	"time"

	"github.com/idelchi/godyl/internal/tokenstore/keyring"
)

// service is the name of the keyring service used to store tokens.
const service = "godyl"

const defaultTimeout = 3 * time.Second

// TokenStore provides methods to manage authentication tokens using the keyring package.
type TokenStore struct {
	Service string
}

// New creates a new TokenStore instance with the default service name.
func New() TokenStore {
	return TokenStore{
		Service: service,
	}
}

// Available checks if the keyring service is available for storing tokens.
func (ts TokenStore) Available() (bool, error) {
	// Try a quick get on a non-existent key
	_, err := keyring.Get(ts.Service, "__health_check__", 1*time.Second)

	return !errors.Is(err, keyring.ErrTimeout), err
}

// GetAll retrieves a map of tokens from the keyring for the specified keys.
// If a key is not found, it is skipped.
func (ts TokenStore) GetAll(keys ...string) (map[string]string, error) {
	tokens := make(map[string]string)

	for _, key := range keys {
		value, err := ts.Get(key)

		switch {
		case errors.Is(err, keyring.ErrNotFound):
			continue
		case err != nil:
			return tokens, err
		}

		tokens[key] = value
	}

	return tokens, nil
}

// SetAll sets multiple tokens in the keyring.
func (ts TokenStore) SetAll(tokens map[string]any) error {
	for key, token := range tokens {
		value, ok := token.(string)
		if !ok {
			return fmt.Errorf("%s: invalid token type %T, expected string", key, token)
		}

		if err := ts.Set(key, value); err != nil {
			return fmt.Errorf("%s: failed to set key: %w", key, err)
		}
	}

	return nil
}

// Set stores a token in the keyring for a specific key.
func (ts TokenStore) Set(key, token string) error {
	return keyring.Set(ts.Service, key, token, defaultTimeout)
}

// Get retrieves a token from the keyring for a specific key.
func (ts TokenStore) Get(key string) (string, error) {
	token, err := keyring.Get(ts.Service, key, defaultTimeout)
	if err != nil {
		return "", err
	}

	return token, nil
}

// Delete removes the given keys from the keyring.
// If no keys are provided, it clears all keys for the service.
func (ts TokenStore) Delete(keys ...string) error {
	if len(keys) == 0 {
		return keyring.DeleteAll(ts.Service, defaultTimeout)
	}

	var errs []error

	for _, key := range keys {
		if err := keyring.Delete(ts.Service, key, defaultTimeout); err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", key, err))
		}
	}

	return errors.Join(errs...)
}
