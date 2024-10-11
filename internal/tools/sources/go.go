package sources

import (
	"fmt"
	"path/filepath"
	"strings"
	"sync"

	"github.com/idelchi/godyl/internal/folder"
	ginstaller "github.com/idelchi/godyl/internal/go"
	"github.com/idelchi/godyl/internal/match"
)

type Go struct {
	github *GitHub

	Command string `yaml:"command"`

	Data Metadata `yaml:"-"`
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

func (g *Go) Install(d InstallData) (output, found string, err error) {
	mu.Lock()
	binary, err := ginstaller.New()
	if err != nil {
		return "", "", err
	}
	mu.Unlock()

	installer := ginstaller.GInstaller{
		Binary: binary,
	}

	var folder folder.Folder
	folder.CreateRandomInTempDir()

	installer.Binary.Env.Append(
		ginstaller.Env{
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
		paths = []string{(strings.Replace(d.Path, fmt.Sprintf("/%s@", name), fmt.Sprintf("/%s/%s@", name, g.Command), 1))}
	}

	for i, path := range paths {
		// Lowercase the path
		paths[i] = strings.ToLower(path)
	}

	for _, path := range paths {
		output, err = installer.Install(path)

		if err == nil {
			d.Path = path
			found, err := FindAndSymlink(folder.Path(), d)

			return output, found, err
		} else {
			fmt.Println(err)
			fmt.Println(output)
		}
	}

	return output, "", err
}
