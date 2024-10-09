package sources

import (
	"fmt"
	"sync"

	"github.com/idelchi/godyl/internal/folder"
	ginstaller "github.com/idelchi/godyl/internal/go"
	"github.com/idelchi/godyl/internal/match"
)

type Go struct {
	github *GitHub

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

func (g Go) Install(d InstallData) (output string, err error) {
	mu.Lock()
	binary, err := ginstaller.New()
	if err != nil {
		return "", err
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

	output, err = installer.Install(d.Path)
	if err != nil {
		return output, err
	}

	return output, FindAndSymlink(folder.Path(), d)
}
