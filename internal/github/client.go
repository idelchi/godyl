package github

import (
	"github.com/google/go-github/v64/github"
)

func NewClient(token string) *github.Client {
	c := github.NewClient(nil)
	if token != "" {
		return c.WithAuthToken(token)
	}

	return c
}
