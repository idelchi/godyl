// Package app provides the entrypoint for the application, creating the root command and executing it.
package app

import (
	"embed"
	"fmt"

	"github.com/idelchi/godyl/internal/cli"
	"github.com/idelchi/godyl/internal/cli/common"
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
	files, err := common.NewEmbeddedFiles(a.embeds)
	if err != nil {
		return fmt.Errorf("creating embedded files: %w", err)
	}

	// Execute the application
	if err := cli.Command(cfg, files, a.version).Execute(); err != nil {
		return err //nolint:wrapcheck 	// Error does not need additional wrapping.
	}

	return nil
}
