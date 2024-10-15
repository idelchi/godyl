// Package download provides functionality for downloading files from URLs.
package download

import (
	"context"
	"time"

	"github.com/hashicorp/go-getter/v2"
	"github.com/idelchi/godyl/pkg/file"
)

type Downloader struct {
	ContextTimeout time.Duration
	ReadTimeout    time.Duration
	HeadTimeout    time.Duration
}

func New() *Downloader {
	return &Downloader{
		ContextTimeout: 5 * time.Minute,
		ReadTimeout:    5 * time.Minute,
		HeadTimeout:    5 * time.Minute,
	}
}

// Download fetches a file from the given URL and saves it to the specified output path.
// If the file is an archive, it will be extracted to the output directory.
// It returns the destination path of the downloaded file (or folder) and any error encountered.
func (d Downloader) Download(url string, output string) (file.File, error) {
	ctx, cancel := context.WithTimeout(context.Background(), d.ContextTimeout)
	defer cancel()

	req := &getter.Request{
		Src:     url,
		Dst:     output,
		GetMode: getter.ModeAny,
	}
	client := &getter.Client{
		Getters: []getter.Getter{
			&getter.HttpGetter{
				Netrc:                 true,
				XTerraformGetDisabled: true,
				HeadFirstTimeout:      d.HeadTimeout,
				ReadTimeout:           d.ReadTimeout,
			},
		},
	}

	res, err := client.Get(ctx, req)
	if err != nil {
		return file.New(), err
	}

	return file.New(res.Dst), nil
}
