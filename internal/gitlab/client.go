package gitlab

import (
	"fmt"
	"net/url"

	gitlab "gitlab.com/gitlab-org/api/client-go"
)

// NewClient creates a new GitLab client.
// If a token is provided, the client is authenticated using the token.
// If baseURL is provided, the client will connect to that GitLab instance instead of gitlab.com.
func NewClient(token, baseURL string) (*gitlab.Client, error) {
	var options []gitlab.ClientOptionFunc

	// If baseURL is provided, configure the client to use it
	if baseURL != "" {
		const apiPath = "/api/v4"

		url, err := url.JoinPath(baseURL, apiPath)
		if err != nil {
			return nil, fmt.Errorf("joining base URL %q with API path %q: %w", url, apiPath, err)
		}

		options = append(options, gitlab.WithBaseURL(url))
	}

	// Create client with token and options
	client, err := gitlab.NewClient(token, options...)
	if err != nil {
		return nil, fmt.Errorf("creating GitLab client at %q: %w", baseURL, err)
	}

	return client, nil
}
