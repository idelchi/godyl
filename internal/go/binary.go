package ginstaller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/idelchi/godyl/internal/detect"
	"github.com/idelchi/godyl/internal/executable"
	"github.com/idelchi/godyl/internal/folder"
	"github.com/idelchi/godyl/internal/match"
	"github.com/idelchi/godyl/pkg/download"
)

type Env map[string]string

func (e Env) ToSlice() []string {
	var env []string
	for k, v := range e {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}

	return env
}

func (e *Env) Append(env Env) {
	for k, v := range env {
		(*e)[k] = v
	}
}

type Binary struct {
	Path executable.Executable
	Dir  folder.Folder
	Env  Env
}

var mu sync.Mutex

func New() (b Binary, err error) {
	mu.Lock()
	defer mu.Unlock()

	dir := folder.Folder(".godyl-go")
	if err := dir.CreateInTempDir(); err != nil && !errors.Is(err, os.ErrExist) {
		return b, fmt.Errorf("creating temp dir: %w", err)
	}

	if binary, err := b.Find(dir.Path()); err == nil {
		b.Path = binary
		b.Dir = folder.Folder(binary.Dir())

		if dir.IsParentOf(b.Dir) {
			fmt.Printf("%q is parent of %q\n", dir.Path(), b.Dir.Path())
			b.Env = Env{
				"GOMODCACHE": filepath.Join(b.Dir.Path(), ".cache"),
				"GOCACHE":    filepath.Join(b.Dir.Path(), ".cache"),
				"GOPATH":     filepath.Join(b.Dir.Path(), ".path"),
			}
		} else {
			b.Env = Env{}
		}
		return b, nil
	} else {
		b.Dir = dir
	}

	release, err := b.Latest()
	if err != nil {
		return b, err
	}

	targets := Targets{}
	for _, file := range release.Files {
		if file.IsArchive() {
			targets.Files = append(targets.Files, file)
		}
	}

	path, err := targets.Match()
	if err != nil {
		return b, err
	}

	err = b.Download(path[0].Name)
	if err != nil {
		return b, err
	}

	b.Env = Env{
		"GOMODCACHE": filepath.Join(b.Dir.Path(), ".cache"),
		"GOCACHE":    filepath.Join(b.Dir.Path(), ".cache"),
		"GOPATH":     filepath.Join(b.Dir.Path(), ".path"),
	}

	return b, nil
}

func (b *Binary) Find(paths ...string) (executable.Executable, error) {
	binary, err := exec.LookPath("go")
	err = errors.New("not found")
	if err != nil {
		for _, path := range paths {
			file := executable.New(path, filepath.Join("go", "bin", "go"))

			if file.Exists() {
				return file, nil
			}
		}

		return executable.Executable{}, fmt.Errorf("go binary not found: %w", err)

	} else {
		fmt.Println(b)
	}

	return executable.Executable{
		Path: binary,
	}, nil
}

func (b *Binary) Download(path string) error {
	url := fmt.Sprintf("https://go.dev/dl/%s", path)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	destination, err := download.Download(ctx, url, b.Dir.Path())
	if err != nil {
		return fmt.Errorf("downloading %q: %w", url, err)
	}

	b.Path = executable.New(destination, filepath.Join("go", "bin", "go"))

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

type Targets struct {
	Files []Target `json:"files"`
}

func (gt Targets) FilterBy(predicate func(Target) bool) Targets {
	var filtered Targets
	for _, file := range gt.Files {
		if predicate(file) {
			filtered.Files = append(filtered.Files, file)
		}
	}
	return filtered
}

func (gt Targets) FilterByOS(os string) Targets {
	return gt.FilterBy(func(file Target) bool {
		return file.OS == os
	})
}

func (gt Targets) FilterByArch(arch string) Targets {
	return gt.FilterBy(func(file Target) bool {
		return file.Arch == arch
	})
}

type Target struct {
	FileName string `json:"filename"`
	Arch     string `json:"arch"`
	OS       string `json:"os"`
	Version  string `json:"version"`
}

func (t Target) IsArchive() bool {
	return strings.HasSuffix(t.FileName, ".tar.gz") || filepath.Ext(t.FileName) == ".zip"
}

type Release struct {
	Version string   `json:"version"`
	Files   []Target `json:"files"`
}

func (t Targets) Match() (match.Results, error) {
	platform := detect.Platform{}
	if err := platform.Detect(); err != nil {
		return nil, fmt.Errorf("detecting platform: %w", err)
	}

	var assets match.Assets

	for _, tt := range t.Files {
		asset := match.Asset{Name: tt.FileName}

		asset.Platfrom.OS.From(tt.OS)
		asset.Platfrom.Architecture.From(tt.Arch, "")

		assets = append(assets, asset)
	}

	hints := []match.Hint{
		match.NewDefaultHint(platform.OS.Name()),
	}

	var err error

	matches := assets.Select(match.Requirements{Platform: platform, Hints: hints})
	switch {
	case !matches.HasQualified():
		err = fmt.Errorf("no qualified file found")
	case matches.IsAmbigious():
		err = fmt.Errorf("ambiguous file selection")
	case !matches.Success():
		err = fmt.Errorf("no matching file found")
	}

	return matches, err
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
