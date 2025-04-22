package github

import (
	"errors"
	"fmt"

	"github.com/hashicorp/go-getter/v2"

	"github.com/idelchi/godyl/internal/github"
	"github.com/idelchi/godyl/internal/match"
	"github.com/idelchi/godyl/internal/tools/sources/common"
	"github.com/idelchi/godyl/pkg/path/file"
)

// GitHub represents a GitHub repository configuration and state.
type GitHub struct {
	Data                common.Metadata `yaml:"-"`
	latestStoredRelease *github.Release
	Repo                string
	Owner               string
	Token               string `mask:"fixed"`
	Pre                 bool
}

// Initialize sets up the GitHub repository configuration from the given name.
// Returns an error if the repository name format is invalid.
func (g *GitHub) Initialize(name string) error {
	if err := g.PopulateOwnerAndRepo(name); err != nil {
		return err
	}

	g.Data.Set("exe", g.Repo)

	return nil
}

// Version fetches the latest release version and stores it in metadata.
func (g *GitHub) Version(_ string) error {
	version, err := g.LatestVersion()
	if err != nil {
		return err
	}

	g.Data.Set("version", version)

	return nil
}

// Path finds a matching release asset and stores its URL in metadata.
// Uses version, extensions, and requirements to find the appropriate asset.
func (g *GitHub) Path(_ string, extensions []string, version string, requirements match.Requirements) error {
	url, err := g.MatchAssetsToRequirements(extensions, version, requirements)
	if err != nil {
		return err
	}

	g.Data.Set("path", url)

	return nil
}

// Install downloads the GitHub release asset using the provided configuration.
// Returns the operation output, downloaded file information, and any errors.
func (g *GitHub) Install(
	d common.InstallData,
	progressListener getter.ProgressTracker,
) (output string, found file.File, err error) {
	// Pass the progress listener down to the common download function
	d.ProgressListener = progressListener

	return common.Download(d)
}

// Get retrieves a metadata attribute value by its key.
func (g *GitHub) Get(attribute string) string {
	return g.Data.Get(attribute)
}

// LatestVersion fetches the latest release version from GitHub.
// Returns the tag name of the latest release, respecting the Pre flag setting.
func (g *GitHub) LatestVersion() (string, error) {
	client := github.NewClient(g.Token)
	repository := github.NewRepository(g.Owner, g.Repo, client)

	var release *github.Release

	var err error

	if g.Pre {
		release, err = repository.GetLatestIncludingPreRelease()
	} else {
		if tag, err := repository.LatestVersionFromWeb(); err == nil {
			return tag, nil
		}

		release, err = repository.LatestRelease()
	}

	if err != nil {
		return "", fmt.Errorf("failed to retrieve latest release: %w", err)
	}

	// Store the latest release for future use
	g.latestStoredRelease = release

	return release.Tag, nil
}

// MatchAssetsToRequirements finds release assets matching the given requirements.
// Returns the download URL of the best matching asset, considering platform,
// architecture, and other specified requirements.
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

// PopulateOwnerAndRepo sets the Owner and Repo fields from a name string.
// Expects name in "owner/repo" format if fields are not already set.
// Returns an error if the format is invalid or fields are partially set.
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
