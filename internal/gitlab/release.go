package gitlab

import (
	"errors"
	"fmt"

	gitlab "gitlab.com/gitlab-org/api/client-go"
)

// ErrRelease is returned when a release issue is encountered.
var ErrRelease = errors.New("release")

// Release represents a GitLab release, containing the release name, tag, and associated assets.
type Release struct {
	// Name is the name of the release.
	Name string `json:"name"`
	// Tag is the tag associated with the release (e.g., version number).
	Tag string `json:"tag_name"` //nolint:tagliatelle // GitLab API uses snake_case field names
	// Assets is a collection of assets attached to the release.
	Assets Assets `json:"assets"`
}

// FromRepositoryRelease converts a GitLab repository release to a Release object.
func (r *Release) FromRepositoryRelease(release *gitlab.Release) error {
	if release == nil {
		return fmt.Errorf("%w: repository release is nil", ErrRelease)
	}

	if release.TagName == "" {
		return fmt.Errorf("%w: release tag name is empty", ErrRelease)
	}

	// Convert GitLab assets to our Asset type
	assets := make(Assets, 0, len(release.Assets.Links))

	for _, link := range release.Assets.Links {
		assets = append(assets, Asset{
			Name: link.Name,
			URL:  link.URL,
			Type: string(link.LinkType), // Convert LinkTypeValue to string
		})
	}

	*r = Release{
		Name:   release.Name,
		Tag:    release.TagName,
		Assets: assets,
	}

	return nil
}
