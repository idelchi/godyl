package github_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/google/go-cmp/cmp"

	internalgithub "github.com/idelchi/godyl/internal/github"
	"github.com/idelchi/godyl/internal/release"
)

func TestParseGitHubReleaseAssets(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		html       string
		wantAssets []release.Asset
	}{
		{
			name: "single asset",
			html: heredoc.Doc(`
				<ul>
				  <li class="Box-row">
				    <a href="/owner/repo/releases/download/v1.0.0/tool-linux-amd64.tar.gz">tool-linux-amd64.tar.gz</a>
				  </li>
				</ul>
			`),
			wantAssets: []release.Asset{
				{
					Name: "tool-linux-amd64.tar.gz",
					URL:  "https://github.com/owner/repo/releases/download/v1.0.0/tool-linux-amd64.tar.gz",
					Type: "text/plain",
				},
			},
		},
		{
			name: "two assets",
			html: heredoc.Doc(`
				<ul>
				  <li class="Box-row">
				    <a href="/owner/repo/releases/download/v2.0.0/tool-linux-amd64.tar.gz">tool-linux-amd64.tar.gz</a>
				  </li>
				  <li class="Box-row">
				    <a href="/owner/repo/releases/download/v2.0.0/tool-darwin-arm64.tar.gz">tool-darwin-arm64.tar.gz</a>
				  </li>
				</ul>
			`),
			wantAssets: []release.Asset{
				{
					Name: "tool-linux-amd64.tar.gz",
					URL:  "https://github.com/owner/repo/releases/download/v2.0.0/tool-linux-amd64.tar.gz",
					Type: "text/plain",
				},
				{
					Name: "tool-darwin-arm64.tar.gz",
					URL:  "https://github.com/owner/repo/releases/download/v2.0.0/tool-darwin-arm64.tar.gz",
					Type: "text/plain",
				},
			},
		},
		{
			name: "asset with valid digest",
			html: heredoc.Doc(`
				<ul>
				  <li class="Box-row">
				    <a href="/owner/repo/releases/download/v1.0.0/tool-linux-amd64.tar.gz">tool-linux-amd64.tar.gz</a>
				    <span class="Truncate-text">sha256:2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824</span>
				  </li>
				</ul>
			`),
			wantAssets: []release.Asset{
				{
					Name:   "tool-linux-amd64.tar.gz",
					URL:    "https://github.com/owner/repo/releases/download/v1.0.0/tool-linux-amd64.tar.gz",
					Type:   "text/plain",
					Digest: "sha256:2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824",
				},
			},
		},
		{
			name: "digest with invalid hex is ignored",
			html: heredoc.Doc(`
				<ul>
				  <li class="Box-row">
				    <a href="/owner/repo/releases/download/v1.0.0/tool-linux-amd64.tar.gz">tool-linux-amd64.tar.gz</a>
				    <span class="Truncate-text">sha256:not-valid-hex!!!</span>
				  </li>
				</ul>
			`),
			wantAssets: []release.Asset{
				{
					Name: "tool-linux-amd64.tar.gz",
					URL:  "https://github.com/owner/repo/releases/download/v1.0.0/tool-linux-amd64.tar.gz",
					Type: "text/plain",
				},
			},
		},
		{
			name:       "empty HTML returns no assets",
			html:       "",
			wantAssets: nil,
		},
		{
			name: "no Box-row elements returns no assets",
			html: `<ul><li class="some-other-class">` +
				`<a href="/owner/repo/releases/download/v1.0.0/tool.tar.gz">tool.tar.gz</a></li></ul>`,
			wantAssets: nil,
		},
		{
			name:       "Box-row without download link returns no assets",
			html:       `<ul><li class="Box-row"><a href="/owner/repo/tree/main">source code</a></li></ul>`,
			wantAssets: nil,
		},
		{
			name: "mixed Box-row elements, only download links counted",
			html: heredoc.Doc(`
				<ul>
				  <li class="Box-row">
				    <a href="/owner/repo/releases/download/v3.0.0/app-linux.tar.gz">app-linux.tar.gz</a>
				  </li>
				  <li class="Box-row">
				    <a href="/owner/repo/commit/abc123">not a download</a>
				  </li>
				  <li class="Box-row">
				    <a href="/owner/repo/releases/download/v3.0.0/app-windows.zip">app-windows.zip</a>
				  </li>
				</ul>
			`),
			wantAssets: []release.Asset{
				{
					Name: "app-linux.tar.gz",
					URL:  "https://github.com/owner/repo/releases/download/v3.0.0/app-linux.tar.gz",
					Type: "text/plain",
				},
				{
					Name: "app-windows.zip",
					URL:  "https://github.com/owner/repo/releases/download/v3.0.0/app-windows.zip",
					Type: "text/plain",
				},
			},
		},
		{
			name: "asset with digest and one without",
			html: heredoc.Doc(`
				<ul>
				  <li class="Box-row">
				    <a href="/owner/repo/releases/download/v1.0.0/app.tar.gz">app.tar.gz</a>
				    <span class="Truncate-text">sha256:deadbeef1234abcd</span>
				  </li>
				  <li class="Box-row">
				    <a href="/owner/repo/releases/download/v1.0.0/app.tar.gz.sha256">app.tar.gz.sha256</a>
				  </li>
				</ul>
			`),
			wantAssets: []release.Asset{
				{
					Name:   "app.tar.gz",
					URL:    "https://github.com/owner/repo/releases/download/v1.0.0/app.tar.gz",
					Type:   "text/plain",
					Digest: "sha256:deadbeef1234abcd",
				},
				{
					Name: "app.tar.gz.sha256",
					URL:  "https://github.com/owner/repo/releases/download/v1.0.0/app.tar.gz.sha256",
					Type: "text/plain",
				},
			},
		},
		{
			name:       "HTML with no list at all",
			html:       `<html><body><p>No releases here</p></body></html>`,
			wantAssets: nil,
		},
		{
			name: "URL prefixed with github.com",
			html: heredoc.Doc(`
				<ul>
				  <li class="Box-row">
				    <a href="/myorg/mytool/releases/download/v0.9.1/mytool-freebsd-amd64">mytool-freebsd-amd64</a>
				  </li>
				</ul>
			`),
			wantAssets: []release.Asset{
				{
					Name: "mytool-freebsd-amd64",
					URL:  "https://github.com/myorg/mytool/releases/download/v0.9.1/mytool-freebsd-amd64",
					Type: "text/plain",
				},
			},
		},
		{
			name:       "Box-row with anchor tag but no href attribute",
			html:       `<ul><li class="Box-row"><a>no href here</a></li></ul>`,
			wantAssets: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, err := internalgithub.ParseGitHubReleaseAssets(tc.html)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if diff := cmp.Diff(tc.wantAssets, []release.Asset(got)); diff != "" {
				t.Errorf("ParseGitHubReleaseAssets() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

// rewriteTransport redirects all requests to the given test server,
// allowing per-Repository transport injection instead of global mutation.
type rewriteTransport struct {
	server *httptest.Server
}

func (rt *rewriteTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.URL.Scheme = "http"
	req.URL.Host = rt.server.Listener.Addr().String()

	return http.DefaultTransport.RoundTrip(req)
}

// newWebTestRepo creates a Repository whose web scraping transport is
// redirected to the given test server handler. The returned Repository
// is safe for use in parallel tests.
func newWebTestRepo(t *testing.T, handler http.Handler) *internalgithub.Repository {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	repo := internalgithub.NewRepository("owner", "repo", internalgithub.NewClient(""))
	internalgithub.SetTransport(repo, &rewriteTransport{server: server})

	return repo
}

func TestGetReleaseFromWeb(t *testing.T) {
	t.Parallel()

	html := heredoc.Doc(`
		<ul>
		  <li class="Box-row">
		    <a href="/owner/repo/releases/download/v1.2.3/tool-linux-amd64.tar.gz">tool-linux-amd64.tar.gz</a>
		  </li>
		  <li class="Box-row">
		    <a href="/owner/repo/releases/download/v1.2.3/tool-darwin-arm64.tar.gz">tool-darwin-arm64.tar.gz</a>
		  </li>
		</ul>
	`)

	repo := newWebTestRepo(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/owner/repo/releases/expanded_assets/v1.2.3" {
			http.NotFound(w, r)

			return
		}

		w.WriteHeader(http.StatusOK)

		if _, err := w.Write([]byte(html)); err != nil {
			http.Error(w, "write failed", http.StatusInternalServerError)
		}
	}))

	got, err := repo.GetReleaseFromWeb(t.Context(), "v1.2.3")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := &release.Release{
		Tag:  "v1.2.3",
		Name: "v1.2.3",
		Assets: release.Assets{
			{
				Name: "tool-linux-amd64.tar.gz",
				URL:  "https://github.com/owner/repo/releases/download/v1.2.3/tool-linux-amd64.tar.gz",
				Type: "text/plain",
			},
			{
				Name: "tool-darwin-arm64.tar.gz",
				URL:  "https://github.com/owner/repo/releases/download/v1.2.3/tool-darwin-arm64.tar.gz",
				Type: "text/plain",
			},
		},
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("GetReleaseFromWeb() mismatch (-want +got):\n%s", diff)
	}
}

const latestReleasePath = "/owner/repo/releases/latest"

func TestLatestVersionFromWebHTML(t *testing.T) {
	t.Parallel()

	repo := newWebTestRepo(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != latestReleasePath {
			http.NotFound(w, r)

			return
		}

		http.Redirect(w, r, "https://github.com/owner/repo/releases/tag/v3.1.4", http.StatusFound)
	}))

	got, err := repo.LatestVersionFromWebHTML(t.Context())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != "v3.1.4" {
		t.Errorf("LatestVersionFromWebHTML() = %q, want %q", got, "v3.1.4")
	}
}

func TestLatestVersionFromWebHTML_NonRedirectStatus(t *testing.T) {
	t.Parallel()

	repo := newWebTestRepo(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != latestReleasePath {
			http.NotFound(w, r)

			return
		}

		http.Error(w, "internal server error", http.StatusInternalServerError)
	}))

	got, err := repo.LatestVersionFromWebHTML(t.Context())
	if err == nil {
		t.Fatalf("expected error for non-302 response, got tag: %q", got)
	}

	if got != "" {
		t.Errorf("expected empty tag on error, got %q", got)
	}
}

func TestLatestVersionFromWebJSON(t *testing.T) {
	t.Parallel()

	t.Run("happy path returns tag_name", func(t *testing.T) {
		t.Parallel()

		repo := newWebTestRepo(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != latestReleasePath {
				http.NotFound(w, r)

				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			_, _ = w.Write([]byte(`{"tag_name":"v4.2.0"}`))
		}))

		got, err := repo.LatestVersionFromWebJSON(t.Context())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got != "v4.2.0" {
			t.Errorf("LatestVersionFromWebJSON() = %q, want %q", got, "v4.2.0")
		}
	})

	t.Run("empty tag_name returns error", func(t *testing.T) {
		t.Parallel()

		repo := newWebTestRepo(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != latestReleasePath {
				http.NotFound(w, r)

				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			_, _ = w.Write([]byte(`{"tag_name":""}`))
		}))

		got, err := repo.LatestVersionFromWebJSON(t.Context())
		if err == nil {
			t.Fatalf("expected error for empty tag_name, got %q", got)
		}

		if got != "" {
			t.Errorf("expected empty string on error, got %q", got)
		}
	})
}

func TestGetReleaseFromWeb_ServerError(t *testing.T) {
	t.Parallel()

	repo := newWebTestRepo(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}))

	got, err := repo.GetReleaseFromWeb(t.Context(), "v1.0.0")
	if err == nil {
		t.Fatalf("expected error for 500 response, got release: %+v", got)
	}

	if got != nil {
		t.Errorf("expected nil release on error, got %+v", got)
	}
}

func TestLatestVersionFromWebHTML_EmptyLocationHeader(t *testing.T) {
	t.Parallel()

	repo := newWebTestRepo(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != latestReleasePath {
			http.NotFound(w, r)

			return
		}

		// Return 302 but without a Location header.
		w.WriteHeader(http.StatusFound)
	}))

	got, err := repo.LatestVersionFromWebHTML(t.Context())
	if err == nil {
		t.Fatalf("expected error for empty Location header, got tag: %q", got)
	}

	if got != "" {
		t.Errorf("expected empty tag on error, got %q", got)
	}
}
