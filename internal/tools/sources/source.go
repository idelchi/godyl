package sources

import (
	"fmt"

	"github.com/idelchi/godyl/internal/match"
	"github.com/idelchi/godyl/internal/tools/sources/command"
	"github.com/idelchi/godyl/internal/tools/sources/common"
	"github.com/idelchi/godyl/internal/tools/sources/github"
	goc "github.com/idelchi/godyl/internal/tools/sources/go"
	"github.com/idelchi/godyl/internal/tools/sources/url"
	"github.com/idelchi/godyl/pkg/file"
)

type Type string

func (t Type) String() string {
	return string(t)
}

func (t *Type) From(name string) {
	*t = Type(name)
}

const (
	GITHUB  Type = "github"
	GITLAB  Type = "gitlab"
	DIRECT  Type = "url"
	COMMAND Type = "command"
	GO      Type = "go"
	RUST    Type = "rust"
)

type Source struct {
	Type     Type
	Github   github.GitHub
	URL      url.URL
	Go       goc.Go
	Commands command.Commands
}

type Populater interface {
	Initialize(string) error
	Exe() error
	Version(string) error
	Path(string, []string, string, match.Requirements) error
	Install(common.InstallData) (string, file.File, error)
	Get(string) string
}

func (s *Source) Installer() (Populater, error) {
	switch s.Type {
	case GITHUB:
		return &s.Github, nil
	case DIRECT:
		return &s.URL, nil
	case COMMAND:
		return &s.Commands, nil
	case GO:
		s.Go.SetGitHub(&s.Github)
		return &s.Go, nil
	default:
		return nil, fmt.Errorf("unknown source type: %s", s.Type)
	}
}
