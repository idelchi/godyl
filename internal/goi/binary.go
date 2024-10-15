package goi

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"sync"

	"github.com/go-resty/resty/v2"

	"github.com/idelchi/godyl/pkg/download"
	"github.com/idelchi/godyl/pkg/file"
	"github.com/idelchi/godyl/pkg/folder"
)

type Binary struct {
	File file.File
	Dir  folder.Folder
	Env  Env
}

var mu sync.Mutex

func New() (binary Binary, err error) {
	mu.Lock()
	defer mu.Unlock()

	dir := folder.New(".godyl-go")
	if err := dir.CreateInTempDir(); err != nil && !errors.Is(err, os.ErrExist) {
		return binary, fmt.Errorf("creating temp dir: %w", err)
	}

	if file, err := binary.Find(dir.Path()); err == nil {
		binary.File = file
		if dir.IsParentOf(file.Dir()) {
			binary.Dir = dir
			binary.Env.Default(binary.Dir.Path())
		} else {
			binary.Env = Env{}
			binary.Dir = folder.New()
		}
		return binary, nil
	} else {
		binary.Dir = dir
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

	err = binary.Download(path[0].Name)
	if err != nil {
		return binary, err
	}

	binary.Env.Default(binary.Dir.Path())

	return binary, nil
}

func (b *Binary) Find(paths ...string) (file.File, error) {
	binary, err := exec.LookPath("go")
	if err != nil {
		for _, path := range paths {
			file := file.New(path, "go", "bin", "go")

			if file.Exists() {
				return file, nil
			}
		}

		return file.File(""), fmt.Errorf("go binary not found: %w", err)

	}

	return file.New(binary), nil
}

func (b *Binary) Download(path string) error {
	url := fmt.Sprintf("https://go.dev/dl/%s", path)

	fmt.Fprintf(os.Stderr, "Downloading %q\n", url)

	downloader := download.New()

	destination, err := downloader.Download(url, b.Dir.Path())
	if err != nil {
		return fmt.Errorf("downloading %q: %w", url, err)
	}

	if !destination.IsFile() {
		return fmt.Errorf("downloaded file is not a file")
	}

	b.File = file.New(destination.String(), "go", "bin", "go")

	return nil
}

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

func (b Binary) Latest() (Release, error) {
	client := resty.New()
	resp, err := client.R().Get("https://go.dev/dl/?mode=json")
	if err != nil {
		return Release{}, err
	}

	var releases []Release

	if err := json.Unmarshal(resp.Body(), &releases); err != nil {
		return Release{}, err
	}

	if len(releases) > 0 {
		return releases[0], nil
	}

	return Release{}, fmt.Errorf("no versions found")
}
