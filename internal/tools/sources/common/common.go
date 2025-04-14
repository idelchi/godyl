package common

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"

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
	// Path is the URL or filesystem path to download from.
	Path string

	// Name is the identifier of the tool or project.
	Name string

	// Exe is the name of the executable file.
	Exe string

	// Patterns contains regex patterns for finding the executable.
	Patterns []string

	// Output specifies the directory where files will be installed.
	Output string

	// Aliases are alternative names for the executable.
	Aliases []string

	// Mode defines the installation behavior (e.g., "find" for locating executables).
	Mode string

	// Env holds environment variables for the installation process.
	Env env.Env

	// NoVerifySSL disables SSL certificate verification when downloading.
	NoVerifySSL bool

	// Header contains HTTP headers for download requests.
	Header http.Header

	// ProgressListener tracks download progress.
	ProgressListener getter.ProgressTracker
}

// Download retrieves files according to the InstallData configuration.
// Creates temporary directories when needed, manages the download process,
// and returns the download output and file information.
func Download(data InstallData) (string, file.File, error) {
	var err error

	var found file.File

	dir := folder.New(data.Output)

	if data.Mode == "find" {
		if dir, err = tmp.GodylCreateRandomDir(); err != nil {
			return "", "", fmt.Errorf("creating random dir: %w", err)
		}

		defer func() {
			dir.Remove() //nolint:gosec 		// TODO(Idelchi): Address this later.
		}()
	}

	downloader := download.New()
	downloader.InsecureSkipVerify = data.NoVerifySSL
	// Pass the progress listener if provided in InstallData
	if data.ProgressListener != nil {
		downloader.ProgressListener = data.ProgressListener
	}

	destination, err := downloader.Download(data.Path, dir.Path(), data.Header)
	if err != nil {
		return "", "", fmt.Errorf("downloading %q: %w", data.Path, err)
	}

	if data.Mode == "find" {
		found, err = FindAndSymlink(destination, data)
	}

	return "", found, err
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
				re, err := regexp.Compile(pattern)
				if err != nil {
					return false, fmt.Errorf("compiling pattern %q: %w", pattern, err)
				}

				matched := re.MatchString(file.Path())

				if matched {
					return true, nil
				}

				return false, nil
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
				"finding executable: no executable matching patterns %v found in %q",
				d.Patterns,
				searchDir,
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

	// Create symlinks for the aliases
	aliases := files.New(d.Output, d.Aliases...)

	return destination, aliases.SymlinksFor(target)
}
