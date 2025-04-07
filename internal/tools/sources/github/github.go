package github

import (
	"errors"
	"fmt"

	"github.com/idelchi/godyl/internal/github"
	"github.com/idelchi/godyl/internal/match"
	"github.com/idelchi/godyl/internal/tools/sources/common"
	"github.com/idelchi/godyl/pkg/path/file"
)

// GitHub represents a GitHub repository with optional authentication token and metadata.
type GitHub struct {
	Repo  string
	Owner string
	Token string `mask:"fixed"`
	Pre   bool

	// Data holds additional metadata related to the repository.
	Data common.Metadata `yaml:"-"`

	latestStoredRelease *github.Release
}

// Get retrieves a specific attribute from the GitHub repository's metadata.
func (g *GitHub) Get(attribute string) string {
	return g.Data.Get(attribute)
}

// LatestVersion fetches the latest release version of the GitHub repository.
func (g *GitHub) LatestVersion() (string, error) {
	client := github.NewClient(g.Token)
	repository := github.NewRepository(g.Owner, g.Repo, client)

	var release *github.Release
	var err error

	if g.Pre {
		release, err = repository.GetLatestIncludingPreRelease()
	} else {
		release, err = repository.LatestRelease()
	}

	if err != nil {
		return "", fmt.Errorf("failed to retrieve latest release: %w", err)
	}

	// Store the latest release for future use
	g.latestStoredRelease = release

	return release.Tag, nil
}

// MatchAssetsToRequirements matches release assets to specific file extensions and requirements,
// returning the URL of the matched asset.
func (g *GitHub) MatchAssetsToRequirements(
	_ []string,
	version string,
	requirements match.Requirements,
) (string, error) {
	client := github.NewClient(g.Token)
	repository := github.NewRepository(g.Owner, g.Repo, client)

	var release *github.Release

	if g.latestStoredRelease == nil {
		var err error

		release, err = repository.GetRelease(version)
		if err != nil {
			return "", fmt.Errorf("failed to get release: %w", err)
		}
	} else {
		release = g.latestStoredRelease
	}

	assets := release.Assets

	matches := assets.Match(requirements)
	if matches.Status() != nil {
		return "", matches.WithoutZero().Status()
	}

	if len(matches) == 0 {
		return "", fmt.Errorf("no assets found for requirements: %v", requirements)
	}

	err := matches.Status()
	if err != nil {
		err = fmt.Errorf("status: %w", err)
	}

	return assets.FilterByName(matches[0].Asset.Name)[0].URL, err
}

// PopulateOwnerAndRepo sets the Owner and Repo fields based on the given name.
// If Owner and Repo are already set, this method does nothing.
func (g *GitHub) PopulateOwnerAndRepo(name string) (err error) {
	// If both Owner and Repo are already set, nothing to do
	if g.Owner != "" && g.Repo != "" {
		return nil
	}

	// If exactly one of Owner or Repo is set (but not both), that's invalid
	if (g.Owner == "") != (g.Repo == "") {
		return errors.New("Either both `owner` and `repo` must be set or `name` must be in the format `owner/repo`")
	}

	g.Owner, g.Repo, err = common.SplitName(name)
	if err != nil {
		return err
	}

	return nil
}

// Initialize populates the GitHub repository's owner and name from the given input.
func (g *GitHub) Initialize(name string) error {
	if err := g.PopulateOwnerAndRepo(name); err != nil {
		return err
	}

	return nil
}

// Exe sets the executable name in the metadata to the repository name.
func (g *GitHub) Exe() error {
	g.Data.Set("exe", g.Repo)

	return nil
}

// Version fetches and sets the latest release version in the metadata.
func (g *GitHub) Version(_ string) error {
	version, err := g.LatestVersion()
	if err != nil {
		return err
	}

	g.Data.Set("version", version)

	return nil
}

// Path sets the download URL of the matched asset in the metadata, based on version, file extensions, and requirements.
func (g *GitHub) Path(_ string, extensions []string, version string, requirements match.Requirements) error {
	url, err := g.MatchAssetsToRequirements(extensions, version, requirements)
	if err != nil {
		return err
	}

	g.Data.Set("path", url)

	return nil
}

// Install downloads the asset from GitHub and returns the output, the found file, and any error encountered.
func (g *GitHub) Install(d common.InstallData) (output string, found file.File, err error) {
	return common.Download(d)
}
