// Package url provides functionality to handle URLs as sources for downloading and managing files.
// It supports initializing URLs, retrieving metadata, and performing installations based on
// download operations. This package integrates with external utilities for file handling and
// matching patterns during installation processes.
package url

import (
	"fmt"
	"net/http"

	"github.com/idelchi/godyl/internal/match"
	"github.com/idelchi/godyl/internal/tools/sources/common"
	"github.com/idelchi/godyl/pkg/path/file"
)

// URL represents a download source with optional headers and tokens for authorization.
type URL struct {
	Token   Token
	Headers http.Header

	// Data holds additional metadata related to the URL.
	Data common.Metadata `yaml:"-"`
}

type Token struct {
	Token  string `mapstructure:"url-token" mask:"fixed"`
	Header string `mapstructure:"url-token-header"`
	Scheme string `mapstructure:"url-token-scheme"`
}

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

// Get retrieves a specific attribute from the URL's metadata.
func (u *URL) Get(attribute string) string {
	return u.Data.Get(attribute)
}

// Initialize prepares the URL based on the given name. (Ineffective).
func (u *URL) Initialize(_ string) error {
	return nil
}

// Exe executes the URL's associated action. (Ineffective).
func (u *URL) Exe() error {
	return nil
}

// Version sets the version for the URL. (Ineffective).
func (u *URL) Version(_ string) error {
	return nil
}

// Path sets the path for the URL, storing the provided name in the metadata.
func (u *URL) Path(name string, _ []string, _ string, _ match.Requirements) error {
	u.Data.Set("path", name)

	return nil
}

// Install downloads the file from the URL and processes it based on the provided InstallData.
// It returns the output, the downloaded file, and any error encountered.
func (u *URL) Install(d common.InstallData) (output string, found file.File, err error) {
	d.Header = u.GetHeaders()

	return common.Download(d)
}
