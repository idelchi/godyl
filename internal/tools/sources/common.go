package sources

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/idelchi/godyl/internal/executable"
	"github.com/idelchi/godyl/internal/folder"
	"github.com/idelchi/godyl/pkg/download"
)

type InstallData struct {
	Path    string
	Name    string
	Exe     string
	Output  string
	Aliases []string
}

func Download(d InstallData) (output string, err error) {
	var tmp folder.Folder
	if err := tmp.CreateRandomInTempDir(); err != nil {
		return "", fmt.Errorf("creating temp dir: %w", err)
	}
	defer tmp.Remove()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	destination, err := download.Download(ctx, d.Path, tmp.Path())
	if err != nil {
		return "", fmt.Errorf("downloading %q: %w", d.Path, err)
	}

	return "", FindAndSymlink(destination, d)
}

func FindAndSymlink(destination string, d InstallData) error {
	// Construct an executables item from all the possible names
	executables := executable.Executables{}.FromStrings("", append([]string{d.Name, d.Exe, filepath.Base(d.Path)}, d.Aliases...)...)

	// Find the specific executable that was downloaded
	download, err := executables.Find(destination)
	if err != nil {
		return fmt.Errorf("finding executable: %w", err)
	}

	folder := folder.Folder(d.Output)
	if !folder.Exists() {
		if err := folder.Create(); err != nil {
			return fmt.Errorf("creating output folder: %w", err)
		}
	}

	target := executable.New(d.Output, d.Exe)

	if err := download.Copy(target.Path); err != nil {
		return fmt.Errorf("copying %q to %q: %w", download.Path, target.Path, err)
	}

	aliases := executable.Executables{}.FromStrings(d.Output, d.Aliases...)

	return aliases.SymlinksFor(target)
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
