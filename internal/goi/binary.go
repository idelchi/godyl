package goi

import (
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"sync"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/go-getter/v2"

	"github.com/idelchi/godyl/internal/data"
	"github.com/idelchi/godyl/internal/debug"
	"github.com/idelchi/godyl/internal/tools/checksum"
	"github.com/idelchi/godyl/pkg/download"
	"github.com/idelchi/godyl/pkg/path/file"
	"github.com/idelchi/godyl/pkg/path/folder"
)

// Binary represents a Go binary, including its associated file, directory, and environment variables.
type Binary struct {
	Env              Env
	File             file.File
	Dir              folder.Folder
	noVerifySSL      bool
	noVerifyChecksum bool
	progress         getter.ProgressTracker
}

// mutex is a mutex to prevent concurrent binary creation.
var mutex sync.Mutex //nolint:gochecknoglobals 		// TODO(Idelchi): Address this later.

// New creates a new Binary instance, setting up the directory, downloading the latest release if necessary,
// and initializing environment variables. It ensures thread-safe execution by using a mutex lock.
func New(
	noVerifySSL, downloadIfMissing, noVerifyChecksum bool,
	progress getter.ProgressTracker,
) (binary Binary, err error) {
	mutex.Lock()
	defer mutex.Unlock()

	binary.noVerifySSL = noVerifySSL
	binary.progress = progress
	binary.noVerifyChecksum = noVerifyChecksum

	// 1: Search for go binary on system
	if path, err := exec.LookPath("go"); err == nil {
		binary.File = file.New(path)
		binary.Env = Env{}
		binary.Dir = folder.New(binary.File.Dir())

		return binary, nil
		// 2: Else search in the (possibly) previously created directory
	} else if !downloadIfMissing {
		return binary, fmt.Errorf("looking for go binary in system: %w", err)
	}

	if path, err := binary.Find(data.GoDir().Path()); err == nil {
		binary.File = path
		binary.Dir = folder.New(binary.File.Dir())
		binary.Env.Default(data.GoDir().Path())

		return binary, nil
	}

	binary.Dir = data.GoDir()
	if err := binary.Dir.Create(); err != nil {
		return binary, fmt.Errorf("creating dir: %w", err)
	}

	release, err := binary.Latest()
	if err != nil {
		return binary, err
	}

	targets := Targets{}

	for _, file := range release.Files {
		if file.IsArchive() {
			targets.Files = append(targets.Files, file)
		}
	}

	path, err := targets.Match()
	if err != nil {
		return binary, err
	}

	toolchain := path[0].Asset.Name

	// Filter the targets to get back the original one
	targets = targets.FilterBy(func(t Target) bool {
		return t.FileName == toolchain
	})

	target := targets.Files[0]

	debug.Debug("Downloading Go toolchain: %q", target.FileName)

	err = binary.Download(target)
	if err != nil {
		return binary, err
	}

	debug.Debug("Go toolchain downloaded to: %q", binary.File)

	binary.Env.Default(binary.Dir.Path())

	return binary, nil
}

// Find searches for the Go binary in the given paths, returning the file if found.
func (b *Binary) Find(paths ...string) (file.File, error) {
	for _, path := range paths {
		file := file.New(path, "go", "bin", "go")

		if file.Exists() {
			return file, nil
		}
	}

	return file.New(""), fmt.Errorf("go binary not found: %w", folder.ErrNotFound)
}

// Download downloads the Go binary from the provided path and saves it to the directory.
// It returns an error if the download or file validation fails.
func (b *Binary) Download(target Target) error {
	url := "https://go.dev/dl/" + target.FileName

	options := []download.Option{}

	if b.noVerifySSL {
		options = append(options, download.WithInsecureSkipVerify())
	}

	if b.progress != nil {
		options = append(options, download.WithProgress(b.progress))
	}

	if !b.noVerifyChecksum {
		checksum := checksum.Checksum{
			Type:  "sha256",
			Value: target.Checksum,
		}

		options = append(
			options,
			download.WithChecksum(checksum.ToQuery()),
		)
	}

	downloader := download.New(options...)

	destination, err := downloader.Download(url, b.Dir.Path())
	if err != nil {
		return fmt.Errorf("downloading %q: %w", url, err)
	}

	b.File = file.New(destination.String(), "go", "bin", "go")

	return nil
}

// CleanUp removes the temporary directory associated with the binary.
// It returns an error if the directory removal fails.
func (b *Binary) CleanUp() error {
	if !b.Dir.IsSet() {
		return nil
	}

	if err := b.Dir.Remove(); err != nil {
		return fmt.Errorf("cleaning up: %w", err)
	}

	b.Dir = folder.New("")

	return nil
}

// Latest fetches the latest Go release information from the official Go download page.
// It returns the most recent release or an error if the process fails.
func (b *Binary) Latest() (Release, error) {
	client := resty.New()

	resp, err := client.R().Get("https://go.dev/dl/?mode=json")
	if err != nil {
		return Release{}, fmt.Errorf("fetching latest Go release: %w", err)
	}

	var releases []Release

	if err := json.Unmarshal(resp.Body(), &releases); err != nil {
		return Release{}, fmt.Errorf("unmarshalling Go releases: %w", err)
	}

	if len(releases) > 0 {
		return releases[0], nil
	}

	return Release{}, errors.New("no versions found")
}
