// Package download provides file‑download utilities with timeout,
// retry, and optional SSL‑verification control.
package download

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/go-getter/v2"
	retryablehttp "github.com/hashicorp/go-retryablehttp"

	"github.com/idelchi/godyl/internal/debug"
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

// Download fetches url to output (archives auto‑extracted).
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

	req := &getter.Request{
		Src:              url,
		Dst:              output,
		GetMode:          getter.ModeAny,
		ProgressListener: d.progressListener,
	}

	res, err := (&getter.Client{Getters: []getter.Getter{httpGetter}}).Get(ctx, req)
	if err != nil {
		debug.Debug("tried to download %q to %q", url, output)
		debug.Debug("error: %v", err)

		return file.New(), fmt.Errorf("%w: getting file: %w", ErrDownload, err)
	}

	debug.Debug("downloaded %q to %q", url, res.Dst)

	return file.New(res.Dst), nil
}
