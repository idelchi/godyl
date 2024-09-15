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
func (ts TokenStore) Available() (error, bool) {
	// Try a quick get on a non-existent key
	_, err := keyring.Get(ts.Service, "__health_check__", 1*time.Second)

	return err, !errors.Is(err, keyring.ErrTimeout)
}

// GetAll retrieves all tokens for the specified users.
func (ts TokenStore) GetAll(users []string) (map[string]string, error) {
	tokens := make(map[string]string)

	for _, user := range users {
		value, err := ts.Get(user)
		switch {
		case err == nil:
		case errors.Is(err, keyring.ErrNotFound):
			continue
		default:
			return tokens, err
		}

		tokens[user] = value
	}

	return tokens, nil
}

// SetAll sets multiple tokens for the specified users.
func (ts TokenStore) SetAll(tokens map[string]any) error {
	for user, token := range tokens {
		if err := ts.Set(user, token.(string)); err != nil {
			return fmt.Errorf("%s: failed to set token: %w", user, err)
		}
	}

	return nil
}

// Set stores a token for a specific user in the keyring.
func (ts TokenStore) Set(user, token string) error {
	return keyring.Set(ts.Service, user, token, defaultTimeout)
}

// Get retrieves a token for a specific user from the keyring.
func (ts TokenStore) Get(user string) (string, error) {
	token, err := keyring.Get(ts.Service, user, defaultTimeout)
	if err != nil {
		return "", err
	}

	return token, nil
}

// Delete removes the given keys from the keyring.
// If no keys are provided, it clears all tokens for the service.
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
