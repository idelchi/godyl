package sources

import (
	"fmt"

	"github.com/idelchi/godyl/internal/match"
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
	Github   GitHub
	URL      URL
	Go       Go
	Commands Commands
}

type Populater interface {
	Initialize(string) error
	Exe() error
	Version(string) error
	Path(string, []string, string, match.Requirements) error
	Install(InstallData) (string, file.File, error)
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
		s.Go.github = &s.Github
		return &s.Go, nil
	default:
		return nil, fmt.Errorf("unknown source type: %s", s.Type)
	}
}
