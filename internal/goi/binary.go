package goi

import (
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"sync"

	"github.com/go-resty/resty/v2"

	"github.com/idelchi/godyl/internal/tmp"
	"github.com/idelchi/godyl/pkg/download"
	"github.com/idelchi/godyl/pkg/file"
	"github.com/idelchi/godyl/pkg/folder"
)

// Binary represents a Go binary, including its associated file, directory, and environment variables.
type Binary struct {
	File file.File     // File holds the file information for the Go binary.
	Dir  folder.Folder // Dir refers to the directory where the binary is stored.
	Env  Env           // Env contains the environment variables for running the binary.

	noVerifySSL bool
}

// mutex is a mutex to prevent concurrent binary creation.
var mutex sync.Mutex //nolint:gochecknoglobals 		// TODO(Idelchi): Address this later.

// New creates a new Binary instance, setting up the directory, downloading the latest release if necessary,
// and initializing environment variables. It ensures thread-safe execution by using a mutex lock.
func New(noVerifySSL bool) (binary Binary, err error) {
	mutex.Lock()
	defer mutex.Unlock()

	binary.noVerifySSL = noVerifySSL

	// Step 1: Search for go binary on system
	if path, err := exec.LookPath("go"); err == nil {
		binary.File = file.New(path)
		binary.Env = Env{}
		binary.Dir = folder.New(binary.File.Dir())

		return binary, nil
		// Step 2: Else search in other possible paths
		// } else if path, err := binary.Find("/some", "/other", "/paths"); err == nil {
		// 	binary.File = path
		// 	binary.Env = Env{}
		// 	binary.Dir = folder.New(binary.File.Dir())

		// 	return binary, nil
		// Step 3: Else search in the (possibly) previously created directory
	} else if path, err := binary.Find(tmp.GodylDir("go").Path()); err == nil {
		binary.File = path
		binary.Dir = folder.New(binary.File.Dir())
		binary.Env.Default(binary.Dir.Path())

		return binary, nil
	}

	binary.Dir = tmp.GodylDir("go")
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

	err = binary.Download(path[0].Asset.Name)
	if err != nil {
		return binary, err
	}

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

	return file.File(""), fmt.Errorf("go binary not found: %w", folder.ErrNotFound)
}

// Download downloads the Go binary from the provided path and saves it to the directory.
// It returns an error if the download or file validation fails.
func (b *Binary) Download(path string) error {
	url := "https://go.dev/dl/" + path

	downloader := download.New()
	downloader.InsecureSkipVerify = b.noVerifySSL

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

	b.Dir = ""

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
