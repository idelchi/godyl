package github

import (
	"path/filepath"
	"strings"
)

// Asset represents a GitHub release asset with its name, download URL, and content type.
type Asset struct {
	Name string `json:"name"`                 // Name is the name of the asset.
	URL  string `json:"browser_download_url"` // URL is the browser download URL for the asset.
	Type string `json:"content_type"`         // Type is the content type of the asset.
}

// Match checks if the asset name matches the given pattern.
func (a Asset) Match(pattern string) (bool, error) {
	return filepath.Match(pattern, a.Name)
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
