package common

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/idelchi/godyl/internal/tmp"
	"github.com/idelchi/godyl/pkg/download"
	"github.com/idelchi/godyl/pkg/env"
	"github.com/idelchi/godyl/pkg/file"
	"github.com/idelchi/godyl/pkg/files"
	"github.com/idelchi/godyl/pkg/folder"
)

// InstallData holds the details required for downloading and installing files,
// including the path, executable name, output directory, and environment settings.
type InstallData struct {
	Path        string   // The URL or path to download from
	Name        string   // The name of the file or project
	Exe         string   // The name of the executable
	Patterns    []string // Patterns to match files for the executable
	Output      string   // Output directory for the installation
	Aliases     []string // Aliases for the executable
	Mode        string   // Mode of operation, such as "find" for locating executables
	Env         env.Env  // Environment variables for the installation process
	NoVerifySSL bool     // Skip SSL verification
}

// Download handles downloading files based on the InstallData configuration.
// It creates a temporary folder if needed and manages the download process.
func Download(data InstallData) (string, file.File, error) {
	var err error

	var found file.File

	dir := folder.New(data.Output)

	if data.Mode == "find" {
		if dir, err = tmp.GodylCreateRandomDir(); err != nil {
			return "", "", fmt.Errorf("creating random dir: %w", err)
		}

		defer func() {
			if err == nil {
				dir.Remove() //nolint:gosec 		// TODO(Idelchi): Address this later.
			}
		}()
	}

	downloader := download.New()
	downloader.InsecureSkipVerify = data.NoVerifySSL

	destination, err := downloader.Download(data.Path, dir.Path())
	if err != nil {
		return "", "", fmt.Errorf("downloading %q: %w", data.Path, err)
	}

	if data.Mode == "find" {
		found, err = FindAndSymlink(destination, data)
	}

	return "", found, err
}

// FindAndSymlink finds the executable within the downloaded folder and creates symlinks for it
// based on the provided InstallData. It handles directories and sets up aliases as needed.
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

				matched := re.MatchString(file.Normalized().Path())

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
