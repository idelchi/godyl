package sources

import (
	"fmt"

	"github.com/idelchi/godyl/internal/match"
)

type Source struct {
	Type     string
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
	Install(InstallData) (string, error)
	Get(string) string
}

func (s *Source) Installer() (Populater, error) {
	switch s.Type {
	case "github":
		return &s.Github, nil
	case "url":
		return &s.URL, nil
	case "command":
		return &s.Commands, nil
	case "go":
		return &Go{github: &s.Github}, nil
	default:
		return nil, fmt.Errorf("unknown source type: %s", s.Type)
	}
}
