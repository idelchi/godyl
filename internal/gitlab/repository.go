package gitlab

import (
	"context"
	"fmt"
	"sort"

	gitlab "gitlab.com/gitlab-org/api/client-go"
)

// Repository represents a GitLab repository with its owner and name.
// It contains a GitLab client and context for making API calls.
type Repository struct {
	Owner  string          // Owner is the owner of the repository (GitLab username or group).
	Repo   string          // Repo is the name of the repository.
	client *gitlab.Client  // client is the GitLab client used to interact with the GitLab API.
	ctx    context.Context // ctx is the context used for API requests.
}

// NewRepository creates a new instance of Repository.
// It requires the repository owner, repository name, and a GitLab client.
func NewRepository(owner, repo string, client *gitlab.Client) *Repository {
	return &Repository{
		Owner:  owner,
		Repo:   repo,
		client: client,
		ctx:    context.Background(),
	}
}

// WithContext returns a copy of the repository with the given context.
func (g *Repository) WithContext(ctx context.Context) *Repository {
	repo := *g
	repo.ctx = ctx

	return &repo
}

// LatestRelease retrieves the latest release for the repository.
func (g *Repository) LatestRelease() (*Release, error) {
	path := fmt.Sprintf("%s/%s", g.Owner, g.Repo)

	releases, _, err := g.client.Releases.ListReleases(path, &gitlab.ListReleasesOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get list releases: %w", err)
	}

	if len(releases) == 0 {
		return nil, fmt.Errorf("no releases found for %s/%s", g.Owner, g.Repo)
	}

	// Get the first release (should be the latest)
	latestRelease := releases[0]

	release := &Release{}
	if err := release.FromRepositoryRelease(latestRelease); err != nil {
		return nil, fmt.Errorf("failed to parse release: %w", err)
	}

	return release, nil
}

// GetRelease retrieves a specific release for the repository based on the provided tag.
func (g *Repository) GetRelease(tag string) (*Release, error) {
	path := fmt.Sprintf("%s/%s", g.Owner, g.Repo)
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

// Languages retrieves the programming languages used in the repository, sorted by usage in descending order.
func (g *Repository) Languages() ([]string, error) {
	path := fmt.Sprintf("%s/%s", g.Owner, g.Repo)
	languages, _, err := g.client.Projects.GetProjectLanguages(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get languages: %w", err)
	}

	// Create a slice of keys to sort
	keys := make([]string, 0, len(*languages))
	for k := range *languages {
		keys = append(keys, k)
	}

	// Sort the keys based on the values in descending order
	sort.Slice(keys, func(i, j int) bool {
		return (*languages)[keys[i]] > (*languages)[keys[j]]
	})

	return keys, nil
}

// GetLatestIncludingPreRelease retrieves the most recently published release for the repository,
// including pre-releases. This returns the newest release by published date, regardless of
// whether it's a regular release or pre-release.
func (g *Repository) GetLatestIncludingPreRelease() (*Release, error) {
	path := fmt.Sprintf("%s/%s", g.Owner, g.Repo)
	releases, _, err := g.client.Releases.ListReleases(path, &gitlab.ListReleasesOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list releases: %w", err)
	}

	if len(releases) == 0 {
		return nil, fmt.Errorf("no releases found for %s/%s", g.Owner, g.Repo)
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
