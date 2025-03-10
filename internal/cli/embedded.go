package cli

import (
	"embed"
	"fmt"

	"github.com/idelchi/godyl/internal/config"
)

// NewEmbeddedFiles loads embedded configuration files and templates.
func NewEmbeddedFiles(embeds embed.FS) (config.Embedded, error) {
	files := config.Embedded{}
	var err error

	// Read embedded files
	if files.Defaults, err = embeds.ReadFile("defaults.yml"); err != nil {
		return files, fmt.Errorf("reading defaults file: %w", err)
	}

	if files.Tools, err = embeds.ReadFile("tools.yml"); err != nil {
		return files, fmt.Errorf("reading tools file: %w", err)
	}

	if files.Template, err = embeds.ReadFile("internal/core/updater/scripts/cleanup.bat.template"); err != nil {
		return files, fmt.Errorf("reading cleanup template: %w", err)
	}

	return files, nil
}
