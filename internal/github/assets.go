package github

import (
	"path"
	"path/filepath"
	"strings"

	"github.com/idelchi/godyl/internal/detect/platform"
	"github.com/idelchi/godyl/internal/match"
	"github.com/idelchi/godyl/internal/tools/checksum"
)

// Assets represents a collection of GitHub release assets.
type Assets []Asset

// FilterByName returns the assets that match the given name.
// It compares asset names in a case-insensitive manner.
func (as Assets) FilterByName(name string) (assets Assets) {
	for _, asset := range as {
		if strings.EqualFold(asset.Name, name) {
			assets = append(assets, asset)
		}
	}

	return assets
}

// Match checks if the assets match the given requirements.
// It processes each asset to extract platform and extension information.
func (as Assets) Match(requirements match.Requirements) (matches match.Results) {
	var assets match.Assets

	for _, a := range as {
		asset := match.Asset{Name: a.Name}
		asset.Parse() // Parse the asset name to extract additional info (platform, architecture, etc.)

		asset.Platform.Extension = platform.Extension(
			filepath.Ext(a.Name),
		) // Assign the file extension to the platform field

		assets = append(assets, asset)
	}

	// Select the assets that satisfy the given requirements.
	matches = assets.Select(requirements)

	return matches
}

// Checksums returns all assets that appear to be checksum files.
func (as Assets) Checksums(pattern string) checksum.Checksums {
	checksums := checksum.Checksums{}

	if pattern != "" {
		for _, checksum := range checksums {
			match, err := path.Match(pattern, checksum)
			if err == nil && match {
				checksums = append(checksums, checksum)
			}
		}

		return checksums
	}

	for _, asset := range as {
		checksums = append(checksums, asset.Name)
	}

	return checksums.IsChecksumLike()
}
