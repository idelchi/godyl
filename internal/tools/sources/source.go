// Package sources provides abstractions for handling various types of installation sources,
// including GitHub repositories, direct URLs, Go projects, and command-based sources.
// The package defines a common interface, Populator, which is implemented by these sources
// to handle initialization, execution, versioning, path setup, and installation processes.
package sources

import (
	"fmt"
	"net/http"

	"github.com/hashicorp/go-getter/v2"

	"github.com/idelchi/godyl/internal/match"
	"github.com/idelchi/godyl/internal/tools/sources/github"
	"github.com/idelchi/godyl/internal/tools/sources/gitlab"
	goc "github.com/idelchi/godyl/internal/tools/sources/go"
	"github.com/idelchi/godyl/internal/tools/sources/install"
	"github.com/idelchi/godyl/internal/tools/sources/none"
	"github.com/idelchi/godyl/internal/tools/sources/url"
	"github.com/idelchi/godyl/pkg/path/file"
)

// Type represents the installation source type for a tool.
type Type string

// String returns the string representation of the Type value.
func (t Type) String() string {
	return string(t)
}

// From sets the Type value from the provided name string.
func (t *Type) From(name string) {
	*t = Type(name)
}

// TODO(Idelchi): go generate the source type strings

const (
	// GITHUB indicates GitHub as the source type.
	GITHUB Type = "github"
	// GITLAB indicates GitLab as the source type.
	GITLAB Type = "gitlab"
	// URL indicates a direct URL as the source type.
	URL Type = "url"
	// NONE indicates no source type.
	NONE Type = "none"
	// GO indicates Go modules as the source type.
	GO Type = "go"
)

// SupportsChecksum returns true if the source type supports checksum verification.
func (t Type) SupportsChecksum() bool {
	switch t {
	case GITHUB, GITLAB, URL:
		return true
	case NONE, GO:
		return false
	default:
		return false
	}
}

// SupportsArtifactReuse returns true if source artifacts can be reused after URL resolution.
func (t Type) SupportsArtifactReuse() bool {
	switch t {
	case GITHUB, GITLAB, URL:
		return true
	case NONE, GO:
		return false
	default:
		return false
	}
}

// Source represents the configuration for various source types used to retrieve tools.
// TODO(Idelchi): Add validation.
type Source struct {
	// GitHub contains GitHub-specific source configuration.
	GitHub github.GitHub
	// URL contains direct-download source configuration.
	URL url.URL
	// Go contains Go-module source configuration.
	Go goc.Go
	// Type selects the active source implementation.
	Type Type `validate:"oneof=github gitlab url none go"`
	// GitLab contains GitLab-specific source configuration.
	GitLab gitlab.GitLab
}

// Headers returns the HTTP headers used by the active source during artifact download.
func (s *Source) Headers() http.Header {
	switch s.Type {
	case GITLAB:
		return s.GitLab.GetHeaders().Clone()
	case URL:
		return s.URL.Headers.Clone()
	case GITHUB, NONE, GO:
		return http.Header{}
	default:
		return http.Header{}
	}
}

// Populator defines the interface that all source types must implement.
// It provides methods for managing the complete lifecycle of tool installation,
// from initialization through execution, versioning, path setup, and installation.
type Populator interface {
	// Initialize derives source state from the configured repository or URL.
	Initialize(repo string) error
	// Version resolves and stores the target version.
	Version(version string) error
	// URL resolves and stores the artifact location.
	URL(name string, extensions []string, version string, requirements match.Requirements) error
	// Install retrieves and installs the resolved artifact.
	Install(data install.Data, progressListener getter.ProgressTracker) (string, file.File, error)
	// Get returns a resolved metadata value.
	Get(key string) string
}

// Installer returns the appropriate Populator implementation for the source Type.
// Returns an error if the source type is unknown or unsupported.
func (s *Source) Installer() (Populator, error) {
	switch s.Type {
	case GITHUB:
		return &s.GitHub, nil
	case GITLAB:
		return &s.GitLab, nil
	case URL:
		return &s.URL, nil
	case NONE:
		return &none.None{}, nil
	case GO:
		s.Go.SetGitHub(&s.GitHub)

		return &s.Go, nil
	default:
		return nil, fmt.Errorf("unknown source type: %s", s.Type)
	}
}
