// Package keyring is a simple wrapper that adds timeouts to the zalando/go-keyring package.
package keyring

import (
	"errors"
	"fmt"
	"time"

	"github.com/zalando/go-keyring"
)

var (
	ErrNotFound = errors.New("secret not found in keyring")
	ErrTimeout  = errors.New("timeout while accessing keyring")
)

// Set secret in keyring for user.
func Set(service, user, secret string, timeout time.Duration) error {
	ch := make(chan error, 1)
	go func() {
		defer close(ch)
		ch <- keyring.Set(service, user, secret)
	}()
	select {
	case err := <-ch:
		return err
	case <-time.After(timeout):
		return fmt.Errorf("setting secret: %w", ErrTimeout)
	}
}

// Get secret from keyring given service and user name.
func Get(service, user string, timeout time.Duration) (string, error) {
	ch := make(chan struct {
		val string
		err error
	}, 1)
	go func() {
		defer close(ch)

		val, err := keyring.Get(service, user)
		ch <- struct {
			val string
			err error
		}{val, err}
	}()
	select {
	case res := <-ch:
		if errors.Is(res.err, keyring.ErrNotFound) {
			return "", ErrNotFound
		}

		return res.val, res.err
	case <-time.After(timeout):
		return "", fmt.Errorf("getting secret: %w", ErrTimeout)
	}
}

// Delete secret from keyring.
func Delete(service, user string, timeout time.Duration) error {
	ch := make(chan error, 1)
	go func() {
		defer close(ch)
		ch <- keyring.Delete(service, user)
	}()
	select {
	case err := <-ch:
		return err
	case <-time.After(timeout):
		return fmt.Errorf("deleting secret: %w", ErrTimeout)
	}
}

func DeleteAll(service string, timeout time.Duration) error {
	ch := make(chan error, 1)
	go func() {
		defer close(ch)
		ch <- keyring.DeleteAll(service)
	}()
	select {
	case err := <-ch:
		return err
	case <-time.After(timeout):
		return fmt.Errorf("deleting all secrets: %w", ErrTimeout)
	}
}
