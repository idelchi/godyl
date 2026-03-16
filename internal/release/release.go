// Package release provides shared types for representing releases and their assets
// across different source providers (GitHub, GitLab, etc.).
package release

import "errors"

// ErrRelease is returned when a release issue is encountered.
var ErrRelease = errors.New("release")

// Release represents a source release, containing the release name, tag, body, and associated assets.
type Release struct {
	Name   string
	Tag    string
	Body   string
	Assets Assets
}
