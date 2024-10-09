package sources

import (
	"fmt"

	"github.com/idelchi/godyl/internal/github"
	"github.com/idelchi/godyl/internal/match"
)

type GitHub struct {
	Repo  string
	Owner string
	Token string

	Data Metadata `yaml:"-"`
}

func (g *GitHub) Get(attribute string) string {
	return g.Data.Get(attribute)
}

func (g *GitHub) LatestVersion() (string, error) {
	client := github.NewClient(g.Token)
	repository := github.NewRepository(g.Owner, g.Repo, client)

	release, err := repository.LatestRelease()
	if err != nil {
		return "", err
	}

	return release.Tag, nil
}

func (g *GitHub) MatchAssetsToRequirements(filters []string, version string, requirements match.Requirements) (string, error) {
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

	// TODO(Idelchi): Will this fail if match is empty?
	return assets.FilterByName(match[0].Name)[0].URL, match.Status()
}

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

func (g *GitHub) Initialize(name string) error {
	if err := g.PopulateOwnerAndRepo(name); err != nil {
		return err
	}

	return nil
}

func (g *GitHub) Exe() error {
	g.Data.Set("exe", g.Repo)

	return nil
}

func (g *GitHub) Version(name string) error {
	version, err := g.LatestVersion()
	if err != nil {
		return err
	}

	g.Data.Set("version", version)

	return nil
}

func (g *GitHub) Path(_ string, extensions []string, version string, requirements match.Requirements) error {
	url, err := g.MatchAssetsToRequirements(extensions, version, requirements)
	if err != nil {
		return err
	}

	g.Data.Set("path", url)

	return nil
}

func (g *GitHub) Install(d InstallData) (output string, err error) {
	return Download(d)
}
