package github

import (
	"errors"
	"fmt"

	"github.com/google/go-github/v74/github"

	"github.com/idelchi/godyl/pkg/generic"
)

// ErrRelease is returned when a release issue is encountered.
var ErrRelease = errors.New("release")

// Release represents a GitHub release, containing the release name, tag, and associated assets.
type Release struct {
	Name   string `json:"name"`
	Tag    string `json:"tag_name"`
	Body   string `json:"body"`
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

	assets := make(Assets, 0, len(release.Assets))

	for _, asset := range release.Assets {
		if asset.Name == nil || asset.BrowserDownloadURL == nil || asset.ContentType == nil {
			continue // Skip assets with missing required fields
		}

		assets = append(assets, Asset{
			Name:   generic.SafeDereference(asset.Name),
			URL:    generic.SafeDereference(asset.BrowserDownloadURL),
			Type:   generic.SafeDereference(asset.ContentType),
			Digest: generic.SafeDereference(asset.Digest),
		})
	}

	*r = Release{
		Name:   generic.SafeDereference(release.Name),
		Tag:    generic.SafeDereference(release.TagName),
		Assets: assets,
		Body:   generic.SafeDereference(release.Body),
	}

	return nil
}
