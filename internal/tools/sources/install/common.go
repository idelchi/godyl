package install

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/go-getter/v2"

	"github.com/idelchi/godyl/internal/data"
	"github.com/idelchi/godyl/internal/tools/checksum"
	"github.com/idelchi/godyl/pkg/download"
	"github.com/idelchi/godyl/pkg/env"
	"github.com/idelchi/godyl/pkg/path/file"
	"github.com/idelchi/godyl/pkg/path/folder"
)

// Data contains configuration for downloading and installing tools.
type Data struct {
	// ProgressListener receives download progress updates.
	ProgressListener getter.ProgressTracker
	// Env contains environment overrides for installer commands.
	Env env.Env
	// Checksum defines artifact integrity verification.
	Checksum checksum.Checksum
	// Header contains HTTP headers sent with artifact requests.
	Header http.Header
	// Path is the resolved artifact URL or module path.
	Path string
	// Name is the configured tool name.
	Name string
	// Exe is the destination executable name.
	Exe string
	// Output is the installation directory.
	Output string
	// Mode controls how downloaded content is installed.
	Mode string
	// Patterns contains candidate executable filenames.
	Patterns []string
	// NoVerifySSL disables TLS certificate verification.
	NoVerifySSL bool
	// NoVerifyChecksum disables artifact checksum verification.
	NoVerifyChecksum bool
	// OS is the target operating system for cross-compilation.
	OS string
	// Arch is the target architecture for cross-compilation.
	Arch string
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

	destination, err := Acquire(d, dir)
	if err != nil {
		return "", err
	}

	if d.Mode == "find" {
		found, err = Find(destination, d)
	}

	return found, err
}

// Acquire downloads and extracts a tool artifact into dir.
func Acquire(d Data, dir folder.Folder) (file.File, error) {
	options := []download.Option{
		download.WithProgress(d.ProgressListener),
		download.WithContextTimeout(download.DefaultTimeout),
	}
	if d.NoVerifySSL {
		options = append(options, download.WithInsecureSkipVerify())
	}

	if d.Checksum.IsMandatory() && !d.NoVerifyChecksum {
		options = append(options, download.WithChecksum(d.Checksum.ToQuery()))
	}

	downloader := download.New(options...)

	destination, err := downloader.Download(d.Path, dir.Path(), d.Header)
	if err != nil {
		return "", fmt.Errorf("downloading %q: %w", d.Path, err)
	}

	return destination, nil
}

// findExecutableInDir searches for an executable file in a directory using the provided patterns.
func findExecutableInDir(destination file.File, patterns []string) (file.File, error) {
	searchDir := folder.New(destination.Dir())

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

	allFiles, _ := searchDir.FindFiles(func(file file.File) (bool, error) {
		return file.Matches("**")
	})

	return destination, fmt.Errorf(
		"finding executable: no executable matching patterns %v found in %q: found\n%v",
		patterns,
		searchDir,
		"- "+strings.Join(allFiles.RelativeTo(searchDir.String()).AsSlice(), "\n- "),
	)
}

// Find locates an executable in the downloaded content.
// It searches directories recursively using provided patterns and copies the executable
// to the output location.
func Find(destination file.File, d Data) (file.File, error) {
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

	return destination, nil
}
