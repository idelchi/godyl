// Package release provides shared types for representing releases and their assets
// across different source providers (GitHub, GitLab, etc.).
package release

import "errors"

// ErrRelease is returned when a release issue is encountered.
var ErrRelease = errors.New("release")

// Release represents a source release, containing the release name, tag, body, and associated assets.
type Release struct {
	// Name is the human-readable release name.
	Name string
	// Tag identifies the release version in its source repository.
	Tag string
	// Body contains the release notes.
	Body string
	// Assets lists the files published with the release.
	Assets Assets
}
