package github

import (
	"encoding/json"
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
	assets := release.Assets

	var releaseAssets Assets
	assetJSON, err := json.Marshal(assets)
	if err != nil {
		return fmt.Errorf("failed to marshal assets: %w", err)
	}

	if err := json.Unmarshal(assetJSON, &releaseAssets); err != nil {
		return fmt.Errorf("failed to unmarshal assets: %w", err)
	}

	var name string
	if release.Name != nil {
		name = *release.Name
	}

	*r = Release{
		Name:   name,
		Tag:    *release.TagName,
		Assets: releaseAssets,
	}

	return nil
}
