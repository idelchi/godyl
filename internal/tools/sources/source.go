// Package sources provides abstractions for handling various types of installation sources,
// including GitHub repositories, direct URLs, Go projects, and command-based sources.
// The package defines a common interface, Populater, which is implemented by these sources
// to handle initialization, execution, versioning, path setup, and installation processes.
package sources

import (
	"fmt"

	"github.com/idelchi/godyl/internal/match"
	"github.com/idelchi/godyl/internal/tools/sources/common"
	"github.com/idelchi/godyl/internal/tools/sources/github"
	"github.com/idelchi/godyl/internal/tools/sources/gitlab"
	goc "github.com/idelchi/godyl/internal/tools/sources/go"
	"github.com/idelchi/godyl/internal/tools/sources/none"
	"github.com/idelchi/godyl/internal/tools/sources/url"
	"github.com/idelchi/godyl/pkg/path/file"
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

// TODO(Idelchi): go generate the source type strings

const (
	GITHUB Type = "github" // GitHub source type
	GITLAB Type = "gitlab" // GitLab source type
	URL    Type = "url"    // URL source type
	NONE   Type = "none"   // No source type
	GO     Type = "go"     // Go source type
	RUST   Type = "rust"   // Rust source type
)

// Source represents a source of installation, which could be GitHub, URL, Go, or command-based.
//
// TODO(Idelchi): Add validation
type Source struct {
	// Type of the source
	Type Type `validate:"oneof=github gitlab url none go rust"`
	// GitHub repository source
	GitHub github.GitHub
	// GitLab repository source
	GitLab gitlab.GitLab
	// URL source for direct downloads
	URL url.URL
	// Go project source
	Go goc.Go
	// None
	None none.None
}

// Populater defines the interface that all source types must implement to handle initialization, execution,
// versioning, path setup, and installation.
type Populater interface {
	Initialize(repo string) error
	Exe() error
	Version(version string) error
	Path(name string, extensions []string, version string, requirements match.Requirements) error
	Install(data common.InstallData) (string, file.File, error)
	Get(key string) string
}

// Installer returns the appropriate Populater implementation based on the source Type.
// It determines the correct handling for GitHub, URL, Go, and command-based sources.
func (s *Source) Installer() (Populater, error) {
	switch s.Type {
	case GITHUB:
		return &s.GitHub, nil
	case GITLAB:
		return &s.GitLab, nil
	case URL:
		return &s.URL, nil
	case NONE:
		return &s.None, nil
	case GO:
		s.Go.SetGitHub(&s.GitHub)

		return &s.Go, nil
	case RUST:
		return nil, fmt.Errorf("source type %s is not yet supported", s.Type)
	default:
		return nil, fmt.Errorf("unknown source type: %s", s.Type)
	}
}
