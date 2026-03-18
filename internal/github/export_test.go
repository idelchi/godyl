package github

import "net/http"

// ParseGitHubReleaseAssets exports the unexported parseGitHubReleaseAssets for use in tests.
var ParseGitHubReleaseAssets = parseGitHubReleaseAssets

// SetTransport sets the HTTP transport on a Repository for testing.
func SetTransport(r *Repository, rt http.RoundTripper) {
	r.transport = rt
}
