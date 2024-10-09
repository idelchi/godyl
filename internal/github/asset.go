package github

import (
	"path/filepath"
	"strings"
)

// Asset represents a GitHub release asset.
type Asset struct {
	Name string `json:"name"`
	URL  string `json:"browser_download_url"`
	Type string `json:"content_type"`
}

// Match checks if the asset name matches the given pattern.
func (a Asset) Match(pattern string) (bool, error) {
	return filepath.Match(pattern, a.Name)
}

// HasExtension checks if the asset has the given extension.
func (a Asset) HasExtension(extension string) (bool, error) {
	if strings.Count(extension, ".") <= 1 {
		return filepath.Ext(a.Name) == extension, nil
	}

	return strings.HasSuffix(a.Name, extension), nil
}
