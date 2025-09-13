// Package download provides file‑download utilities with timeout,
// retry, and optional SSL‑verification control.
package download

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/go-getter/v2"
	retryablehttp "github.com/hashicorp/go-retryablehttp"

	"github.com/idelchi/godyl/internal/debug"
	"github.com/idelchi/godyl/internal/tools/checksum"
	"github.com/idelchi/godyl/pkg/generic"
	"github.com/idelchi/godyl/pkg/path/file"
)

// Downloader manages download settings.
type Downloader struct {
	progressListener   getter.ProgressTracker
	contextTimeout     time.Duration
	readTimeout        time.Duration
	headTimeout        time.Duration
	insecureSkipVerify bool
	checksum           checksum.Checksum

	// retry settings
	maxRetries   int
	retryWaitMin time.Duration
	retryWaitMax time.Duration
}

// New creates a new Downloader with default settings and applies the given options.
func New(opts ...Option) *Downloader {
	const defaultTimeout = 5 * time.Minute

	const defaultRetries = 3

	const defaultRetryWaitMin = time.Second

	const defaultRetryWaitMax = 30 * time.Second

	downloader := &Downloader{
		contextTimeout: defaultTimeout,
		readTimeout:    defaultTimeout,
		headTimeout:    defaultTimeout,

		maxRetries:   defaultRetries,
		retryWaitMin: defaultRetryWaitMin,
		retryWaitMax: defaultRetryWaitMax,
	}

	// Apply all provided options
	for _, opt := range opts {
		opt(downloader)
	}

	return downloader
}

// ErrDownload indicates a download operation failed.
var ErrDownload = errors.New("download error")

// URLWithChecksum  appends the checksum query parameter to the URL if a checksum is provided.
func URLWithChecksum(url string, c checksum.Checksum) string {
	// Append checksum if configured
	if c.Type != "" {
		// Direct hash value
		param := "checksum=" + c.Type + ":" + c.Value

		if strings.Contains(url, "?") {
			url += "&" + param
		} else {
			url += "?" + param
		}
	}

	return url
}

// Download fetches url to output (archives auto‑extracted).
//
//nolint:gocognit // Complex error handling
func (d Downloader) Download(url, output string, header ...http.Header) (file.File, error) {
	ctx, cancel := context.WithTimeout(context.Background(), d.contextTimeout)
	defer cancel()

	// retryable HTTP client
	client := retryablehttp.NewClient()

	client.Logger = nil // silence default logging
	client.RetryMax = d.maxRetries
	client.RetryWaitMin = d.retryWaitMin
	client.RetryWaitMax = d.retryWaitMax
	client.HTTPClient = cleanhttp.DefaultClient()

	httpClient := client.StandardClient()

	if d.insecureSkipVerify {
		httpClient.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, //nolint:gosec   // Only set if explicitly requested
		}
	}

	// merge headers
	headers := make(http.Header)

	for _, h := range header {
		for k, vv := range h {
			for _, v := range vv {
				headers.Add(k, v)
			}
		}
	}

	httpGetter := &getter.HttpGetter{
		Netrc:                 true,
		XTerraformGetDisabled: true,
		HeadFirstTimeout:      d.headTimeout,
		ReadTimeout:           d.readTimeout,
		Client:                httpClient,
		Header:                headers,
	}

	if !generic.IsURL(url) {
		return file.New(), fmt.Errorf("%w: invalid URL: %q", ErrDownload, url)
	}

	rawURL := url

	url = URLWithChecksum(url, d.checksum)

	req := &getter.Request{
		Src:              url,
		Dst:              output,
		GetMode:          getter.ModeAny,
		ProgressListener: d.progressListener,
	}

	debug.Debug("downloading %q to %q", url, output)

	res, err := (&getter.Client{Getters: []getter.Getter{httpGetter}}).Get(ctx, req)
	if err != nil { //nolint:nestif // Complex error handling
		debug.Debug("tried to download %q to %q", url, output)
		debug.Debug("error: %v", err)

		var checksumErr *getter.ChecksumError

		if errors.As(err, &checksumErr) || d.checksum.Optional {
			if errors.As(err, &checksumErr) {
				debug.Debug("checksum mismatch for %s (got=%x expected=%x)",
					checksumErr.File, checksumErr.Actual, checksumErr.Expected)
			}

			if d.checksum.Optional {
				if errors.As(err, &checksumErr) {
					debug.Debug("continuing despite checksum mismatch as it is marked optional")
				} else {
					debug.Debug("retrying in case the failure was due to some checksum fetch issue")
				}

				url = rawURL

				req := &getter.Request{
					Src:              url,
					Dst:              output,
					GetMode:          getter.ModeAny,
					ProgressListener: d.progressListener,
				}

				debug.Debug("downloading %q to %q", url, output)

				res, err = (&getter.Client{Getters: []getter.Getter{httpGetter}}).Get(ctx, req)
				if err != nil {
					debug.Debug("tried to download without checksum %q to %q", req.Src, output)
					debug.Debug("error: %v", err)
				}
			}
		}
	}

	if err != nil {
		return file.New(), fmt.Errorf("%w: getting file: %w", ErrDownload, err)
	}

	debug.Debug("downloaded %q to %q", rawURL, res.Dst)

	return file.New(res.Dst), nil
}
