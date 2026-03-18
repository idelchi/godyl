package github

import (
	"fmt"

	"github.com/google/go-github/v74/github"

	"github.com/idelchi/godyl/internal/release"
	"github.com/idelchi/godyl/pkg/generic"
)

// FromRepositoryRelease converts a GitHub repository release to a shared Release object.
func FromRepositoryRelease(repoRelease *github.RepositoryRelease) (*release.Release, error) {
	if repoRelease == nil {
		return nil, fmt.Errorf("%w: repository release is nil", release.ErrRelease)
	}

	if repoRelease.TagName == nil {
		return nil, fmt.Errorf("%w: release tag name is nil", release.ErrRelease)
	}

	assets := make(release.Assets, 0, len(repoRelease.Assets))

	for _, asset := range repoRelease.Assets {
		if asset == nil || asset.Name == nil || asset.BrowserDownloadURL == nil || asset.ContentType == nil {
			continue // Skip assets with missing required fields
		}

		assets = append(assets, release.Asset{
			Name:   generic.SafeDereference(asset.Name),
			URL:    generic.SafeDereference(asset.BrowserDownloadURL),
			Type:   generic.SafeDereference(asset.ContentType),
			Digest: generic.SafeDereference(asset.Digest),
		})
	}

	return &release.Release{
		Name:   generic.SafeDereference(repoRelease.Name),
		Tag:    generic.SafeDereference(repoRelease.TagName),
		Assets: assets,
		Body:   generic.SafeDereference(repoRelease.Body),
	}, nil
}
