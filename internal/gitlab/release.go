package gitlab

import (
	"fmt"

	"github.com/idelchi/godyl/internal/release"

	gitlab "gitlab.com/gitlab-org/api/client-go"
)

// FromRepositoryRelease converts a GitLab repository release to a shared Release object.
func FromRepositoryRelease(repoRelease *gitlab.Release) (*release.Release, error) {
	if repoRelease == nil {
		return nil, fmt.Errorf("%w: repository release is nil", release.ErrRelease)
	}

	if repoRelease.TagName == "" {
		return nil, fmt.Errorf("%w: release tag name is empty", release.ErrRelease)
	}

	// Convert GitLab assets to our Asset type
	assets := make(release.Assets, 0, len(repoRelease.Assets.Links))

	for _, link := range repoRelease.Assets.Links {
		assets = append(assets, release.Asset{
			Name: link.Name,
			URL:  link.DirectAssetURL,
			Type: string(link.LinkType), // Convert LinkTypeValue to string
		})
	}

	return &release.Release{
		Name:   repoRelease.Name,
		Tag:    repoRelease.TagName,
		Assets: assets,
	}, nil
}
