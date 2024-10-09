package github

import (
	"path/filepath"

	"github.com/idelchi/godyl/internal/detect/platform"
	"github.com/idelchi/godyl/internal/match"
	"github.com/idelchi/godyl/pkg/compare"
)

// Assets represents a collection of GitHub release assets.
type Assets []Asset

// FilterByName returns the assets that match the given name.
func (as Assets) FilterByName(name string) (assets Assets) {
	for _, asset := range as {
		if compare.Lower(asset.Name, name) {
			assets = append(assets, asset)
		}
	}

	return assets
}

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
func (as Assets) Match(requirements match.Requirements) (matches match.Results, err error) {
	var assets match.Assets

	for _, a := range as {
		asset := match.Asset{Name: a.Name}
		asset.Parse()
		asset.Platfrom.Extension = platform.Extension(filepath.Ext(a.Name))

		assets = append(assets, asset)
	}

	matches = assets.Select(requirements)
	return matches, matches.Status()
}
