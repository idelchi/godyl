package common

import (
	"fmt"

	"github.com/idelchi/godyl/pkg/download"
	"github.com/idelchi/godyl/pkg/env"
	"github.com/idelchi/godyl/pkg/file"
	"github.com/idelchi/godyl/pkg/folder"
)

// InstallData holds the details required for downloading and installing files,
// including the path, executable name, output directory, and environment settings.
type InstallData struct {
	Path     string   // The URL or path to download from
	Name     string   // The name of the file or project
	Exe      string   // The name of the executable
	Patterns []string // Patterns to match files for the executable
	Output   string   // Output directory for the installation
	Aliases  []string // Aliases for the executable
	Mode     string   // Mode of operation, such as "find" for locating executables
	Env      env.Env  // Environment variables for the installation process
}

// Download handles downloading files based on the InstallData configuration.
// It creates a temporary folder if needed and manages the download process.
func Download(d InstallData) (string, file.File, error) {
	var err error
	var found file.File

	folder := folder.Folder(d.Output)

	if d.Mode == "find" {
		if err := folder.CreateRandomInTempDir(); err != nil {
			return "", "", fmt.Errorf("creating temp dir: %w", err)
		}
		defer func() {
			if err == nil {
				folder.Remove()
			}
		}()
	}

	downloader := download.New()

	destination, err := downloader.Download(d.Path, folder.Path())
	if err != nil {
		return "", "", fmt.Errorf("downloading %q: %w", d.Path, err)
	}

	if d.Mode == "find" {
		found, err = FindAndSymlink(destination, d)
	}

	return "", found, err
}

// FindAndSymlink finds the executable within the downloaded folder and creates symlinks for it
// based on the provided InstallData. It handles directories and sets up aliases as needed.
func FindAndSymlink(destination file.File, d InstallData) (file.File, error) {
	if destination.IsDir() {
		folder := folder.New(destination.Name())

		// Construct files item from the specified patterns
		files := file.Files{}.FromStrings("", d.Patterns...)

		// Find the specific executable that was downloaded
		var err error
		destination, err = files.Find(folder.Path())
		if err != nil {
			return destination, fmt.Errorf("finding executable: %w", err)
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
	aliases := file.Files{}.FromStrings(d.Output, d.Aliases...)
	return destination, aliases.SymlinksFor(target)
}
