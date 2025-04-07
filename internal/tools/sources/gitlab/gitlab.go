package gitlab

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/idelchi/godyl/internal/gitlab"
	"github.com/idelchi/godyl/internal/match"
	"github.com/idelchi/godyl/internal/tools/sources/common"
	"github.com/idelchi/godyl/pkg/path/file"
)

// GitLab represents a GitLab repository with optional authentication token and metadata.
type GitLab struct {
	Repo   string
	Owner  string
	Token  string `mask:"fixed"`
	Pre    bool
	Server string

	// Data holds additional metadata related to the repository.
	Data common.Metadata `yaml:"-"`

	latestStoredRelease *gitlab.Release
}

// Get retrieves a specific attribute from the GitLab repository's metadata.
func (g *GitLab) Get(attribute string) string {
	return g.Data.Get(attribute)
}

// LatestVersion fetches the latest release version of the GitLab repository.
func (g *GitLab) LatestVersion() (string, error) {
	client, err := gitlab.NewClient(g.Token, g.Server)
	if err != nil {
		return "", fmt.Errorf("failed to create GitLab client: %w", err)
	}

	repository := gitlab.NewRepository(g.Owner, g.Repo, client)

	var release *gitlab.Release

	if g.Pre {
		release, err = repository.GetLatestIncludingPreRelease()
	} else {
		release, err = repository.LatestRelease()
	}

	if err != nil {
		return "", fmt.Errorf("failed to get latest release: %w", err)
	}

	// Store the latest release for future use
	g.latestStoredRelease = release

	return release.Tag, nil
}

// MatchAssetsToRequirements matches release assets to specific file extensions and requirements,
// returning the URL of the matched asset.
func (g *GitLab) MatchAssetsToRequirements(
	_ []string,
	version string,
	requirements match.Requirements,
) (string, error) {
	client, err := gitlab.NewClient(g.Token, g.Server)
	if err != nil {
		return "", fmt.Errorf("failed to create GitLab client: %w", err)
	}

	repository := gitlab.NewRepository(g.Owner, g.Repo, client)

	var release *gitlab.Release

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

	err = matches.Status()
	if err != nil {
		err = fmt.Errorf("status: %w", err)
	}

	return assets.FilterByName(matches[0].Asset.Name)[0].URL, err
}

// PopulateOwnerAndRepo sets the Owner and Repo fields based on the given name.
// If Owner and Repo are already set, this method does nothing.
func (g *GitLab) PopulateOwnerAndRepo(name string) error {
	switch {
	case g.Owner != "" && g.Repo != "":
		return nil
	case g.Owner == "" && g.Repo == "":
	default:
		return errors.New("Either both `owner` and `repo` must be set or `name` must be in the format `owner/repo`")
	}

	parts, err := gitlab.SplitName(name)
	if err != nil {
		return err
	}

	// Set the owner and repo fields
	g.Owner = parts[0]
	g.Repo = parts[1]

	return nil
}

// Initialize populates the GitLab repository's owner and name from the given input.
func (g *GitLab) Initialize(name string) error {
	if err := g.PopulateOwnerAndRepo(name); err != nil {
		return err
	}

	return nil
}

// Exe sets the executable name in the metadata to the repository name.
func (g *GitLab) Exe() error {
	g.Data.Set("exe", g.Repo)

	return nil
}

// Version fetches and sets the latest release version in the metadata.
func (g *GitLab) Version(_ string) error {
	version, err := g.LatestVersion()
	if err != nil {
		return err
	}

	g.Data.Set("version", version)

	return nil
}

// Path sets the download URL of the matched asset in the metadata, based on version, file extensions, and requirements.
func (g *GitLab) Path(_ string, extensions []string, version string, requirements match.Requirements) error {
	url, err := g.MatchAssetsToRequirements(extensions, version, requirements)
	if err != nil {
		return err
	}

	g.Data.Set("path", url)

	return nil
}

// Install downloads the asset from GitLab and returns the output, the found file, and any error encountered.
func (g *GitLab) Install(d common.InstallData) (output string, found file.File, err error) {
	d.Header = http.Header{
		"PRIVATE-TOKEN": []string{g.Token},
	}

	return common.Download(d)
}
