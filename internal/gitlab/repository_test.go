package gitlab_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	internalgitlab "github.com/idelchi/godyl/internal/gitlab"
	"github.com/idelchi/godyl/internal/release"
)

type gitlabReleaseJSON struct {
	TagName     string           `json:"tag_name"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	CreatedAt   *time.Time       `json:"created_at,omitempty"`
	Assets      gitlabAssetsJSON `json:"assets"`
}

type gitlabAssetsJSON struct {
	Links []gitlabLinkJSON `json:"links"`
}

type gitlabLinkJSON struct {
	Name           string `json:"name"`
	URL            string `json:"url"`
	DirectAssetURL string `json:"direct_asset_url"`
	LinkType       string `json:"link_type"`
}

func newGitLabTestServer(t *testing.T, mux *http.ServeMux) *internalgitlab.Repository {
	t.Helper()

	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	client, err := internalgitlab.NewClient("", server.URL)
	if err != nil {
		t.Fatalf("creating client: %v", err)
	}

	return internalgitlab.NewRepository("mygroup", "myrepo", client)
}

func TestGetRelease(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v4/projects/mygroup%2Fmyrepo/releases/v1.2.3", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)

			return
		}

		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(gitlabReleaseJSON{
			TagName:     "v1.2.3",
			Name:        "Release v1.2.3",
			Description: "release description",
			Assets: gitlabAssetsJSON{
				Links: []gitlabLinkJSON{
					{
						Name:           "tool-linux-amd64.tar.gz",
						URL:            "https://example.com/files/tool-linux-amd64.tar.gz",
						DirectAssetURL: "https://example.com/direct/tool-linux-amd64.tar.gz",
						LinkType:       "package",
					},
					{
						Name:           "tool-darwin-arm64.tar.gz",
						URL:            "https://example.com/files/tool-darwin-arm64.tar.gz",
						DirectAssetURL: "https://example.com/direct/tool-darwin-arm64.tar.gz",
						LinkType:       "package",
					},
				},
			},
		}); err != nil {
			http.Error(w, "encode failed", http.StatusInternalServerError)
		}
	})

	repo := newGitLabTestServer(t, mux)

	got, err := repo.GetRelease(t.Context(), "v1.2.3")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := &release.Release{
		Tag:  "v1.2.3",
		Name: "Release v1.2.3",
		Assets: release.Assets{
			{
				Name: "tool-linux-amd64.tar.gz",
				URL:  "https://example.com/direct/tool-linux-amd64.tar.gz",
				Type: "package",
			},
			{
				Name: "tool-darwin-arm64.tar.gz",
				URL:  "https://example.com/direct/tool-darwin-arm64.tar.gz",
				Type: "package",
			},
		},
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("GetRelease() mismatch (-want +got):\n%s", diff)
	}
}

func TestGetRelease_NotFound(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v4/projects/mygroup%2Fmyrepo/releases/v9.9.9", func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, `{"message": "404 Not Found"}`, http.StatusNotFound)
	})

	repo := newGitLabTestServer(t, mux)

	got, err := repo.GetRelease(t.Context(), "v9.9.9")
	if err == nil {
		t.Fatalf("expected error for 404 response, got release: %+v", got)
	}

	if got != nil {
		t.Errorf("expected nil release on error, got %+v", got)
	}

	errStr := err.Error()
	if !strings.Contains(errStr, "404") {
		t.Errorf("expected error string to contain %q, got %q", "404", errStr)
	}
}

func TestLatestRelease(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v4/projects/mygroup%2Fmyrepo/releases", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)

			return
		}

		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode([]gitlabReleaseJSON{
			{
				TagName: "v3.0.0",
				Name:    "Release v3.0.0",
				Assets: gitlabAssetsJSON{
					Links: []gitlabLinkJSON{
						{
							Name:           "tool-linux-amd64.tar.gz",
							URL:            "https://example.com/files/tool-linux-amd64.tar.gz",
							DirectAssetURL: "https://example.com/direct/tool-linux-amd64.tar.gz",
							LinkType:       "package",
						},
					},
				},
			},
			{
				TagName: "v2.0.0",
				Name:    "Release v2.0.0",
				Assets:  gitlabAssetsJSON{},
			},
		}); err != nil {
			http.Error(w, "encode failed", http.StatusInternalServerError)
		}
	})

	repo := newGitLabTestServer(t, mux)

	got, err := repo.LatestRelease(t.Context())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := &release.Release{
		Tag:  "v3.0.0",
		Name: "Release v3.0.0",
		Assets: release.Assets{
			{
				Name: "tool-linux-amd64.tar.gz",
				URL:  "https://example.com/direct/tool-linux-amd64.tar.gz",
				Type: "package",
			},
		},
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("LatestRelease() mismatch (-want +got):\n%s", diff)
	}
}

func TestLatestRelease_NoReleases(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v4/projects/mygroup%2Fmyrepo/releases", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, _ = w.Write([]byte("[]"))
	})

	repo := newGitLabTestServer(t, mux)

	got, err := repo.LatestRelease(t.Context())
	if err == nil {
		t.Fatalf("expected error when no releases exist, got: %+v", got)
	}

	if got != nil {
		t.Errorf("expected nil release, got %+v", got)
	}
}

func TestLatestRelease_ServerError(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v4/projects/mygroup%2Fmyrepo/releases", func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, `{"message": "Service Unavailable"}`, http.StatusServiceUnavailable)
	})

	repo := newGitLabTestServer(t, mux)

	got, err := repo.LatestRelease(t.Context())
	if err == nil {
		t.Fatalf("expected error for 503 response, got release: %+v", got)
	}

	if got != nil {
		t.Errorf("expected nil release on error, got %+v", got)
	}
}

func TestGetLatestIncludingPreRelease_EmptyReleases(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v4/projects/mygroup%2Fmyrepo/releases", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, _ = w.Write([]byte("[]"))
	})

	repo := newGitLabTestServer(t, mux)

	got, err := repo.GetLatestIncludingPreRelease(t.Context(), 100)
	if err == nil {
		t.Fatalf("expected error for empty releases, got: %+v", got)
	}

	if got != nil {
		t.Errorf("expected nil release, got %+v", got)
	}
}

func TestGetLatestIncludingPreRelease(t *testing.T) {
	t.Parallel()

	// Releases are listed out of order: oldest, newest, middle.
	// The function must pick newest by CreatedAt regardless of list position.
	oldest := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	middle := time.Date(2024, 3, 15, 0, 0, 0, 0, time.UTC)
	newest := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v4/projects/mygroup%2Fmyrepo/releases", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)

			return
		}

		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode([]gitlabReleaseJSON{
			{
				TagName:   "v1.0.0",
				Name:      "Release v1.0.0",
				CreatedAt: &oldest,
				Assets:    gitlabAssetsJSON{},
			},
			{
				TagName:   "v2.0.0-beta.1",
				Name:      "Pre-release v2.0.0-beta.1",
				CreatedAt: &newest,
				Assets: gitlabAssetsJSON{
					Links: []gitlabLinkJSON{
						{
							Name:           "tool-linux-amd64.tar.gz",
							URL:            "https://example.com/files/tool-linux-amd64.tar.gz",
							DirectAssetURL: "https://example.com/direct/tool-linux-amd64.tar.gz",
							LinkType:       "package",
						},
					},
				},
			},
			{
				TagName:   "v1.5.0",
				Name:      "Release v1.5.0",
				CreatedAt: &middle,
				Assets:    gitlabAssetsJSON{},
			},
		}); err != nil {
			http.Error(w, "encode failed", http.StatusInternalServerError)
		}
	})

	repo := newGitLabTestServer(t, mux)

	got, err := repo.GetLatestIncludingPreRelease(t.Context(), 100)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := &release.Release{
		Tag:  "v2.0.0-beta.1",
		Name: "Pre-release v2.0.0-beta.1",
		Assets: release.Assets{
			{
				Name: "tool-linux-amd64.tar.gz",
				URL:  "https://example.com/direct/tool-linux-amd64.tar.gz",
				Type: "package",
			},
		},
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("GetLatestIncludingPreRelease() mismatch (-want +got):\n%s", diff)
	}
}

func TestGetLatestIncludingPreRelease_NilCreatedAt(t *testing.T) {
	t.Parallel()

	// The GitLab GetLatestIncludingPreRelease implementation does:
	//   latestRelease.CreatedAt.Before(*release.CreatedAt)
	// When any release beyond the first has a nil CreatedAt, dereferencing it
	// causes a nil-pointer panic. This test documents that current behavior.
	//
	// A single-release list is used so that the loop body that dereferences
	// CreatedAt is never reached (the i==0 guard short-circuits). Then a
	// two-release list — where the second release has nil CreatedAt — is used
	// to trigger the actual panic path.

	// Case 1: only one release, nil CreatedAt — no panic (i==0 branch taken).
	t.Run("single release with nil CreatedAt returns release", func(t *testing.T) {
		t.Parallel()

		mux := http.NewServeMux()
		mux.HandleFunc("/api/v4/projects/mygroup%2Fmyrepo/releases", func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")

			// Omit created_at so it deserialises as nil.
			if err := json.NewEncoder(w).Encode([]gitlabReleaseJSON{
				{
					TagName: "v1.0.0",
					Name:    "Release v1.0.0",
					Assets:  gitlabAssetsJSON{},
					// CreatedAt deliberately absent (nil after JSON decode)
				},
			}); err != nil {
				http.Error(w, "encode failed", http.StatusInternalServerError)
			}
		})

		repo := newGitLabTestServer(t, mux)

		got, err := repo.GetLatestIncludingPreRelease(t.Context(), 100)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		want := &release.Release{
			Tag:    "v1.0.0",
			Name:   "Release v1.0.0",
			Assets: release.Assets{},
		}

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("GetLatestIncludingPreRelease() mismatch (-want +got):\n%s", diff)
		}
	})

	// Case 2: two releases where the second has nil CreatedAt — the comparison
	//   latestRelease.CreatedAt.Before(*release.CreatedAt)
	// dereferences a nil pointer and panics. The test recovers the panic and
	// documents it as the current (unfixed) behaviour.
	t.Run("second release with nil CreatedAt panics", func(t *testing.T) {
		t.Parallel()

		firstTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

		mux := http.NewServeMux()
		mux.HandleFunc("/api/v4/projects/mygroup%2Fmyrepo/releases", func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")

			if err := json.NewEncoder(w).Encode([]gitlabReleaseJSON{
				{
					TagName:   "v1.0.0",
					Name:      "Release v1.0.0",
					CreatedAt: &firstTime,
					Assets:    gitlabAssetsJSON{},
				},
				{
					TagName: "v2.0.0",
					Name:    "Release v2.0.0",
					Assets:  gitlabAssetsJSON{},
					// CreatedAt deliberately absent (nil after JSON decode)
				},
			}); err != nil {
				http.Error(w, "encode failed", http.StatusInternalServerError)
			}
		})

		repo := newGitLabTestServer(t, mux)

		// Catch the expected nil-pointer panic from the unguarded CreatedAt dereference.
		defer func() {
			r := recover()
			if r == nil {
				t.Error("expected a panic from nil CreatedAt dereference, but no panic occurred")
			}
		}()

		_, _ = repo.GetLatestIncludingPreRelease(t.Context(), 100)
	})
}
