package github_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	gogithub "github.com/google/go-github/v74/github"

	internalgithub "github.com/idelchi/godyl/internal/github"
	"github.com/idelchi/godyl/internal/release"
)

type releaseJSON struct {
	TagName     string      `json:"tag_name"`
	Name        string      `json:"name"`
	Body        string      `json:"body"`
	PublishedAt *time.Time  `json:"published_at,omitempty"`
	Assets      []assetJSON `json:"assets,omitempty"`
}

type assetJSON struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
	ContentType        string `json:"content_type"`
}

func newTestServer(t *testing.T, owner, repo string, mux *http.ServeMux) *internalgithub.Repository {
	t.Helper()

	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	client := internalgithub.NewClient("", server.URL+"/")

	return internalgithub.NewRepository(owner, repo, client)
}

func TestLatestRelease(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("/repos/myowner/myrepo/releases/latest", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)

			return
		}

		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(releaseJSON{
			TagName: "v2.0.0",
			Name:    "Release v2.0.0",
			Body:    "latest release notes",
			Assets: []assetJSON{
				{
					Name:               "tool-linux-amd64.tar.gz",
					BrowserDownloadURL: "https://example.com/releases/download/v2.0.0/tool-linux-amd64.tar.gz",
					ContentType:        "application/gzip",
				},
			},
		}); err != nil {
			http.Error(w, "encode failed", http.StatusInternalServerError)
		}
	})

	repo := newTestServer(t, "myowner", "myrepo", mux)

	got, err := repo.LatestRelease(t.Context())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := &release.Release{
		Tag:  "v2.0.0",
		Name: "Release v2.0.0",
		Body: "latest release notes",
		Assets: release.Assets{
			{
				Name: "tool-linux-amd64.tar.gz",
				URL:  "https://example.com/releases/download/v2.0.0/tool-linux-amd64.tar.gz",
				Type: "application/gzip",
			},
		},
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("LatestRelease() mismatch (-want +got):\n%s", diff)
	}
}

func TestGetRelease(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("/repos/myowner/myrepo/releases/tags/v1.5.0", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)

			return
		}

		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(releaseJSON{
			TagName: "v1.5.0",
			Name:    "Release v1.5.0",
			Body:    "specific tag release notes",
			Assets: []assetJSON{
				{
					Name:               "tool-darwin-arm64.tar.gz",
					BrowserDownloadURL: "https://example.com/releases/download/v1.5.0/tool-darwin-arm64.tar.gz",
					ContentType:        "application/gzip",
				},
				{
					Name:               "tool-windows-amd64.zip",
					BrowserDownloadURL: "https://example.com/releases/download/v1.5.0/tool-windows-amd64.zip",
					ContentType:        "application/zip",
				},
			},
		}); err != nil {
			http.Error(w, "encode failed", http.StatusInternalServerError)
		}
	})

	repo := newTestServer(t, "myowner", "myrepo", mux)

	got, err := repo.GetRelease(t.Context(), "v1.5.0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := &release.Release{
		Tag:  "v1.5.0",
		Name: "Release v1.5.0",
		Body: "specific tag release notes",
		Assets: release.Assets{
			{
				Name: "tool-darwin-arm64.tar.gz",
				URL:  "https://example.com/releases/download/v1.5.0/tool-darwin-arm64.tar.gz",
				Type: "application/gzip",
			},
			{
				Name: "tool-windows-amd64.zip",
				URL:  "https://example.com/releases/download/v1.5.0/tool-windows-amd64.zip",
				Type: "application/zip",
			},
		},
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("GetRelease() mismatch (-want +got):\n%s", diff)
	}
}

func assertGitHubError(t *testing.T, err error, wantStatus int) {
	t.Helper()

	var ghErr *gogithub.ErrorResponse
	if !errors.As(err, &ghErr) {
		t.Fatalf("expected *github.ErrorResponse, got %T: %v", err, err)
	}

	if ghErr.Response.StatusCode != wantStatus {
		t.Errorf("expected status %d, got %d", wantStatus, ghErr.Response.StatusCode)
	}
}

func TestLatestRelease_NotFound(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("/repos/myowner/myrepo/releases/latest", func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, `{"message": "Not Found"}`, http.StatusNotFound)
	})

	repo := newTestServer(t, "myowner", "myrepo", mux)

	got, err := repo.LatestRelease(t.Context())
	if err == nil {
		t.Fatalf("expected error for 404 response, got release: %+v", got)
	}

	assertGitHubError(t, err, http.StatusNotFound)

	if got != nil {
		t.Errorf("expected nil release on error, got %+v", got)
	}
}

func TestLatestRelease_ServerError(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("/repos/myowner/myrepo/releases/latest", func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, `{"message": "Service Unavailable"}`, http.StatusServiceUnavailable)
	})

	repo := newTestServer(t, "myowner", "myrepo", mux)

	got, err := repo.LatestRelease(t.Context())
	if err == nil {
		t.Fatalf("expected error for 503 response, got release: %+v", got)
	}

	assertGitHubError(t, err, http.StatusServiceUnavailable)

	if got != nil {
		t.Errorf("expected nil release on error, got %+v", got)
	}
}

func TestGetRelease_NotFound(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("/repos/myowner/myrepo/releases/tags/v9.9.9", func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, `{"message": "Not Found"}`, http.StatusNotFound)
	})

	repo := newTestServer(t, "myowner", "myrepo", mux)

	got, err := repo.GetRelease(t.Context(), "v9.9.9")
	if err == nil {
		t.Fatalf("expected error for 404 response, got release: %+v", got)
	}

	assertGitHubError(t, err, http.StatusNotFound)

	if got != nil {
		t.Errorf("expected nil release on error, got %+v", got)
	}
}

func TestLatestIncludingPreRelease_EmptyReleases(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("/repos/myowner/myrepo/releases", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, _ = w.Write([]byte("[]"))
	})

	repo := newTestServer(t, "myowner", "myrepo", mux)

	got, err := repo.LatestIncludingPreRelease(t.Context(), 100)
	if err == nil {
		t.Fatalf("expected error for empty releases, got: %+v", got)
	}

	if got != nil {
		t.Errorf("expected nil release, got %+v", got)
	}
}

func TestLatestIncludingPreRelease(t *testing.T) {
	t.Parallel()

	// Releases are listed out of order: oldest, middle, newest.
	// The function must pick newest by PublishedAt regardless of list position.
	oldest := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	middle := time.Date(2024, 3, 15, 0, 0, 0, 0, time.UTC)
	newest := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)

	releases := []releaseJSON{
		{
			TagName:     "v1.0.0",
			Name:        "Release v1.0.0",
			PublishedAt: &oldest,
			Assets:      []assetJSON{},
		},
		{
			TagName:     "v2.0.0-beta.1",
			Name:        "Pre-release v2.0.0-beta.1",
			PublishedAt: &newest,
			Assets: []assetJSON{
				{
					Name:               "tool-linux-amd64.tar.gz",
					BrowserDownloadURL: "https://example.com/releases/download/v2.0.0-beta.1/tool-linux-amd64.tar.gz",
					ContentType:        "application/gzip",
				},
			},
		},
		{
			TagName:     "v1.5.0",
			Name:        "Release v1.5.0",
			PublishedAt: &middle,
			Assets:      []assetJSON{},
		},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/repos/myowner/myrepo/releases", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)

			return
		}

		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(releases); err != nil {
			http.Error(w, "encode failed", http.StatusInternalServerError)
		}
	})

	repo := newTestServer(t, "myowner", "myrepo", mux)

	got, err := repo.LatestIncludingPreRelease(t.Context(), 100)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := &release.Release{
		Tag:  "v2.0.0-beta.1",
		Name: "Pre-release v2.0.0-beta.1",
		Assets: release.Assets{
			{
				Name: "tool-linux-amd64.tar.gz",
				URL:  "https://example.com/releases/download/v2.0.0-beta.1/tool-linux-amd64.tar.gz",
				Type: "application/gzip",
			},
		},
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("LatestIncludingPreRelease() mismatch (-want +got):\n%s", diff)
	}
}

func TestGetReleasesByWildcard(t *testing.T) {
	t.Parallel()

	// Shared release list served by the test server.
	// v1.0.0, v1.5.0, v1.9.0 are v1.x; v2.0.0 is a different major; v0.9.0 is pre-v1.
	releases := []releaseJSON{
		{
			TagName: "v0.9.0",
			Name:    "Release v0.9.0",
			Assets:  []assetJSON{},
		},
		{
			TagName: "v1.0.0",
			Name:    "Release v1.0.0",
			Assets:  []assetJSON{},
		},
		{
			TagName: "v1.5.0",
			Name:    "Release v1.5.0",
			Assets:  []assetJSON{},
		},
		{
			TagName: "v1.9.0",
			Name:    "Release v1.9.0",
			Assets: []assetJSON{
				{
					Name:               "tool-linux-amd64.tar.gz",
					BrowserDownloadURL: "https://example.com/releases/download/v1.9.0/tool-linux-amd64.tar.gz",
					ContentType:        "application/gzip",
				},
			},
		},
		{
			TagName: "v2.0.0",
			Name:    "Release v2.0.0",
			Assets:  []assetJSON{},
		},
		{
			TagName: "not-semver",
			Name:    "Non-semver tag",
			Assets:  []assetJSON{},
		},
	}

	tests := []struct {
		name        string
		pattern     string
		want        *release.Release
		wantErrFrag string // non-empty means we expect an error containing this substring
	}{
		{
			name:    "wildcard v1.* returns highest v1.x release",
			pattern: "v1.*",
			want: &release.Release{
				Tag:  "v1.9.0",
				Name: "Release v1.9.0",
				Assets: release.Assets{
					{
						Name: "tool-linux-amd64.tar.gz",
						URL:  "https://example.com/releases/download/v1.9.0/tool-linux-amd64.tar.gz",
						Type: "application/gzip",
					},
				},
			},
		},
		{
			name:    "exact version v1.5.0 returns that release",
			pattern: "v1.5.0",
			want: &release.Release{
				Tag:    "v1.5.0",
				Name:   "Release v1.5.0",
				Assets: release.Assets{},
			},
		},
		{
			name:        "no matching releases returns error",
			pattern:     "v3.*",
			wantErrFrag: "no releases match",
		},
		{
			name:        "invalid wildcard pattern returns constraint error",
			pattern:     "not-valid-*",
			wantErrFrag: "invalid version pattern",
		},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/repos/myowner/myrepo/releases", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)

			return
		}

		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(releases); err != nil {
			http.Error(w, "encode failed", http.StatusInternalServerError)
		}
	})

	repo := newTestServer(t, "myowner", "myrepo", mux)

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, err := repo.GetReleasesByWildcard(t.Context(), tc.pattern, 100)

			if tc.wantErrFrag != "" {
				if err == nil {
					t.Fatalf("expected error containing %q, got nil (release: %+v)", tc.wantErrFrag, got)
				}

				if !strings.Contains(err.Error(), tc.wantErrFrag) {
					t.Errorf("error %q does not contain %q", err.Error(), tc.wantErrFrag)
				}

				if got != nil {
					t.Errorf("expected nil release on error, got %+v", got)
				}

				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("GetReleasesByWildcard() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestLatestIncludingPreRelease_MultiPage(t *testing.T) {
	t.Parallel()

	// Page 1 has older releases; page 2 has the newest.
	// The function must find the newest across both pages.
	olderTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	newerTime := time.Date(2024, 12, 1, 0, 0, 0, 0, time.UTC)

	page1 := []releaseJSON{
		{
			TagName:     "v1.0.0",
			Name:        "Release v1.0.0",
			PublishedAt: &olderTime,
			Assets:      []assetJSON{},
		},
	}
	page2 := []releaseJSON{
		{
			TagName:     "v2.0.0",
			Name:        "Release v2.0.0",
			PublishedAt: &newerTime,
			Assets: []assetJSON{
				{
					Name:               "tool-linux-amd64.tar.gz",
					BrowserDownloadURL: "https://example.com/releases/download/v2.0.0/tool-linux-amd64.tar.gz",
					ContentType:        "application/gzip",
				},
			},
		},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/repos/myowner/myrepo/releases", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)

			return
		}

		w.Header().Set("Content-Type", "application/json")

		page := r.URL.Query().Get("page")
		if page == "2" {
			if err := json.NewEncoder(w).Encode(page2); err != nil {
				http.Error(w, "encode failed", http.StatusInternalServerError)
			}

			return
		}

		// Page 1: include a Link header pointing at page 2. The go-github
		// client parses the Link header to determine resp.NextPage.
		// The URL host/scheme don't matter — only the ?page= value is used.
		w.Header().Set("Link", `</repos/myowner/myrepo/releases?page=2>; rel="next"`)

		if err := json.NewEncoder(w).Encode(page1); err != nil {
			http.Error(w, "encode failed", http.StatusInternalServerError)
		}
	})

	repo := newTestServer(t, "myowner", "myrepo", mux)

	got, err := repo.LatestIncludingPreRelease(t.Context(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := &release.Release{
		Tag:  "v2.0.0",
		Name: "Release v2.0.0",
		Assets: release.Assets{
			{
				Name: "tool-linux-amd64.tar.gz",
				URL:  "https://example.com/releases/download/v2.0.0/tool-linux-amd64.tar.gz",
				Type: "application/gzip",
			},
		},
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("LatestIncludingPreRelease_MultiPage() mismatch (-want +got):\n%s", diff)
	}
}

func TestLatestRelease_EmptyTagName(t *testing.T) {
	t.Parallel()

	// The GitHub API returns a release whose tag_name is an empty string (non-nil pointer).
	// FromRepositoryRelease only nil-checks the pointer, so a release with Tag="" is returned.
	mux := http.NewServeMux()
	mux.HandleFunc("/repos/myowner/myrepo/releases/latest", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)

			return
		}

		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(releaseJSON{
			TagName: "",
			Name:    "Unnamed release",
			Body:    "body text",
			Assets:  []assetJSON{},
		}); err != nil {
			http.Error(w, "encode failed", http.StatusInternalServerError)
		}
	})

	repo := newTestServer(t, "myowner", "myrepo", mux)

	got, err := repo.LatestRelease(t.Context())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// The empty tag_name propagates through — Tag is "" in the result.
	want := &release.Release{
		Tag:    "",
		Name:   "Unnamed release",
		Body:   "body text",
		Assets: release.Assets{},
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("LatestRelease_EmptyTagName() mismatch (-want +got):\n%s", diff)
	}
}
