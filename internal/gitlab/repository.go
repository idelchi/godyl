package gitlab

import (
	"context"
	"fmt"

	gitlab "gitlab.com/gitlab-org/api/client-go"
)

// Repository represents a GitLab repository with its owner and name.
// It contains a GitLab client for making API calls.
type Repository struct {
	client    *gitlab.Client
	Namespace string
	Repo      string
}

// NewRepository creates a new instance of Repository.
// It requires the repository owner, repository name, and a GitLab client.
func NewRepository(namespace, repo string, client *gitlab.Client) *Repository {
	return &Repository{
		Namespace: namespace,
		Repo:      repo,
		client:    client,
	}
}

// GetRelease retrieves a specific release for the repository based on the provided tag.
func (g *Repository) GetRelease(_ context.Context, tag string) (*Release, error) {
	path := fmt.Sprintf("%s/%s", g.Namespace, g.Repo)

	gitlabRelease, _, err := g.client.Releases.GetRelease(path, tag)
	if err != nil {
		return nil, fmt.Errorf("failed to get assets for release tag %q: %w", tag, err)
	}

	release := &Release{}
	if err := release.FromRepositoryRelease(gitlabRelease); err != nil {
		return nil, fmt.Errorf("failed to parse release: %w", err)
	}

	return release, nil
}

// LatestRelease retrieves the latest release for the repository.
func (g *Repository) LatestRelease(ctx context.Context) (*Release, error) {
	const PerPage = 1000

	releases, err := g.getReleasesWithOptions(ctx, PerPage)
	if err != nil {
		return nil, err
	}

	// Get the first release (should be the latest)
	latestRelease := releases[0]

	release := &Release{}
	if err := release.FromRepositoryRelease(latestRelease); err != nil {
		return nil, fmt.Errorf("failed to parse release: %w", err)
	}

	return release, nil
}

// GetLatestIncludingPreRelease retrieves the most recently published release for the repository,
// including pre-releases. This returns the newest release by published date, regardless of
// whether it's a regular release or pre-release.
func (g *Repository) GetLatestIncludingPreRelease(ctx context.Context, perPage int) (*Release, error) {
	releases, err := g.getReleasesWithOptions(ctx, perPage)
	if err != nil {
		return nil, err
	}

	// Find the most recent release by published date
	var latestRelease *gitlab.Release

	for i, release := range releases {
		if i == 0 || latestRelease.CreatedAt.Before(*release.CreatedAt) {
			latestRelease = release
		}
	}

	// Convert to our Release type
	release := &Release{}
	if err := release.FromRepositoryRelease(latestRelease); err != nil {
		return nil, fmt.Errorf("failed to parse release: %w", err)
	}

	return release, nil
}

// getReleasesWithOptions retrieves releases for the repository using the provided options.
func (g *Repository) getReleasesWithOptions(_ context.Context, perPage int) ([]*gitlab.Release, error) {
	path := fmt.Sprintf("%s/%s", g.Namespace, g.Repo)

	releases, _, err := g.client.Releases.ListReleases(
		path,
		&gitlab.ListReleasesOptions{ListOptions: gitlab.ListOptions{PerPage: int64(perPage)}},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list releases: %w", err)
	}

	if len(releases) == 0 {
		return nil, fmt.Errorf("no releases found for %s/%s", g.Namespace, g.Repo)
	}

	return releases, nil
}
