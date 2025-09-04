// Package url provides functionality to handle URLs as sources for downloading and managing files.
// It supports initializing URLs, retrieving metadata, and performing installations based on
// download operations. This package integrates with external utilities for file handling and
// matching patterns during installation processes.
package url

import (
	"net/http"

	"github.com/hashicorp/go-getter/v2"

	"github.com/idelchi/godyl/internal/match"
	"github.com/idelchi/godyl/internal/tools/sources/install"
	"github.com/idelchi/godyl/pkg/path/file"
)

// URL represents a URL-based download source configuration.
type URL struct {
	Headers http.Header      `mapstructure:"headers" yaml:"headers"`
	Data    install.Metadata `mapstructure:"-"       yaml:"-"`
	Token   string           `mapstructure:"token"   yaml:"token"   mask:"fixed"`
}

// Initialize is a no-op implementation of the Populator interface.
func (u *URL) Initialize(_ string) error {
	return nil
}

// Version is a no-op implementation of the Populator interface.
func (u *URL) Version(_ string) error {
	return nil
}

// URL stores the provided URL in the metadata.
// The URL will be used as the download source during installation.
func (u *URL) URL(name string, _ []string, _ string, _ match.Requirements) error {
	u.Data.Set("url", name)

	return nil
}

// Install downloads a file from the configured URL.
// Handles authentication, downloads the file, and processes it according to InstallData.
// Returns the operation output, downloaded file information, and any errors.
func (u *URL) Install(
	d install.Data,
	progressListener getter.ProgressTracker,
) (output string, found file.File, err error) {
	d.Header = u.Headers
	// Pass the progress listener down
	d.ProgressListener = progressListener

	found, err = install.Download(d)

	return "", found, err
}

// Get retrieves a metadata attribute value by its key.
func (u *URL) Get(attribute string) string {
	return u.Data.Get(attribute)
}
