// Package app provides the entrypoint for the application.
// It creates the root command and executes the application,
// handling or returning any errors that occur.
package app

import (
	"embed"
	"fmt"

	"github.com/idelchi/godyl/internal/cli"
	"github.com/idelchi/godyl/internal/config"
)

// Application hold the version string and embedded files.
type Application struct {
	embeds  embed.FS
	version string
}

// New creates a new Application instance.
func New(version string, embeds embed.FS) *Application {
	return &Application{
		version: version,
		embeds:  embeds,
	}
}

// Execute runs the root command.
func (a *Application) Execute() error {
	cfg := &config.Config{}

	root, err := cli.NewCommand(cfg, a.version, a.embeds)
	if err != nil {
		return fmt.Errorf("application failed to initialize: %w", err)
	}

	// Execute the application
	if err := root.Run(); err != nil {
		return err //nolint:wrapcheck 		// Wrapping here adds no value.
	}

	return nil
}
