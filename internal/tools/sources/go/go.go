// Package goc provides functionality for handling Go-based installations and
// managing Go commands using GitHub repositories. It integrates with GitHub
// to fetch and install Go projects, set paths, and manage metadata related
// to the installation process.
package goc

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"sync"

	"github.com/hashicorp/go-getter/v2"

	"github.com/idelchi/godyl/internal/debug"
	"github.com/idelchi/godyl/internal/detect/platform"
	"github.com/idelchi/godyl/internal/goi"
	"github.com/idelchi/godyl/internal/match"
	"github.com/idelchi/godyl/internal/tools/sources/github"
	"github.com/idelchi/godyl/internal/tools/sources/install"
	progresspkg "github.com/idelchi/godyl/pkg/download/progress"
	"github.com/idelchi/godyl/pkg/path/file"
	"github.com/idelchi/godyl/pkg/path/folder"
)

// Go represents a Go project configuration that can be installed from GitHub.
type Go struct {
	github            *github.GitHub
	Data              install.Metadata `yaml:"-"`
	Command           string           `yaml:"command"`
	Base              string           `yaml:"base"`
	DownloadIfMissing bool             `yaml:"download_if_missing"`
}

// Initialize sets up the Go project configuration from the given name.
// Uses the associated GitHub repository for initialization.
//
// TODO(Idelchi): This should be ignored if the version is already set. As a workaround, just return nil for now.
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
var mu sync.Mutex //nolint:gochecknoglobals // Global mutex for thread-safe access across package

// Install downloads and builds the Go project using 'go install'.
// Handles temporary directory creation, environment setup, and file linking.
// Returns the installation output, installed file information, and any errors.
//
//nolint:funlen // TODO(Idelchi): Refactor later.
func (g *Go) Install(
	d install.Data,
	progressListener getter.ProgressTracker,
) (output string, found file.File, err error) {
	mu.Lock()

	debug.Debug("Searching for go binary...")

	binary, err := goi.New(d.NoVerifySSL, g.DownloadIfMissing, d.NoVerifyChecksum, progressListener)
	if err != nil {
		mu.Unlock()

		return "", "", err
	}

	debug.Debug("Go binary found: %q", binary.File)

	mu.Unlock()

	installer := goi.Installer{
		Binary: binary,
		GOOS:   d.OS,
		GOARCH: d.Arch,
	}

	crossCompiling := d.OS != runtime.GOOS || d.Arch != runtime.GOARCH

	var outFolder folder.Folder

	if value, ok := installer.Binary.Env["GOPATH"]; ok {
		outFolder = folder.New(filepath.Join(value, "bin"))
	} else if bin, ok := os.LookupEnv("GOBIN"); ok {
		outFolder = folder.New(bin)
	} else if goPath, ok := os.LookupEnv("GOPATH"); ok {
		outFolder = folder.New(filepath.Join(goPath, "bin"))
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			home = os.Getenv("HOME")
		}

		outFolder = folder.New(filepath.Join(home, "go", "bin"))
	}

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

	debug.Debug("Setting up progress tracker for 'go install'...")

	stopProgress := startGoInstallProgress(progressListener, paths[0])
	defer stopProgress()

	var (
		os        platform.OS
		extension platform.Extension
	)

	_ = os.ParseFrom(runtime.GOOS)

	for _, pth := range paths {
		debug.Debug("Attempting to install from path: %q", pth)

		output, err = installer.Install(pth)
		debug.Debug("go install output: %q", output)
		debug.Debug("successful path was: %q", pth)

		if err == nil {
			d.Path = pth

			last := path.Base(pth)
			name, _, _ := strings.Cut(last, "@")

			if crossCompiling {
				outFolder = outFolder.Join(fmt.Sprintf("%s_%s", d.OS, d.Arch))

				_ = os.ParseFrom(d.OS)
			}

			extension.ParseFrom(os)

			name = fmt.Sprintf("%s%s", name, extension)

			d.Patterns = []string{name}

			debug.Debug("Searching in %q for %q", outFolder.Path(), d.Patterns)

			debug.Debug("Linking to %q", filepath.Join(outFolder.Path(), name))

			found, findErr := install.FindAndSymlink(
				file.New(outFolder.Path()),
				d,
			)

			return output, found, findErr
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

func startGoInstallProgress(progressListener getter.ProgressTracker, label string) func() {
	valueLabel := "(go install)"
	speedLabel := "(n/a)"
	message := fmt.Sprintf("%-45s %s", file.New(label).Unescape(), valueLabel)

	tracker, ok := progressListener.(*progresspkg.Tracker)
	if !ok {
		return func() {}
	}

	const stallFraction = 0.8

	return progresspkg.StartSynthetic(tracker, label, message, speedLabel, speedLabel, stallFraction)
}
