package github

import (
	"path/filepath"
	"strings"

	"github.com/idelchi/godyl/internal/debug"
	"github.com/idelchi/godyl/internal/detect/platform"
	"github.com/idelchi/godyl/internal/match"
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

// Checksum returns the first asset that appears to be a checksum file, or nil if none found.
func (as Assets) Checksum() *Asset {
	for _, asset := range as {
		if asset.IsChecksumLike() {
			debug.Debug("found checksum asset: %q", asset.Name)

			return &asset
		}
	}

	return nil
}
