// Package goc provides functionality for handling Go-based installations and
// managing Go commands using GitHub repositories. It integrates with GitHub
// to fetch and install Go projects, set paths, and manage metadata related
// to the installation process.
package goc

import (
	"fmt"
	"path/filepath"
	"slices"
	"strings"
	"sync"

	"github.com/hashicorp/go-getter/v2"

	"github.com/idelchi/godyl/internal/goi"
	"github.com/idelchi/godyl/internal/match"
	"github.com/idelchi/godyl/internal/tmp"
	"github.com/idelchi/godyl/internal/tools/sources/common"
	"github.com/idelchi/godyl/internal/tools/sources/github"
	"github.com/idelchi/godyl/pkg/path/file"
)

// Go represents a Go project configuration that can be installed from GitHub.
type Go struct {
	github  *github.GitHub
	Data    common.Metadata `yaml:"-"`
	Command string          `yaml:"command"`
	Base    string          `yaml:"base"`
}

// Initialize sets up the Go project configuration from the given name.
// Uses the associated GitHub repository for initialization.
//
// TODO(Idelchi): This should be ignored if the version is already set.
// As a workaround, just return nil for now.
func (g *Go) Initialize(name string) error {
	if g.Base != "github.com" {
		g.github.Repo = name

		g.Data.Set("exe", name)

		return nil
	}

	return g.github.Initialize(name)
}

// Version fetches the latest release version and stores it in metadata.
func (g *Go) Version(name string) error {
	return g.github.Version(name)
}

// URL constructs and stores the Go module path in metadata.
// Uses the format github.com/{owner}/{repo}@{version}.
func (g *Go) URL(_ string, _ []string, version string, _ match.Requirements) error {
	parts := []string{g.Base, g.github.Owner, g.github.Repo}

	parts = slices.DeleteFunc(parts, func(s string) bool { return s == "" })

	g.github.Data.Set("url", fmt.Sprintf("%s@%s", strings.Join(parts, "/"), version))

	return nil
}

// TODO(Idelchi): Remove this at some point - why do we need it?
var mu sync.Mutex

// Install downloads and builds the Go project using 'go install'.
// Handles temporary directory creation, environment setup, and file linking.
// Progress listener is accepted but not used as go install doesn't support it.
// Returns the installation output, installed file information, and any errors.
func (g *Go) Install(d common.InstallData, _ getter.ProgressTracker) (output string, found file.File, err error) {
	mu.Lock()

	binary, err := goi.New(d.NoVerifySSL)
	if err != nil {
		mu.Unlock()

		return "", "", err
	}
	mu.Unlock()

	installer := goi.Installer{
		Binary: binary,
	}

	folder, err := tmp.GodylCreateRandomDir()
	if err != nil {
		return "", "", fmt.Errorf("creating random dir: %w", err)
	}

	installer.Binary.Env.Append(
		goi.Env{
			"GOBIN": folder.Path(),
		},
	)

	name := strings.TrimSuffix(d.Exe, filepath.Ext(d.Exe))

	mod, version, ok := strings.Cut(d.Path, "@")
	if !ok {
		// fallback or panic, depending on your context
		return "", "", fmt.Errorf("invalid module path: %s", d.Path)
	}

	paths := []string{
		fmt.Sprintf("%s@%s", mod, version),
		fmt.Sprintf("%s/cmd/%s@%s", mod, name, version),
		fmt.Sprintf("%s/cmd@%s", mod, version),
	}

	if g.Command != "" {
		paths = []string{
			fmt.Sprintf("%s/%s@%s", mod, g.Command, version),
		}
	}

	defer func() {
		if err == nil {
			folder.Remove() //nolint:gosec 		// TODO(Idelchi): Address this later.
		}
	}()

	for _, path := range paths {
		output, err = installer.Install(path)

		if err == nil {
			d.Path = path
			found, err := common.FindAndSymlink(file.New(folder.Path()), d)

			return output, found, err
		}
	}

	return output, "", err
}

// Get retrieves a metadata attribute value by its key.
func (g *Go) Get(attribute string) string {
	return g.github.Data.Get(attribute)
}

// SetGitHub configures the GitHub repository for the Go project.
func (g *Go) SetGitHub(gh *github.GitHub) {
	g.github = gh
}
