// Package app provides the main application functionality.
package app

import (
	"fmt"

	"github.com/idelchi/godyl/internal/cli"
	"github.com/idelchi/godyl/internal/config"
)

// Application represents the main application.
type Application struct {
	version      string
	defaultsFile []byte
	toolsFile    []byte
	embeds       interface{}
}

// New creates a new Application instance.
func New(version string, defaultsFile, toolsFile []byte, embeds interface{}) *Application {
	return &Application{
		version:      version,
		defaultsFile: defaultsFile,
		toolsFile:    toolsFile,
		embeds:       embeds,
	}
}

// Execute runs the application.
func (a *Application) Execute() error {
	// Create a new configuration
	cfg := &config.Config{}

	// Create and configure the root command
	root := cli.NewRootCmd(cfg, a.version, a.defaultsFile, a.toolsFile, a.embeds)

	// Execute the command
	if err := root.Execute(); err != nil {
		return fmt.Errorf("executing command: %w", err)
	}

	return nil
}
