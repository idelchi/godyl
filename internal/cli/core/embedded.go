package core

import (
	"embed"
	"fmt"
)

// Embedded holds the embedded files for the application.
type Embedded struct {
	// Defaults to be set for each tool (not flags).
	Defaults []byte
	// Default list of tools that can be used to either view or dump out.
	Tools []byte
	// A template for the cleanup script.
	Template []byte
}

// NewEmbeddedFiles loads embedded configuration files and templates.
func NewEmbeddedFiles(embeds embed.FS) (*Embedded, error) {
	files := &Embedded{}

	var err error

	// Read embedded files
	if files.Defaults, err = embeds.ReadFile("defaults.yml"); err != nil {
		return files, fmt.Errorf("reading defaults file: %w", err)
	}

	if files.Tools, err = embeds.ReadFile("tools.yml"); err != nil {
		return files, fmt.Errorf("reading tools file: %w", err)
	}

	if files.Template, err = embeds.ReadFile("internal/updater/scripts/cleanup.bat.template"); err != nil {
		return files, fmt.Errorf("reading cleanup template: %w", err)
	}

	return files, nil
}
