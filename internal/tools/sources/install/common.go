package install

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-getter/v2"

	"github.com/idelchi/godyl/internal/data"
	"github.com/idelchi/godyl/pkg/download"
	"github.com/idelchi/godyl/pkg/env"
	"github.com/idelchi/godyl/pkg/path/file"
	"github.com/idelchi/godyl/pkg/path/files"
	"github.com/idelchi/godyl/pkg/path/folder"
)

// Data contains configuration for downloading and installing tools.
type Data struct {
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
func Download(d Data) (found file.File, err error) {
	dir := folder.New(d.Output)

	if d.Mode == "find" {
		if dir, err = data.CreateUniqueDirIn(); err != nil {
			return "", fmt.Errorf("creating random dir: %w", err)
		}

		defer func() {
			err = errors.Join(err, dir.Remove())
		}()
	}

	options := []download.Option{download.WithProgress(d.ProgressListener)}
	if d.NoVerifySSL {
		options = append(options, download.WithInsecureSkipVerify())
	}

	downloader := download.New(options...)

	destination, err := downloader.Download(d.Path, dir.Path(), d.Header)
	if err != nil {
		return "", fmt.Errorf("downloading %q: %w", d.Path, err)
	}

	if d.Mode == "find" {
		found, err = FindAndSymlink(destination, d)
	}

	return found, err
}

// findExecutableInDir searches for an executable file in a directory using the provided patterns.
func findExecutableInDir(destination file.File, patterns []string) (file.File, error) {
	searchDir := folder.New(destination.Dir())

	folders, err := searchDir.ListFolders()
	if err != nil {
		return destination, fmt.Errorf("listing folders in %q: %w", searchDir, err)
	}

	files, err := searchDir.ListFiles()
	if err != nil {
		return destination, fmt.Errorf("listing files in %q: %w", searchDir, err)
	}

	// If there's only one folder and no files, search within that folder
	if len(folders) == 1 && len(files) == 0 {
		searchDir = folders[0]
	}

	// Try each pattern in order
	for _, pattern := range patterns {
		match := func(file file.File) (bool, error) {
			return file.Matches(pattern)
		}

		found, err := searchDir.FindFile(match)
		if err != nil {
			if !errors.Is(err, folder.ErrNotFound) {
				return destination, fmt.Errorf("finding executable: %w", err)
			}

			continue
		}

		return found, nil
	}

	return destination, fmt.Errorf(
		"finding executable: no executable matching patterns %v found in %q: found\n%v",
		patterns,
		searchDir,
		files,
	)
}

// FindAndSymlink locates an executable in the downloaded content and sets up symlinks.
// It searches directories recursively using provided patterns, copies the executable
// to the output location, and creates symlinks for any configured aliases.
func FindAndSymlink(destination file.File, d Data) (file.File, error) {
	if destination.IsDir() {
		var err error

		destination, err = findExecutableInDir(destination, d.Patterns)
		if err != nil {
			return destination, err
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
	if len(d.Aliases) > 0 {
		aliases := files.New(d.Output, d.Aliases...)

		return destination, aliases.LinksFor(target)
	}

	return destination, nil
}
