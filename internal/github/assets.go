package github

import (
	"path/filepath"

	"github.com/idelchi/godyl/internal/detect/platform"
	"github.com/idelchi/godyl/internal/match"
	"github.com/idelchi/godyl/pkg/utils"
)

// Assets represents a collection of GitHub release assets.
type Assets []Asset

// FilterByName returns the assets that match the given name.
// It compares asset names in a case-insensitive manner.
func (as Assets) FilterByName(name string) (assets Assets) {
	for _, asset := range as {
		if utils.EqualLower(asset.Name, name) {
			assets = append(assets, asset)
		}
	}

	return assets
}

// FilterByExtension returns the assets that have one of the given extensions.
// If no extensions are provided, it returns all assets.
func (as Assets) FilterByExtension(extensions []string) (assets Assets) {
	if len(extensions) == 0 {
		return as
	}

	for _, asset := range as {
		for _, pattern := range extensions {
			match, err := asset.HasExtension(pattern)
			if err != nil {
				continue
			}
			if match {
				assets = append(assets, asset)
				break
			}
		}
	}

	return assets
}

// FilterByMatch returns the assets that match one of the given patterns.
// If no filters are provided, it returns all assets.
func (as Assets) FilterByMatch(filters []string) (assets Assets) {
	if len(filters) == 0 {
		return as
	}

	for _, asset := range as {
		for _, pattern := range filters {
			match, err := asset.Match(pattern)
			if err != nil {
				continue
			}
			if match {
				assets = append(assets, asset)
				break
			}
		}
	}

	return assets
}

// Match checks if the assets match the given requirements.
// It processes each asset to extract platform and extension information.
func (as Assets) Match(requirements match.Requirements) (matches match.Results, err error) {
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
	return matches, matches.Status()
}
