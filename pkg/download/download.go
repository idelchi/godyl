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

// Downloader manages the configuration for downloading files, including
// timeouts for different stages of the process.
type Downloader struct {
	// ContextTimeout is the maximum duration to wait for the download context.
	ContextTimeout time.Duration
	// ReadTimeout is the maximum duration to wait for reading data from the URL.
	ReadTimeout time.Duration
	// HeadTimeout is the maximum duration to wait for the HTTP HEAD request.
	HeadTimeout time.Duration
	// InsecureSkipVerify controls whether to verify SSL certificates.
	// WARNING: Setting this to true is insecure and should only be used in testing.
	InsecureSkipVerify bool
}

// New returns a new Downloader instance with default timeout values set to 5 minutes.
func New() *Downloader {
	const defaultTimeout = 5 * time.Minute

	return &Downloader{
		ContextTimeout:     defaultTimeout,
		ReadTimeout:        defaultTimeout,
		HeadTimeout:        defaultTimeout,
		InsecureSkipVerify: false,
	}
}

// Download fetches a file from the given URL and saves it to the specified output path.
// If the file is an archive, it will be extracted to the output directory.
// It returns the destination path of the downloaded file (or folder) and any error encountered.
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
