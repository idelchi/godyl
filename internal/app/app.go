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

	// Get the embedded files
	files, err := config.NewEmbeddedFiles(a.embeds)
	if err != nil {
		return fmt.Errorf("creating embedded files: %w", err)
	}

	root := cli.Command(cfg, files, a.version)

	// Execute the application
	if err := root.Execute(); err != nil {
		return err //nolint:wrapcheck 		// Wrapping here adds no value.
	}

	return nil
}
