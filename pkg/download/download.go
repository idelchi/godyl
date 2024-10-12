// Package download provides functionality for downloading files from URLs.
package download

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/hashicorp/go-getter/v2"
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
func (d Downloader) Download(url string, output string) (string, error) {
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

	res, err := client.Get(
		ctx, req,
	)
	if err != nil {
		return "", err
	}

	return res.Dst, err
}

// Type represents the type of a file.
type Type int

const (
	Unknown Type = iota
	File
	Directory
)

type Result struct {
	Path string
	Type Type
}

func NewResult(path string) (Result, error) {
	r := Result{
		Path: path,
	}

	// Check if the path is a directory
	info, err := os.Stat(r.Path)
	if err != nil {
		return r, fmt.Errorf("file info for path %q: %v", path, err)
	}

	switch mode := info.Mode(); {
	case mode.IsDir():
		r.Type = Directory
		return r, nil
	case mode.IsRegular():
		r.Type = File
		return r, nil
	default:
		return r, fmt.Errorf("unknown type %s", mode)
	}
}
