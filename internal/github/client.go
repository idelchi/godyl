package github

import (
	"net/url"

	"github.com/google/go-github/v74/github"
)

// NewClient creates a new GitHub client.
// If a token is provided, the client is authenticated using the token.
// Otherwise, an unauthenticated client is returned.
// An optional baseURL may be provided to redirect API requests to a custom endpoint
// (useful for testing with httptest servers).
func NewClient(token string, baseURL ...string) *github.Client {
	c := github.NewClient(nil)

	if token != "" {
		c = c.WithAuthToken(token)
	}

	if len(baseURL) > 0 && baseURL[0] != "" {
		u, err := url.Parse(baseURL[0])
		if err == nil {
			c.BaseURL = u
		}
	}

	return c
}
