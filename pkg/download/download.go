// Package download provides file‑download utilities with timeout,
// retry, and optional SSL‑verification control.
package download

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/go-getter/v2"
	retryablehttp "github.com/hashicorp/go-retryablehttp"

	"github.com/idelchi/godyl/internal/debug"
	"github.com/idelchi/godyl/pkg/path/file"
	"github.com/idelchi/godyl/pkg/utils"
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

	d := &Downloader{
		contextTimeout: defaultTimeout,
		readTimeout:    defaultTimeout,
		headTimeout:    defaultTimeout,

		maxRetries:   3,
		retryWaitMin: 1 * time.Second,
		retryWaitMax: 30 * time.Second,
	}

	// Apply all provided options
	for _, opt := range opts {
		opt(d)
	}

	return d
}

// Download fetches url to output (archives auto‑extracted).
func (d Downloader) Download(url, output string, header ...http.Header) (file.File, error) {
	ctx, cancel := context.WithTimeout(context.Background(), d.contextTimeout)
	defer cancel()

	// retryable HTTP client
	rc := retryablehttp.NewClient()
	rc.Logger = nil // silence default logging
	rc.RetryMax = d.maxRetries
	rc.RetryWaitMin = d.retryWaitMin
	rc.RetryWaitMax = d.retryWaitMax
	rc.HTTPClient = cleanhttp.DefaultClient()

	httpClient := rc.StandardClient()

	if d.insecureSkipVerify {
		httpClient.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	// merge headers
	h := make(http.Header)
	for _, src := range header {
		for k, vv := range src {
			for _, v := range vv {
				h.Add(k, v)
			}
		}
	}

	httpGetter := &getter.HttpGetter{
		Netrc:                 true,
		XTerraformGetDisabled: true,
		HeadFirstTimeout:      d.headTimeout,
		ReadTimeout:           d.readTimeout,
		Client:                httpClient,
		Header:                h,
	}

	if !utils.IsURL(url) {
		return file.New(), fmt.Errorf("invalid URL: %q", url)
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
		return file.New(), fmt.Errorf("getting file: %w", err)
	}

	debug.Debug("downloaded %q to %q", url, res.Dst)

	return file.New(res.Dst), nil
}
