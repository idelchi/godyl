package github

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"

	"github.com/idelchi/godyl/pkg/path/file"
)

const gitHubLatestReleaseURLFormat = "https://github.com/%s/%s/releases/latest"

// WebReleaseInfo stores information about a release fetched from the GitHub web interface.
type WebReleaseInfo struct {
	Tag string `json:"tag_name"`
}

// newHTTPClient returns a new HTTP client with reasonable timeouts.
func newHTTPClient() *http.Client {
	const Timeout = 10 * time.Second

	return &http.Client{
		Timeout: Timeout,
		CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
			// Don't follow redirects, we just want the Location header
			return http.ErrUseLastResponse
		},
	}
}

// GetReleaseFromWeb retrieves a specific release for the repository based on the provided tag.
func (r *Repository) GetReleaseFromWeb(ctx context.Context, tag string) (*Release, error) {
	url := fmt.Sprintf("https://github.com/%s/%s/releases/expanded_assets/%s", r.Owner, r.Repo, tag)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set Accept header to get JSON-like response
	req.Header.Set("Accept", "application/json")

	client := newHTTPClient()

	client.CheckRedirect = nil // Allow redirects for this request

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("incorrect status code: %d", res.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse the HTML to extract assets
	assets, err := parseGitHubReleaseAssets(string(body))
	if err != nil {
		return nil, fmt.Errorf("failed to parse assets: %w", err)
	}

	// Create the Release object
	release := &Release{
		Tag:    tag,
		Name:   tag, // GitHub web interface doesn't provide release name in this endpoint
		Assets: assets,
		// Body is not available from this endpoint
	}

	return release, nil
}

// LatestVersionFromWebHTML retrieves the latest release for the repository using the GitHub website
// instead of the API, avoiding rate limits.
func (r *Repository) LatestVersionFromWebHTML(ctx context.Context) (string, error) {
	webReleaseInfo, err := r.getLatestReleaseFromWebHTML(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get latest release info from web: %w", err)
	}

	// Now that we have the tag, we can use the existing method to get the full release details
	// If you want to completely avoid the API, you would need to parse the HTML of the release page
	return webReleaseInfo.Tag, nil
}

// LatestVersionFromWebJSON retrieves the latest release for the repository using the GitHub website
// instead of the API, avoiding rate limits.
func (r *Repository) LatestVersionFromWebJSON(ctx context.Context) (string, error) {
	webReleaseInfo, err := r.getLatestReleaseInfoFromWebJSON(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get latest release info from web: %w", err)
	}

	// Now that we have the tag, we can use the existing method to get the full release details
	// If you want to completely avoid the API, you would need to parse the HTML of the release page
	return webReleaseInfo.Tag, nil
}

// getLatestReleaseInfoFromWeb gets the latest release tag by making a HEAD request to the GitHub releases page.
func (r *Repository) getLatestReleaseFromWebHTML(ctx context.Context) (*WebReleaseInfo, error) {
	url := fmt.Sprintf(gitHubLatestReleaseURLFormat, r.Owner, r.Repo)

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
		return nil, errors.New("unable to determine release version (empty Location header)")
	}

	// Extract the tag from the Location header
	// The URL format is typically: https://github.com/owner/repo/releases/tag/v1.2.3
	parts := strings.Split(loc, "/")

	const expectedParts = 2

	if len(parts) < expectedParts {
		return nil, fmt.Errorf("unable to parse release tag from URL: %q", loc)
	}

	tag := parts[len(parts)-1]

	return &WebReleaseInfo{
		Tag: tag,
	}, nil
}

func (r *Repository) getLatestReleaseInfoFromWebJSON(ctx context.Context) (*WebReleaseInfo, error) {
	url := fmt.Sprintf(gitHubLatestReleaseURLFormat, r.Owner, r.Repo)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")

	client := newHTTPClient()
	// Remove the CheckRedirect since we want to follow the redirect and get JSON
	client.CheckRedirect = nil

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("incorrect status code: %d", res.StatusCode)
	}

	release := &WebReleaseInfo{}

	if err := json.NewDecoder(res.Body).Decode(&release); err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}

	if release.Tag == "" {
		return nil, errors.New("tag_name is empty in response")
	}

	return release, nil
}

func parseGitHubReleaseAssets(html string) ([]Asset, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	var assets []Asset

	// Find all list items in the assets box
	doc.Find("li.Box-row").Each(func(_ int, s *goquery.Selection) {
		// Find the download link
		link := s.Find("a[href*='/releases/download/']")
		if link.Length() > 0 {
			href, exists := link.Attr("href")
			if exists {
				// Extract the filename from the link text
				name := file.New(href).Base()

				// Build the full URL
				url := "https://github.com" + href

				assets = append(assets, Asset{
					Name: name,
					URL:  url,
				})
			}
		}
	})

	return assets, nil
}
