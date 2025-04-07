package gitlab

import (
	"fmt"
	"path/filepath"
	"strings"
)

// Asset represents a GitLab release asset with its name, download URL, and content type.
type Asset struct {
	// Name is the name of the asset.
	Name string `json:"name"`
	// URL is the browser download URL for the asset.
	URL string `json:"url"` //nolint:tagliatelle
	// Type is the content type of the asset.
	Type string `json:"content_type"` //nolint:tagliatelle
}

// Match checks if the asset name matches the given pattern.
func (a Asset) Match(pattern string) (bool, error) {
	match, err := filepath.Match(pattern, a.Name)
	if err != nil {
		return false, fmt.Errorf("failed to match pattern: %w", err)
	}

	return match, nil
}

// HasExtension checks if the asset has the given file extension.
func (a Asset) HasExtension(extension string) (bool, error) {
	// If the extension contains one or fewer dots, check the file extension.
	if strings.Count(extension, ".") <= 1 {
		return filepath.Ext(a.Name) == extension, nil
	}

	// Otherwise, check if the name ends with the specified extension.
	return strings.HasSuffix(a.Name, extension), nil
}