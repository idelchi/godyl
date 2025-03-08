package github

import (
	"errors"
	"fmt"

	"github.com/google/go-github/v64/github"
)

// ErrRelease is returned when a release issue is encountered.
var ErrRelease = errors.New("release")

// Release represents a GitHub release, containing the release name, tag, and associated assets.
type Release struct {
	// Name is the name of the release.
	Name string `json:"name"`
	// Tag is the tag associated with the release (e.g., version number).
	Tag string `json:"tag_name"` //nolint:tagliatelle
	// Assets is a collection of assets attached to the release.
	Assets Assets `json:"assets"`
}

// FromRepositoryRelease converts a GitHub repository release to a Release object.
func (r *Release) FromRepositoryRelease(release *github.RepositoryRelease) error {
	if release == nil {
		return fmt.Errorf("%w: repository release is nil", ErrRelease)
	}

	if release.TagName == nil {
		return fmt.Errorf("%w: release tag name is nil", ErrRelease)
	}

	// Convert GitHub assets to our Asset type
	assets := make(Assets, 0, len(release.Assets))

	for _, asset := range release.Assets {
		if asset.Name == nil || asset.BrowserDownloadURL == nil || asset.ContentType == nil {
			continue // Skip assets with missing required fields
		}

		assets = append(assets, Asset{
			Name: *asset.Name,
			URL:  *asset.BrowserDownloadURL,
			Type: *asset.ContentType,
		})
	}

	// Get release name, defaulting to empty string if nil
	var name string
	if release.Name != nil {
		name = *release.Name
	}

	*r = Release{
		Name:   name,
		Tag:    *release.TagName,
		Assets: assets,
	}

	return nil
}
