package github

import (
	"context"
	"fmt"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/google/go-github/v74/github"
)

// Repository represents a GitHub repository with its owner and name.
// It contains a GitHub client for making API calls.
type Repository struct {
	client *github.Client
	Owner  string
	Repo   string
}

// NewRepository creates a new instance of Repository.
// It requires the repository owner, repository name, and a GitHub client.
func NewRepository(owner, repo string, client *github.Client) *Repository {
	return &Repository{
		Owner:  owner,
		Repo:   repo,
		client: client,
	}
}

// LatestRelease retrieves the latest release for the repository.
func (r *Repository) LatestRelease(ctx context.Context) (*Release, error) {
	repositoryRelease, _, err := r.client.Repositories.GetLatestRelease(ctx, r.Owner, r.Repo)
	if err != nil {
		return nil, fmt.Errorf("getting latest release: %w", err)
	}

	release := &Release{}
	if err := release.FromRepositoryRelease(repositoryRelease); err != nil {
		return nil, fmt.Errorf("parsing release: %w", err)
	}

	return release, nil
}

// GetRelease retrieves a specific release for the repository based on the provided tag.
func (r *Repository) GetRelease(ctx context.Context, tag string) (*Release, error) {
	repositoryRelease, _, err := r.client.Repositories.GetReleaseByTag(ctx, r.Owner, r.Repo, tag)
	if err != nil {
		return nil, fmt.Errorf("getting assets for release tag %q: %w", tag, err)
	}

	release := &Release{}
	if err := release.FromRepositoryRelease(repositoryRelease); err != nil {
		return nil, fmt.Errorf("parsing release: %w", err)
	}

	return release, nil
}

// LatestIncludingPreRelease retrieves the most recently published release for the repository,
// including pre-releases. This returns the newest release by published date, regardless of
// whether it's a regular release or pre-release.
func (r *Repository) LatestIncludingPreRelease(ctx context.Context, perPage int) (*Release, error) {
	var allReleases []*github.RepositoryRelease

	page := 1

	for {
		opts := &github.ListOptions{
			PerPage: perPage,
			Page:    page,
		}

		releases, resp, err := r.client.Repositories.ListReleases(ctx, r.Owner, r.Repo, opts)
		if err != nil {
			return nil, fmt.Errorf("listing releases (page %d): %w", page, err)
		}

		allReleases = append(allReleases, releases...)

		// Check if there are more pages
		if resp.NextPage == 0 {
			break
		}

		page = resp.NextPage
	}

	if len(allReleases) == 0 {
		return nil, fmt.Errorf("no releases found for %s/%s", r.Owner, r.Repo)
	}

	// Find the most recent release by published date
	var latestRelease *github.RepositoryRelease

	for i, release := range allReleases {
		if i == 0 || release.PublishedAt == nil || latestRelease.PublishedAt == nil {
			latestRelease = release

			continue
		}

		// Compare the timestamps - need to use the Time property of Timestamp
		if release.PublishedAt.After(latestRelease.PublishedAt.Time) {
			latestRelease = release
		}
	}

	// Convert to our Release type
	release := &Release{}
	if err := release.FromRepositoryRelease(latestRelease); err != nil {
		return nil, fmt.Errorf("parsing release: %w", err)
	}

	return release, nil
}

// GetReleasesByWildcard retrieves the latest release matching a wildcard pattern.
// It returns the highest version that matches the pattern.
func (r *Repository) GetReleasesByWildcard(ctx context.Context, pattern string, perPage int) (*Release, error) {
	pattern = strings.ReplaceAll(pattern, "*", "X")

	c, err := semver.NewConstraint(pattern)
	if err != nil {
		return nil, fmt.Errorf("invalid version pattern %q: %w", pattern, err)
	}

	var allReleases []*github.RepositoryRelease

	page := 1

	for {
		opts := &github.ListOptions{
			PerPage: perPage,
			Page:    page,
		}

		releases, resp, err := r.client.Repositories.ListReleases(ctx, r.Owner, r.Repo, opts)
		if err != nil {
			return nil, fmt.Errorf("listing releases (page %d): %w", page, err)
		}

		allReleases = append(allReleases, releases...)

		// Check if there are more pages
		if resp.NextPage == 0 {
			break
		}

		page = resp.NextPage
	}

	if len(allReleases) == 0 {
		return nil, fmt.Errorf("no releases found for %s/%s", r.Owner, r.Repo)
	}

	var (
		highestVersion *semver.Version
		highestRelease *github.RepositoryRelease
	)

	for _, release := range allReleases {
		if release.TagName == nil {
			continue
		}

		// Parse version (handles v prefix automatically)
		v, err := semver.NewVersion(*release.TagName)
		if err != nil {
			continue // Skip non-semver tags
		}

		// Check if version matches constraint
		if !c.Check(v) {
			continue
		}

		// Track the highest matching version
		if highestVersion == nil || v.GreaterThan(highestVersion) {
			highestVersion = v
			highestRelease = release
		}
	}

	if highestRelease == nil {
		return nil, fmt.Errorf("no releases match pattern %q", pattern)
	}

	// Convert to our Release type
	release := &Release{}
	if err := release.FromRepositoryRelease(highestRelease); err != nil {
		return nil, fmt.Errorf("parsing release: %w", err)
	}

	return release, nil
}
