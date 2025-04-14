package github

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const gitHubLatestReleaseURLFormat = "https://github.com/%s/%s/releases/latest"

// WebReleaseInfo stores information about a release fetched from the GitHub web interface.
type WebReleaseInfo struct {
	Tag string
	URL string
}

// newHTTPClient returns a new HTTP client with reasonable timeouts.
func newHTTPClient() *http.Client {
	return &http.Client{
		Timeout: 10 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// Don't follow redirects, we just want the Location header
			return http.ErrUseLastResponse
		},
	}
}

// LatestReleaseFromWeb retrieves the latest release for the repository using the GitHub website
// instead of the API, avoiding rate limits.
func (g *Repository) LatestVersionFromWeb() (string, error) {
	webReleaseInfo, err := g.getLatestReleaseInfoFromWeb(g.ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get latest release info from web: %w", err)
	}

	// Now that we have the tag, we can use the existing method to get the full release details
	// If you want to completely avoid the API, you would need to parse the HTML of the release page
	return webReleaseInfo.Tag, nil
}

// getLatestReleaseInfoFromWeb gets the latest release tag by making a HEAD request to the GitHub releases page.
func (g *Repository) getLatestReleaseInfoFromWeb(ctx context.Context) (*WebReleaseInfo, error) {
	url := fmt.Sprintf(gitHubLatestReleaseURLFormat, g.Owner, g.Repo)

	req, err := http.NewRequestWithContext(ctx, http.MethodHead, url, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	client := newHTTPClient()
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusFound {
		return nil, fmt.Errorf("incorrect status code: %d", res.StatusCode)
	}

	loc := res.Header.Get("Location")
	if loc == "" {
		return nil, fmt.Errorf("unable to determine release version (empty Location header)")
	}

	// Extract the tag from the Location header
	// The URL format is typically: https://github.com/owner/repo/releases/tag/v1.2.3
	parts := strings.Split(loc, "/")
	if len(parts) < 2 {
		return nil, fmt.Errorf("unable to parse release tag from URL: %s", loc)
	}
	tag := parts[len(parts)-1]

	return &WebReleaseInfo{
		Tag: tag,
		URL: loc,
	}, nil
}
