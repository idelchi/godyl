// Package download provides functionality for downloading files from URLs,
// with support for various protocols and automatic extraction of archives.
// The Downloader struct allows configuration of timeout settings for the download
// context, read operations, and HTTP HEAD requests.
//
// This package is built on top of HashiCorp's go-getter library, which supports
// downloading files from a variety of protocols (HTTP, HTTPS, FTP, etc.), and
// includes automatic handling of archives such as zip or tar files.
//
// Example usage:
//
//	package main
//
//	import (
//	    "log"
//	    "github.com/idelchi/godyl/pkg/download"
//	)
//
//	func main() {
//	    d := download.New()
//	    file, err := d.Download("https://example.com/file.zip", "/path/to/output")
//	    if err != nil {
//	        log.Fatal(err)
//	    }
//	    log.Println("Downloaded to:", file)
//	}
package download

import (
	"context"
	"crypto/tls"
	"net/http"
	"time"

	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/go-getter/v2"

	"github.com/idelchi/godyl/pkg/path/file"
)

// Downloader manages file download operations and configuration.
// Provides configurable timeouts, SSL verification control, and
// progress tracking for downloading files from various sources.
type Downloader struct {
	// ContextTimeout limits the total download operation time.
	ContextTimeout time.Duration

	// ReadTimeout limits time for reading data chunks.
	ReadTimeout time.Duration

	// HeadTimeout limits time for initial HTTP HEAD request.
	HeadTimeout time.Duration

	// InsecureSkipVerify disables SSL certificate validation.
	// WARNING: Setting this to true is insecure and should only
	// be used in testing environments.
	InsecureSkipVerify bool

	// ProgressListener receives download progress updates.
	ProgressListener getter.ProgressTracker
}

// New creates a Downloader with default settings.
// Returns a Downloader configured with 5-minute timeouts for
// context, read operations, and HEAD requests. SSL verification
// is enabled by default.
func New() *Downloader {
	const defaultTimeout = 5 * time.Minute

	return &Downloader{
		ContextTimeout:     defaultTimeout,
		ReadTimeout:        defaultTimeout,
		HeadTimeout:        defaultTimeout,
		InsecureSkipVerify: false,
	}
}

// Download retrieves and processes a file from a URL.
// Downloads from the specified URL to the output path, handling
// various protocols (HTTP, HTTPS, FTP). Automatically extracts
// archives (zip, tar) to the output directory. Returns the final
// destination path and any errors encountered.
func (d Downloader) Download(url, output string, header ...http.Header) (file.File, error) {
	ctx, cancel := context.WithTimeout(context.Background(), d.ContextTimeout)
	defer cancel()

	httpClient := cleanhttp.DefaultClient()

	headers := make(http.Header)

	// Merge all headers from the variadic argument
	for _, h := range header {
		for key, values := range h {
			for _, value := range values {
				headers.Add(key, value)
			}
		}
	}

	httpGetter := &getter.HttpGetter{
		Netrc:                 true,
		XTerraformGetDisabled: true,
		HeadFirstTimeout:      d.HeadTimeout,
		ReadTimeout:           d.ReadTimeout,
		Client:                httpClient,
		Header:                headers,
	}

	// Modify the default HTTP client's transport to skip SSL verification if requested
	if d.InsecureSkipVerify {
		httpClient.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
	}

	req := &getter.Request{
		Src:     url,
		Dst:     output,
		GetMode: getter.ModeAny,
	}

	// Pass the progress listener if provided
	if d.ProgressListener != nil {
		req.ProgressListener = d.ProgressListener
	}

	// TODO(Idelchi): Go-Getter messes up ? queries etc and doesn't seem to follow redirects then,
	// or perhaps messes up the whole URL
	client := &getter.Client{
		Getters: []getter.Getter{
			httpGetter,
		},
	}

	res, err := client.Get(ctx, req)
	if err != nil {
		return file.New(), err
	}

	return file.New(res.Dst), nil
}
