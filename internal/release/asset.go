package release

import (
	"fmt"
	"path/filepath"
	"strings"
)

// Asset represents a release asset with its name, download URL, content type, and optional digest.
type Asset struct {
	// Name is the release asset filename.
	Name string
	// URL is the asset download location.
	URL string
	// Type is the asset media type.
	Type string
	// Digest is the optional publisher-provided checksum.
	Digest string
}

// Match checks if the asset name matches the given pattern.
func (a Asset) Match(pattern string) (bool, error) {
	match, err := filepath.Match(pattern, a.Name)
	if err != nil {
		return false, fmt.Errorf("matching pattern: %w", err)
	}

	return match, nil
}

// HasExtension checks if the asset has the given file extension.
func (a Asset) HasExtension(extension string) (bool, error) {
	if strings.Count(extension, ".") <= 1 {
		return filepath.Ext(a.Name) == extension, nil
	}

	return strings.HasSuffix(a.Name, extension), nil
}
