package gitlab

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-getter/v2"

	"github.com/idelchi/godyl/internal/gitlab"
	"github.com/idelchi/godyl/internal/match"
	"github.com/idelchi/godyl/internal/tools/sources/install"
	"github.com/idelchi/godyl/pkg/path/file"
)

// GitLab represents a GitLab project configuration and state.
type GitLab struct {
	Data                install.Metadata `mapstructure:"-"         yaml:"-"`
	latestStoredRelease *gitlab.Release
	Project             string `mapstructure:"project"   yaml:"project"`
	Namespace           string `mapstructure:"namespace" yaml:"namespace"`
	Token               string `mapstructure:"token"     yaml:"token"     mask:"fixed"`
	Server              string `mapstructure:"server"    yaml:"server"`
	Pre                 bool   `mapstructure:"pre"       yaml:"pre"`
	NoToken             bool   `mapstructure:"no-token"  yaml:"no-token"`
}

// Initialize sets up the GitLab project configuration from the given name.
// Returns an error if the project name format is invalid.
func (g *GitLab) Initialize(name string) error {
	if err := g.PopulateNamespaceAndRepo(name); err != nil {
		return err
	}

	if g.NoToken {
		g.Token = ""
	}

	g.Data.Set("exe", g.Project)

	return nil
}

// Version fetches the latest release version and stores it in metadata.
func (g *GitLab) Version(_ string) error {
	ctx := context.Background()

	version, err := g.LatestVersion(ctx)
	if err != nil {
		return err
	}

	g.Data.Set("version", version)

	return nil
}

// URL finds a matching release asset and stores its URL in metadata.
// Uses version, extensions, and requirements to find the appropriate asset.
func (g *GitLab) URL(_ string, extensions []string, version string, requirements match.Requirements) error {
	ctx := context.Background()

	url, err := g.MatchAssetsToRequirements(ctx, extensions, version, requirements)
	if err != nil {
		return err
	}

	g.Data.Set("url", url)

	return nil
}

// Install downloads the GitLab release asset using the provided configuration.
// Returns the operation output, downloaded file information, and any errors.
func (g *GitLab) Install(
	d install.Data,
	progressListener getter.ProgressTracker,
) (output string, found file.File, err error) {
	d.Header = g.GetHeaders()
	// Pass the progress listener down
	d.ProgressListener = progressListener

	found, err = install.Download(d)

	return "", found, err
}

// Get retrieves a metadata attribute value by its key.
func (g *GitLab) Get(attribute string) string {
	return g.Data.Get(attribute)
}

// LatestVersion fetches the latest release version from GitLab.
// Returns the tag name of the latest release, respecting the Pre flag setting.
func (g *GitLab) LatestVersion(ctx context.Context) (string, error) {
	client, err := gitlab.NewClient(g.Token, g.Server)
	if err != nil {
		return "", fmt.Errorf("failed to create GitLab client: %w", err)
	}

	repository := gitlab.NewRepository(g.Namespace, g.Project, client)

	var release *gitlab.Release

	if g.Pre {
		const PerPage = 1000

		release, err = repository.GetLatestIncludingPreRelease(ctx, PerPage)
	} else {
		release, err = repository.LatestRelease(ctx)
	}

	if err != nil {
		return "", fmt.Errorf("failed to get latest release: %w", err)
	}

	// Store the latest release for future use
	g.latestStoredRelease = release

	return release.Tag, nil
}

// MatchAssetsToRequirements finds release assets matching the given requirements.
// Returns the download URL of the best matching asset, considering platform,
// architecture, and other specified requirements.
func (g *GitLab) MatchAssetsToRequirements(
	ctx context.Context,
	_ []string,
	version string,
	requirements match.Requirements,
) (string, error) {
	client, err := gitlab.NewClient(g.Token, g.Server)
	if err != nil {
		return "", fmt.Errorf("failed to create GitLab client: %w", err)
	}

	repository := gitlab.NewRepository(g.Namespace, g.Project, client)

	var release *gitlab.Release

	if g.latestStoredRelease == nil {
		var releaseErr error

		release, releaseErr = repository.GetRelease(ctx, version)
		if releaseErr != nil {
			return "", fmt.Errorf("failed to get release: %w", releaseErr)
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

// PopulateNamespaceAndRepo sets the Namespace and Project fields from a name string.
// Expects name in "namespace/project" format if fields are not already set.
// Returns an error if the format is invalid or fields are partially set.
func (g *GitLab) PopulateNamespaceAndRepo(name string) (err error) {
	// If both Owner and Repo are already set, nothing to do
	if g.Namespace != "" && g.Project != "" {
		return nil
	}

	// If exactly one of Owner or Repo is set (but not both), that's invalid
	if (g.Namespace == "") != (g.Project == "") {
		return errors.New(
			"either both `namespace` and `repo` must be set or `name` must be in the format `namespace/repo`",
		)
	}

	g.Namespace, g.Project, err = install.CutName(name)
	if err != nil {
		return err
	}

	return nil
}

// GetHeaders returns the HTTP headers required for GitLab API authentication.
func (g *GitLab) GetHeaders() http.Header {
	if g.NoToken {
		return http.Header{}
	}

	return http.Header{
		"PRIVATE-TOKEN": []string{g.Token},
	}
}
