// Package parse generates and executes the command-line interface for the application.
package parse

import (
	"errors"
	"fmt"

	"github.com/idelchi/godyl/internal/app"
	"github.com/idelchi/godyl/internal/config"
	"github.com/idelchi/gogen/pkg/cobraext"
)

// Execute creates and configures the command-line interface.
// It runs the root command with all subcommands and flags configured.
// This function is kept for backward compatibility and delegates to the app package.
func Execute(version string, defaultsFile, toolsFile []byte, embeds interface{}) error {
	// Create a new application instance
	application := app.New(version, defaultsFile, toolsFile, embeds)

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

// For backward compatibility, we keep this function signature
// but implement it using the new app package
func createRootCommand(cfg *config.Config, version string, defaultsFile, toolsFile []byte, embeds interface{}) interface{ Execute() error } {
	return app.New(version, defaultsFile, toolsFile, embeds)
}
