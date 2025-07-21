package common

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-getter/v2"

	"github.com/idelchi/godyl/internal/tmp"
	"github.com/idelchi/godyl/pkg/download"
	"github.com/idelchi/godyl/pkg/env"
	"github.com/idelchi/godyl/pkg/path/file"
	"github.com/idelchi/godyl/pkg/path/files"
	"github.com/idelchi/godyl/pkg/path/folder"
)

// InstallData contains configuration for downloading and installing tools.
type InstallData struct {
	ProgressListener getter.ProgressTracker
	Env              env.Env
	Header           http.Header
	Path             string
	Name             string
	Exe              string
	Output           string
	Mode             string
	Patterns         []string
	Aliases          []string
	NoVerifySSL      bool
	OS               string // Target operating system for cross-compilation.
	Arch             string // Target architecture for cross-compilation.
}

// Download retrieves files according to the InstallData configuration.
// Creates temporary directories when needed, manages the download process,
// and returns the download output and file information.
func Download(data InstallData) (file.File, error) {
	var err error

	var found file.File

	dir := folder.New(data.Output)

	if data.Mode == "find" {
		if dir, err = tmp.GodylCreateRandomDir(); err != nil {
			return "", fmt.Errorf("creating random dir: %w", err)
		}

		defer func() {
			dir.Remove() //nolint:gosec 		// TODO(Idelchi): Address this later.
		}()
	}

	options := []download.Option{download.WithProgress(data.ProgressListener)}
	if data.NoVerifySSL {
		options = append(options, download.WithInsecureSkipVerify())
	}

	downloader := download.New(options...)

	destination, err := downloader.Download(data.Path, dir.Path(), data.Header)
	if err != nil {
		return "", fmt.Errorf("downloading %q: %w", data.Path, err)
	}

	if data.Mode == "find" {
		found, err = FindAndSymlink(destination, data)
	}

	return found, err
}

// FindAndSymlink locates an executable in the downloaded content and sets up symlinks.
// It searches directories recursively using provided patterns, copies the executable
// to the output location, and creates symlinks for any configured aliases.
//
//nolint:gocognit   // TODO(Idelchi): Address this later.
func FindAndSymlink(destination file.File, d InstallData) (file.File, error) {
	if destination.IsDir() {
		searchDir := folder.New(destination.Dir())

		folders, err := searchDir.ListFolders()
		if err != nil {
			return destination, fmt.Errorf("listing folders in %q: %w", searchDir, err)
		}

		files, err := searchDir.ListFiles()
		if err != nil {
			return destination, fmt.Errorf("listing files in %q: %w", searchDir, err)
		}

		if len(folders) == 1 && len(files) == 0 {
			searchDir = folders[0]
		}

		var found bool
		// Match patterns in order of priority
		for _, pattern := range d.Patterns {
			match := func(file file.File) (bool, error) {
				return file.Matches(pattern)
			}

			var err error

			destination, err = searchDir.FindFile(match)
			if err != nil {
				if !errors.Is(err, folder.ErrNotFound) {
					return destination, fmt.Errorf("finding executable: %w", err)
				}

				continue
			} else {
				found = true

				break
			}
		}

		if !found {
			return destination, fmt.Errorf(
				"finding executable: no executable matching patterns %v found in %q: found\n%v",
				d.Patterns,
				searchDir,
				files,
			)
		}
	}

	folder := folder.New(d.Output)
	if !folder.Exists() {
		if err := folder.Create(); err != nil {
			return destination, fmt.Errorf("creating output folder: %w", err)
		}
	}

	// Copy the executable to the output directory
	target := file.New(d.Output, d.Exe)
	if err := destination.Copy(target); err != nil {
		return destination, fmt.Errorf("copying %q to %q: %w", destination, target, err)
	}

	if ok, _ := target.IsExecutable(); !ok {
		if err := target.MakeExecutable(); err != nil {
			return destination, fmt.Errorf("making %q executable: %w", target, err)
		}
	}

	// Create symlinks for the aliases
	aliases := files.New(d.Output, d.Aliases...)

	return destination, aliases.LinksFor(target)
}
