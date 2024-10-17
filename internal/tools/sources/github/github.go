package github

import (
	"fmt"

	"github.com/idelchi/godyl/internal/github"
	"github.com/idelchi/godyl/internal/match"
	"github.com/idelchi/godyl/internal/tools/sources/common"
	"github.com/idelchi/godyl/pkg/file"
)

// GitHub represents a GitHub repository with optional authentication token and metadata.
type GitHub struct {
	Repo  string
	Owner string
	Token string `mask:"fixed"`

	// Data holds additional metadata related to the repository.
	Data common.Metadata `yaml:"-"`
}

// Get retrieves a specific attribute from the GitHub repository's metadata.
func (g *GitHub) Get(attribute string) string {
	return g.Data.Get(attribute)
}

// LatestVersion fetches the latest release version of the GitHub repository.
func (g *GitHub) LatestVersion() (string, error) {
	client := github.NewClient(g.Token)
	repository := github.NewRepository(g.Owner, g.Repo, client)

	release, err := repository.LatestRelease()
	if err != nil {
		return "", err
	}

	return release.Tag, nil
}

// MatchAssetsToRequirements matches release assets to specific file extensions and requirements,
// returning the URL of the matched asset.
func (g *GitHub) MatchAssetsToRequirements(
	filters []string,
	version string,
	requirements match.Requirements,
) (string, error) {
	client := github.NewClient(g.Token)
	repository := github.NewRepository(g.Owner, g.Repo, client)

	release, err := repository.GetRelease(version)
	if err != nil {
		return "", err
	}

	assets := release.Assets

	assets = assets.FilterByExtension(filters)

	match, err := assets.Match(requirements)
	if err != nil {
		return "", err
	}

	// TODO: Will this fail if match is empty?
	return assets.FilterByName(match[0].Name)[0].URL, match.Status()
}

// PopulateOwnerAndRepo sets the Owner and Repo fields based on the given name.
// If Owner and Repo are already set, this method does nothing.
func (g *GitHub) PopulateOwnerAndRepo(name string) error {
	switch {
	case g.Owner != "" && g.Repo != "":
		return nil
	case g.Owner == "" && g.Repo == "":
	default:
		return fmt.Errorf("Either both `owner` and `repo` must be set or `name` must be in the format `owner/repo`")
	}

	parts, err := SplitName(name)
	if err != nil {
		return err
	}

	// Set the owner and repo fields
	g.Owner = parts[0]
	g.Repo = parts[1]

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
func (g *GitHub) Version(name string) error {
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
