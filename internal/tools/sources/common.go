package sources

import (
	"fmt"
	"strings"

	"github.com/idelchi/godyl/pkg/download"
	"github.com/idelchi/godyl/pkg/env"
	"github.com/idelchi/godyl/pkg/file"
	"github.com/idelchi/godyl/pkg/folder"
)

type InstallData struct {
	Path     string
	Name     string
	Exe      string
	Patterns []string
	Output   string
	Aliases  []string
	Mode     string
	Env      env.Env
}

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

func FindAndSymlink(destination file.File, d InstallData) (file.File, error) {
	if destination.IsDir() {
		folder := folder.New(destination.Name())
		// Construct an files item from all the possible names

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

	target := file.New(d.Output, d.Exe)

	if err := destination.Copy(target); err != nil {
		return destination, fmt.Errorf("copying %q to %q: %w", destination, target, err)
	}

	aliases := file.Files{}.FromStrings(d.Output, d.Aliases...)

	return destination, aliases.SymlinksFor(target)
}

func SplitName(name string) (parts [2]string, err error) {
	// Split name by first '/'
	split := strings.Split(name, "/")

	// Check if the name is in the correct format
	if len(split) != 2 {
		return parts, fmt.Errorf("invalid source name: %s", name)
	}

	// Set parts to the split values
	parts[0] = split[0] // owner
	parts[1] = split[1] // repo

	return parts, nil
}
