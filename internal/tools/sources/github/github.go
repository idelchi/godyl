package github

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/go-getter/v2"

	"github.com/idelchi/godyl/internal/debug"
	"github.com/idelchi/godyl/internal/github"
	"github.com/idelchi/godyl/internal/match"
	"github.com/idelchi/godyl/internal/tools/sources/install"
	"github.com/idelchi/godyl/pkg/path/file"
)

// GitHub represents a GitHub repository configuration and state.
type GitHub struct {
	Data                install.Metadata `mapstructure:"-" yaml:"-"`
	latestStoredRelease *github.Release
	Repo                string `mapstructure:"repo"  yaml:"repo"`
	Owner               string `mapstructure:"owner" yaml:"owner"`
	Token               string `mapstructure:"token" mask:"fixed" yaml:"token"`
	Pre                 bool   `mapstructure:"pre"   yaml:"pre"`
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
func (g *GitHub) Version(version string) error {
	ctx := context.Background()

	version, err := g.LatestVersion(ctx, version)
	if err != nil {
		return err
	}

	g.Data.Set("version", version)

	return nil
}

// URL finds a matching release asset and stores its URL in metadata.
// Uses version, extensions, and requirements to find the appropriate asset.
func (g *GitHub) URL(_ string, extensions []string, version string, requirements match.Requirements) error {
	ctx := context.Background()

	url, err := g.MatchAssetsToRequirements(ctx, extensions, version, requirements)
	if err != nil {
		return err
	}

	g.Data.Set("url", url)

	return nil
}

// Install downloads the GitHub release asset using the provided configuration.
// Returns the operation output, downloaded file information, and any errors.
func (g *GitHub) Install(
	d install.Data,
	progressListener getter.ProgressTracker,
) (output string, found file.File, err error) {
	// Pass the progress listener down to the common download function
	d.ProgressListener = progressListener

	found, err = install.Download(d)

	return "", found, err
}

// Get retrieves a metadata attribute value by its key.
func (g *GitHub) Get(attribute string) string {
	return g.Data.Get(attribute)
}

// LatestVersion fetches the latest release version from GitHub.
// Returns the tag name of the latest release, respecting the Pre flag setting.
func (g *GitHub) LatestVersion(ctx context.Context, version string) (string, error) {
	client := github.NewClient(g.Token)
	repository := github.NewRepository(g.Owner, g.Repo, client)

	var release *github.Release

	var err error

	switch {
	case strings.Contains(version, "*"):
		const PerPage = 100

		release, err = repository.GetReleasesByWildcard(ctx, version, PerPage)

	case g.Pre:
		const PerPage = 100

		release, err = repository.LatestIncludingPreRelease(
			ctx,
			PerPage,
		)
	default:
		if g.Token == "" {
			if tag, webErr := repository.LatestVersionFromWebJSON(ctx); webErr == nil {
				return tag, nil
			}
		}

		release, err = repository.LatestRelease(ctx)
	}

	if err != nil {
		return "", fmt.Errorf("failed to retrieve latest release: %w", err)
	}

	// Store the latest release for future use
	g.latestStoredRelease = release
	g.Data.Set("body", release.Body)

	return release.Tag, nil
}

// MatchAssetsToRequirements finds release assets matching the given requirements.
// Returns the download URL of the best matching asset, considering platform,
// architecture, and other specified requirements.
func (g *GitHub) MatchAssetsToRequirements(
	ctx context.Context,
	_ []string,
	version string,
	requirements match.Requirements,
) (string, error) {
	client := github.NewClient(g.Token)
	repository := github.NewRepository(g.Owner, g.Repo, client)

	var release *github.Release

	if g.latestStoredRelease == nil { //nolint:nestif // Multiple checks are necessary
		var err error

		if g.Token == "" {
			release, err = repository.GetReleaseFromWeb(ctx, version)
		}

		if err != nil || release == nil {
			release, err = repository.GetRelease(ctx, version)
		}

		if err != nil {
			return "", fmt.Errorf("failed to get release: %w", err)
		}

		if release == nil {
			return "", errors.New("failed to get release: release is nil")
		}
	} else {
		release = g.latestStoredRelease
	}

	assets := release.Assets

	matches := assets.Match(requirements)

	if matches.HasErrors() {
		return "", matches.Errors()[0]
	}

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

	asset := assets.FilterByName(matches[0].Asset.Name)[0]

	// Check inline digest
	if asset.Digest != "" {
		debug.Debug("found asset with digest: %q", asset.Digest)
		g.Data.Set("checksum", asset.Digest)
	} else if checksums := assets.Checksums(requirements.Checksum); len(checksums) > 0 {
		debug.Debug("found checksum assets: %q", checksums)

		preferred := checksums.Preferred(asset.Name)
		if preferred != "" {
			checksum := assets.FilterByName(preferred)[0]
			g.Data.Set("checksum", checksum.URL)
			debug.Debug("using preferred checksum asset: %q from %q", checksum.URL, asset.Name)
		}
	}

	return asset.URL, err
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
		return errors.New("either both `owner` and `repo` must be set or `name` must be in the format `owner/repo`")
	}

	g.Owner, g.Repo, err = install.SplitName(name)
	if err != nil {
		return err
	}

	return nil
}
