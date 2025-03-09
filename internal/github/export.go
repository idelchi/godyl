package github

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// ExportConfig holds configuration for exporting GitHub data.
type ExportConfig struct {
	ExportPath string // Path to export data to
	ImportPath string // Path to import data from
}

// DefaultExportConfig returns the default export configuration.
func DefaultExportConfig() ExportConfig {
	return ExportConfig{
		ExportPath: "tests/assets2.json",
		ImportPath: "tests/assets.json",
	}
}

// exportMutex is a mutex to prevent concurrent writes to the export file.
var exportMutex sync.Mutex //nolint:gochecknoglobals

// Export retrieves the latest release for the repository and stores its assets in a JSON file.
func (g *Repository) Export(release *Release, config ExportConfig) error {
	exportMutex.Lock()
	defer exportMutex.Unlock()

	const permsDir = 0o750

	// Ensure the directory exists
	if err := os.MkdirAll(filepath.Dir(config.ExportPath), permsDir); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Load existing data
	var data map[string][]Asset
	if fileData, err := os.ReadFile(config.ExportPath); err == nil {
		if err := json.Unmarshal(fileData, &data); err != nil {
			return fmt.Errorf("failed to unmarshal existing data: %w", err)
		}
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("failed to read existing file: %w", err)
	}

	if data == nil {
		data = make(map[string][]Asset)
	}

	// Add or update the entry
	key := fmt.Sprintf("%s/%s", g.Owner, g.Repo)
	data[key] = release.Assets

	// Save the updated data
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	const permsFile = 0o600

	if err := os.WriteFile(config.ExportPath, jsonData, permsFile); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// ExportWithDefaults exports the release with default configuration.
func (g *Repository) ExportWithDefaults(release *Release) error {
	return g.Export(release, DefaultExportConfig())
}

// LatestReleaseFromExport retrieves the latest release information from the exported JSON file.
func (g *Repository) LatestReleaseFromExport(config ExportConfig) (*Release, error) {
	// Read the exported file
	fileData, err := os.ReadFile(config.ImportPath)
	if err != nil {
		return g.LatestRelease()
	}

	var data map[string][]Asset
	if err := json.Unmarshal(fileData, &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal export data: %w", err)
	}

	key := fmt.Sprintf("%s/%s", g.Owner, g.Repo)
	assets, ok := data[key]

	if !ok {
		return nil, fmt.Errorf("%w: no data found for repository %s", ErrExporter, key)
	}

	release := &Release{
		Assets: assets,
	}

	return release, nil
}

// LatestReleaseFromExportWithDefaults retrieves the latest release with default configuration.
func (g *Repository) LatestReleaseFromExportWithDefaults() (*Release, error) {
	return g.LatestReleaseFromExport(DefaultExportConfig())
}

// ErrExporter is an error returned when an exporter operation fails.
var ErrExporter = errors.New("exporter error")
