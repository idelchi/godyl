// Package ierrors provides error types for internal use.
package ierrors

import (
	"errors"
)

// ErrUsage is returned when there is an error in the configuration.
var ErrUsage = errors.New("usage error")

// ErrChangesPending indicates that status found tools that need installation.
var ErrChangesPending = errors.New("changes pending")
