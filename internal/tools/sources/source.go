// Package sources provides abstractions for handling various types of installation sources,
// including GitHub repositories, direct URLs, Go projects, and command-based sources.
// The package defines a common interface, Populater, which is implemented by these sources
// to handle initialization, execution, versioning, path setup, and installation processes.
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

// Type represents the source type, such as GitHub, URL, Go, or command-based sources.
type Type string

// String returns the string representation of the Type.
func (t Type) String() string {
	return string(t)
}

// From sets the Type from the provided name.
func (t *Type) From(name string) {
	*t = Type(name)
}

const (
	GITHUB  Type = "github"  // GitHub source type
	GITLAB  Type = "gitlab"  // GitLab source type
	DIRECT  Type = "url"     // URL source type
	COMMAND Type = "command" // Command-based source type
	GO      Type = "go"      // Go source type
	RUST    Type = "rust"    // Rust source type
)

// Source represents a source of installation, which could be GitHub, URL, Go, or command-based.
type Source struct {
	Type     Type             // Type of the source
	Github   github.GitHub    // GitHub repository source
	URL      url.URL          // URL source for direct downloads
	Go       goc.Go           // Go project source
	Commands command.Commands // Command-based source
}

// Populater defines the interface that all source types must implement to handle initialization, execution,
// versioning, path setup, and installation.
type Populater interface {
	Initialize(string) error
	Exe() error
	Version(string) error
	Path(string, []string, string, match.Requirements) error
	Install(common.InstallData) (string, file.File, error)
	Get(string) string
}

// Installer returns the appropriate Populater implementation based on the source Type.
// It determines the correct handling for GitHub, URL, Go, and command-based sources.
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
