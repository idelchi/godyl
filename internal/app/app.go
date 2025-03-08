// Package app provides the main application functionality.
package app

import (
	"embed"
	"errors"
	"fmt"

	"github.com/idelchi/godyl/internal/cli"
	"github.com/idelchi/godyl/internal/config"
	"github.com/idelchi/gogen/pkg/cobraext"
)

// Application represents the main application,
// holding the version and embedded files.
type Application struct {
	version string
	embeds  embed.FS
}

// New creates a new Application instance.
func New(version string, embeds embed.FS) *Application {
	return &Application{
		version: version,
		embeds:  embeds,
	}
}

// Execute runs the application.
func (a *Application) Execute() error {
	cfg := &config.Config{}

	root, err := cli.NewRootCmd(cfg, a.version, a.embeds)
	if err != nil {
		return fmt.Errorf("creating root command: %w", err)
	}

	// Execute the application
	switch err := root.Execute(); {
	case errors.Is(err, cobraext.ErrExitGracefully):
		return nil
	case err != nil:
		return fmt.Errorf("executing command: %w", err)
	}

	return nil
}
