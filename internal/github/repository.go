package github

import (
	"context"
	"fmt"
	"sort"

	"github.com/google/go-github/v64/github"
)

// Repository represents a GitHub repository with its owner and name.
// It contains a GitHub client and context for making API calls.
type Repository struct {
	Owner  string          // Owner is the owner of the repository (GitHub username or organization).
	Repo   string          // Repo is the name of the repository.
	client *github.Client  // client is the GitHub client used to interact with the GitHub API.
	ctx    context.Context // ctx is the context used for API requests.
}

// NewRepository creates a new instance of Repository.
// It requires the repository owner, repository name, and a GitHub client.
func NewRepository(owner, repo string, client *github.Client) *Repository {
	ctx := context.Background()

	return &Repository{
		Owner:  owner,
		Repo:   repo,
		client: client,
		ctx:    ctx,
	}
}

// LatestRelease retrieves the latest release for the repository.
func (g *Repository) LatestRelease() (*Release, error) {
	ctx := context.TODO()

	repositoryRelease, _, err := g.client.Repositories.GetLatestRelease(ctx, g.Owner, g.Repo)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest release: %w", err)
	}

	release := &Release{}

	if err := release.FromRepositoryRelease(repositoryRelease); err != nil {
		return nil, fmt.Errorf("failed to parse release: %w", err)
	}

	return release, nil
}

// GetRelease retrieves a specific release for the repository based on the provided tag.
func (g *Repository) GetRelease(tag string) (*Release, error) {
	ctx := context.TODO()

	repositoryRelease, _, err := g.client.Repositories.GetReleaseByTag(ctx, g.Owner, g.Repo, tag)
	if err != nil {
		return nil, fmt.Errorf("failed to get assets for release tag %q: %w", tag, err)
	}

	release := &Release{}

	if err := release.FromRepositoryRelease(repositoryRelease); err != nil {
		return nil, fmt.Errorf("failed to parse release: %w", err)
	}

	return release, nil
}

// Languages retrieves the programming languages used in the repository, sorted by usage in descending order.
func (g *Repository) Languages() ([]string, error) {
	ctx := context.TODO()

	languages, _, err := g.client.Repositories.ListLanguages(ctx, g.Owner, g.Repo)
	if err != nil {
		return nil, fmt.Errorf("failed to get languages: %w", err)
	}

	// Create a slice of keys to sort
	keys := make([]string, 0, len(languages))
	for k := range languages {
		keys = append(keys, k)
	}

	// Sort the keys based on the values in descending order
	sort.Slice(keys, func(i, j int) bool {
		return languages[keys[i]] > languages[keys[j]]
	})

	return keys, nil
}
