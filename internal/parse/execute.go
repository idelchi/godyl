// Package parse generates and executes the command-line interface for the application.
package parse

import (
	"embed"
	"errors"
	"fmt"

	"github.com/idelchi/godyl/internal/app"
	"github.com/idelchi/gogen/pkg/cobraext"
)

// Execute creates and configures the command-line interface.
// It runs the root command with all subcommands and flags configured.
// This function is kept for backward compatibility and delegates to the app package.
func Execute(version string, embeds embed.FS) error {
	application := app.New(version, embeds)

	// Execute the application
	switch err := application.Execute(); {
	case errors.Is(err, cobraext.ErrExitGracefully):
		return nil
	case err != nil:
		return fmt.Errorf("executing command: %w", err)
	default:
		return nil
	}
}
