package github

import (
	"fmt"

	"github.com/google/go-github/v64/github"
)

// Release represents a GitHub release, containing the release name, tag, and associated assets.
type Release struct {
	Name   string `json:"name"`     // Name is the name of the release.
	Tag    string `json:"tag_name"` // Tag is the tag associated with the release (e.g., version number).
	Assets Assets `json:"assets"`   // Assets is a collection of assets attached to the release.
}

// FromRepositoryRelease converts a GitHub repository release to a Release object.
func (r *Release) FromRepositoryRelease(release *github.RepositoryRelease) error {
	if release == nil {
		return fmt.Errorf("repository release is nil")
	}

	if release.TagName == nil {
		return fmt.Errorf("release tag name is nil")
	}

	// Convert GitHub assets to our Asset type
	assets := make(Assets, 0, len(release.Assets))
	for _, a := range release.Assets {
		if a.Name == nil || a.BrowserDownloadURL == nil || a.ContentType == nil {
			continue // Skip assets with missing required fields
		}

		assets = append(assets, Asset{
			Name: *a.Name,
			URL:  *a.BrowserDownloadURL,
			Type: *a.ContentType,
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
