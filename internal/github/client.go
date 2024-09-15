package github

import (
	"github.com/google/go-github/v64/github"
)

// NewClient creates a new GitHub client.
// If a token is provided, the client is authenticated using the token.
// Otherwise, an unauthenticated client is returned.
func NewClient(token string) *github.Client {
	c := github.NewClient(nil)
	if token != "" {
		return c.WithAuthToken(token) // Authenticate the client with the provided token.
	}

	return c // Return unauthenticated client if no token is provided.
}
