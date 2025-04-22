package updater

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-resty/resty/v2"
)

type Latest struct {
	Version   string
	Changelog string
}

// Get reaches out to https://idelchi.github.io/godyl to check if there's a new version available.
func (l *Latest) Get(pre bool) error {
	client := resty.New()

	url := "https://idelchi.github.io/godyl/_versions/latest"
	if pre {
		url = "https://idelchi.github.io/godyl/_versions/pre"
	}

	resp, err := client.R().Get(url)
	if err != nil {
		return fmt.Errorf("checking for latest version: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}

	// Unmarshal the response body into the versions slice
	err = json.Unmarshal(resp.Body(), l)
	if err != nil {
		return fmt.Errorf("unmarshaling versions: %w", err)
	}

	// Trim the changelog to remove any leading or trailing whitespace
	l.Changelog = strings.TrimSpace(l.Changelog)

	return nil
}
