// Package url provides functionality to handle URLs as sources for downloading and managing files.
// It supports initializing URLs, retrieving metadata, and performing installations based on
// download operations. This package integrates with external utilities for file handling and
// matching patterns during installation processes.
package url

import (
	"fmt"
	"net/http"

	"github.com/hashicorp/go-getter/v2"

	"github.com/idelchi/godyl/internal/match"
	"github.com/idelchi/godyl/internal/tools/sources/common"
	"github.com/idelchi/godyl/pkg/path/file"
)

// URL represents a URL-based download source configuration.
type URL struct {
	Headers http.Header
	Data    common.Metadata `yaml:"-"`
	Token   Token
}

// Token contains authentication configuration for URL requests.
type Token struct {
	// Token is the authentication token value.
	Token string

	// Header is the HTTP header name for the token.
	Header string

	// Scheme is the authentication scheme (e.g., "Bearer").
	Scheme string
}

// Initialize is a no-op implementation of the Populator interface.
func (u *URL) Initialize(_ string) error {
	return nil
}

// Version is a no-op implementation of the Populator interface.
func (u *URL) Version(_ string) error {
	return nil
}

// Path stores the provided URL in the metadata.
// The URL will be used as the download source during installation.
func (u *URL) URL(name string, _ []string, _ string, _ match.Requirements) error {
	u.Data.Set("url", name)

	return nil
}

// Install downloads a file from the configured URL.
// Handles authentication, downloads the file, and processes it according to InstallData.
// Returns the operation output, downloaded file information, and any errors.
func (u *URL) Install(
	d common.InstallData,
	progressListener getter.ProgressTracker,
) (output string, found file.File, err error) {
	d.Header = u.GetHeaders()
	// Pass the progress listener down
	d.ProgressListener = progressListener

	return common.Download(d)
}

// Get retrieves a metadata attribute value by its key.
func (u *URL) Get(attribute string) string {
	return u.Data.Get(attribute)
}

// GetHeaders returns HTTP headers for URL requests, including authentication.
// Combines configured headers with token-based authentication if configured.
func (u *URL) GetHeaders() http.Header {
	headers := make(http.Header)

	// Clone existing headers if any
	if u.Headers != nil {
		headers = u.Headers.Clone()
	}

	// Add token to headers if both token and header are specified
	if u.Token.Token != "" && u.Token.Header != "" {
		tokenValue := u.Token.Token
		if u.Token.Scheme != "" {
			tokenValue = fmt.Sprintf("%s %s", u.Token.Scheme, u.Token.Token)
		}

		headers.Set(u.Token.Header, tokenValue)
	}

	return headers
}
