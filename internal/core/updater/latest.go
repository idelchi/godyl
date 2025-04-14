package updater

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-resty/resty/v2"
)

// Latest reaches out to https://idelchi.github.io/godyl/latest to check if there's a new version available.
func Latest() (string, error) {
	client := resty.New()
	resp, err := client.R().Get("https://idelchi.github.io/godyl/latest")
	if err != nil {
		return "", fmt.Errorf("checking for latest version: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}

	version := strings.TrimSpace(string(resp.Body()))
	return version, nil
}

// AllVersions reaches out to https://idelchi.github.io/godyl/versions to retrieve all versions.
func AllVersions() ([]string, error) {
	var versions []string

	client := resty.New()
	resp, err := client.R().Get("https://idelchi.github.io/godyl/versions") // Note: URL changed to /versions as per function comment
	if err != nil {
		return nil, fmt.Errorf("checking for versions: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}

	// Unmarshal the response body into the versions slice
	err = json.Unmarshal(resp.Body(), &versions)
	if err != nil {
		return nil, fmt.Errorf("unmarshaling versions: %w", err)
	}

	return versions, nil
}
