package goc

import (
	"fmt"
	"path/filepath"
	"strings"
	"sync"

	"github.com/idelchi/godyl/internal/goi"
	"github.com/idelchi/godyl/internal/match"
	"github.com/idelchi/godyl/internal/tools/sources/common"
	"github.com/idelchi/godyl/internal/tools/sources/github"
	"github.com/idelchi/godyl/pkg/file"
	"github.com/idelchi/godyl/pkg/folder"
)

type Go struct {
	github *github.GitHub

	Command string `yaml:"command"`

	Data common.Metadata `yaml:"-"`
}

func (g *Go) SetGitHub(gh *github.GitHub) {
	g.github = gh
}

func (g *Go) Get(attribute string) string {
	return g.github.Data.Get(attribute)
}

func (g *Go) Initialize(name string) error {
	return g.github.Initialize(name)
}

func (g *Go) Exe() error {
	return g.github.Exe()
}

func (g *Go) Version(name string) error {
	return g.github.Version(name)
}

func (g *Go) Path(_ string, _ []string, version string, _ match.Requirements) error {
	g.github.Data.Set("path", fmt.Sprintf("github.com/%s/%s@%s", g.github.Owner, g.github.Repo, version))

	return nil
}

var mu sync.Mutex

func (g *Go) Install(d common.InstallData) (output string, found file.File, err error) {
	mu.Lock()
	binary, err := goi.New()
	if err != nil {
		return "", "", err
	}
	mu.Unlock()

	installer := goi.Installer{
		Binary: binary,
	}

	var folder folder.Folder
	folder.CreateRandomInTempDir()

	installer.Binary.Env.Append(
		goi.Env{
			"GOBIN": folder.Path(),
		},
	)

	name := strings.TrimSuffix(d.Exe, filepath.Ext(d.Exe))

	paths := []string{
		d.Path,
		strings.Replace(d.Path, fmt.Sprintf("/%s@", name), fmt.Sprintf("/%s/cmd/%s@", name, name), 1),
		strings.Replace(d.Path, fmt.Sprintf("/%s@", name), fmt.Sprintf("/%s/cmd@", name), 1),
	}

	if g.Command != "" {
		paths = []string{
			(strings.Replace(d.Path, fmt.Sprintf("/%s@", name), fmt.Sprintf("/%s/%s@", name, g.Command), 1)),
		}
	}

	for i, path := range paths {
		// Lowercase the path
		paths[i] = strings.ToLower(path)
	}

	for _, path := range paths {
		output, err = installer.Install(path)

		if err == nil {
			d.Path = path
			found, err := common.FindAndSymlink(file.New(folder.Path()), d)

			return output, found, err
		} else {
			fmt.Println(err)
			fmt.Println(output)
		}
	}

	return output, "", err
}
