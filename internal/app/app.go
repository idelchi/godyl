// Package app provides the entrypoint for the application, creating the root command and executing it.
package app

import (
	"embed"

	"github.com/idelchi/godyl/internal/cli"
	"github.com/idelchi/godyl/internal/cli/core"
)

// Execute runs the root command.
func Execute(version string, files embed.FS) error {
	// Get the embedded files
	embedded, err := core.NewEmbeddedFiles(files)
	if err != nil {
		return err
	}

	// Execute the application
	if err := cli.Command(embedded, version).Execute(); err != nil {
		return err
	}

	return nil
}
