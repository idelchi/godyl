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

	"github.com/idelchi/godyl/pkg/path/file"
)

// Downloader manages download settings.
type Downloader struct {
	ProgressListener   getter.ProgressTracker
	ContextTimeout     time.Duration
	ReadTimeout        time.Duration
	HeadTimeout        time.Duration
	InsecureSkipVerify bool

	// retry settings
	MaxRetries   int
	RetryWaitMin time.Duration
	RetryWaitMax time.Duration
}

// New returns a Downloader with sane defaults.
func New() *Downloader {
	const defaultTimeout = 5 * time.Minute

	return &Downloader{
		ContextTimeout:     defaultTimeout,
		ReadTimeout:        defaultTimeout,
		HeadTimeout:        defaultTimeout,
		InsecureSkipVerify: false,

		MaxRetries:   3,
		RetryWaitMin: 1 * time.Second,
		RetryWaitMax: 30 * time.Second,
	}
}

// Download fetches url to output (archives auto‑extracted).
func (d Downloader) Download(url, output string, header ...http.Header) (file.File, error) {
	ctx, cancel := context.WithTimeout(context.Background(), d.ContextTimeout)
	defer cancel()

	// retryable HTTP client
	rc := retryablehttp.NewClient()
	rc.Logger = nil // silence default logging
	rc.RetryMax = d.MaxRetries
	rc.RetryWaitMin = d.RetryWaitMin
	rc.RetryWaitMax = d.RetryWaitMax
	rc.HTTPClient = cleanhttp.DefaultClient()

	httpClient := rc.StandardClient()

	if d.InsecureSkipVerify {
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
		HeadFirstTimeout:      d.HeadTimeout,
		ReadTimeout:           d.ReadTimeout,
		Client:                httpClient,
		Header:                h,
	}

	req := &getter.Request{
		Src:     url,
		Dst:     output,
		GetMode: getter.ModeAny,
	}

	if d.ProgressListener != nil {
		req.ProgressListener = d.ProgressListener
	}

	res, err := (&getter.Client{Getters: []getter.Getter{httpGetter}}).Get(ctx, req)
	if err != nil {
		return file.New(), fmt.Errorf("getting file: %w", err)
	}

	return file.New(res.Dst), nil
}
