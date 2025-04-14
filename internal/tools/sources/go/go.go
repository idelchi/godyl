// Package goc provides functionality for handling Go-based installations and
// managing Go commands using GitHub repositories. It integrates with GitHub
// to fetch and install Go projects, set paths, and manage metadata related
// to the installation process.
package goc

import (
	"fmt"
	"path/filepath"
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
	// Github is the underlying GitHub repository configuration.
	github *github.GitHub

	// Command is an optional custom command to install instead of the default(s).
	Command string `yaml:"command"`

	// Data contains additional metadata about the Go project.
	Data common.Metadata `yaml:"-"`
}

// SetGitHub configures the GitHub repository for the Go project.
func (g *Go) SetGitHub(gh *github.GitHub) {
	g.github = gh
}

// Get retrieves a metadata attribute value by its key.
func (g *Go) Get(attribute string) string {
	return g.github.Data.Get(attribute)
}

// Initialize sets up the Go project configuration from the given name.
// Uses the associated GitHub repository for initialization.
func (g *Go) Initialize(name string) error {
	return g.github.Initialize(name)
}

// Exe stores the project name as the executable name in metadata.
func (g *Go) Exe() error {
	return g.github.Exe()
}

// Version fetches the latest release version and stores it in metadata.
func (g *Go) Version(name string) error {
	return g.github.Version(name)
}

// Path constructs and stores the Go module path in metadata.
// Uses the format github.com/{owner}/{repo}@{version}.
func (g *Go) Path(_ string, _ []string, version string, _ match.Requirements) error {
	g.github.Data.Set("path", fmt.Sprintf("github.com/%s/%s@%s", g.github.Owner, g.github.Repo, version))

	return nil
}

var mu sync.Mutex

// Install downloads and builds the Go project using 'go install'.
// Handles temporary directory creation, environment setup, and file linking.
// Progress listener is accepted but not used as go install doesn't support it.
// Returns the installation output, installed file information, and any errors.
func (g *Go) Install(d common.InstallData, _ getter.ProgressTracker) (output string, found file.File, err error) {
	mu.Lock()

	binary, err := goi.New(d.NoVerifySSL)
	if err != nil {
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

	paths := []string{
		d.Path,
		strings.Replace(d.Path, fmt.Sprintf("/%s@", name), fmt.Sprintf("/%s/cmd/%s@", name, name), 1),
		strings.Replace(d.Path, fmt.Sprintf("/%s@", name), fmt.Sprintf("/%s/cmd@", name), 1),
	}

	if g.Command != "" {
		paths = []string{
			strings.Replace(d.Path, fmt.Sprintf("/%s@", name), fmt.Sprintf("/%s/%s@", name, g.Command), 1),
		}
	}

	for i, path := range paths {
		paths[i] = strings.ToLower(path)
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
