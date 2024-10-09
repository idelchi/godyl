package download

import (
	"context"
	"time"

	"github.com/hashicorp/go-getter/v2"
)

func Download(ctx context.Context, url string, output string) (string, error) {
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
				HeadFirstTimeout:      1 * time.Minute,
				ReadTimeout:           2 * time.Minute,
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
