package download_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/idelchi/godyl/pkg/download"
)

func TestURLWithChecksum(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		url   string
		query string
		want  string
	}{
		{
			name:  "plain URL gets query param appended",
			url:   "https://example.com/file.tar.gz",
			query: "checksum=sha256:abc123",
			want:  "https://example.com/file.tar.gz?checksum=sha256:abc123",
		},
		{
			name:  "URL with trailing slash and checksum",
			url:   "https://example.com/releases/",
			query: "checksum=md5:deadbeef",
			want:  "https://example.com/releases/?checksum=md5:deadbeef",
		},
		{
			name:  "URL with path segments and no params",
			url:   "https://example.com/owner/repo/releases/download/v1.0.0/binary",
			query: "checksum=sha256:0000",
			want:  "https://example.com/owner/repo/releases/download/v1.0.0/binary?checksum=sha256:0000",
		},
		{
			name:  "URL with existing query appends with ampersand",
			url:   "https://example.com/file.tar.gz?foo=bar",
			query: "checksum=sha256:abc123",
			want:  "https://example.com/file.tar.gz?foo=bar&checksum=sha256:abc123",
		},
		{
			name:  "URL with multiple existing params",
			url:   "https://example.com/file.zip?token=xyz&version=1",
			query: "checksum=sha256:abc",
			want:  "https://example.com/file.zip?token=xyz&version=1&checksum=sha256:abc",
		},
		{
			name:  "URL with existing empty value param",
			url:   "https://example.com/file?key=",
			query: "checksum=sha256:111",
			want:  "https://example.com/file?key=&checksum=sha256:111",
		},
		{
			name:  "empty query returns URL unchanged",
			url:   "https://example.com/file.tar.gz",
			query: "",
			want:  "https://example.com/file.tar.gz",
		},
		{
			name:  "empty query with existing params returns URL unchanged",
			url:   "https://example.com/file.tar.gz?existing=param",
			query: "",
			want:  "https://example.com/file.tar.gz?existing=param",
		},
		{
			name:  "URL already contains the same checksum param",
			url:   "https://example.com/file.tar.gz?checksum=sha256:abc123",
			query: "checksum=sha256:abc123",
			want:  "https://example.com/file.tar.gz?checksum=sha256:abc123&checksum=sha256:abc123",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := download.URLWithChecksum(tc.url, tc.query)

			if got != tc.want {
				t.Errorf("URLWithChecksum(%q, %q) = %q, want %q", tc.url, tc.query, got, tc.want)
			}
		})
	}
}

func TestDownload(t *testing.T) {
	t.Parallel()

	const body = "hello download"

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(w, body)
	}))
	t.Cleanup(srv.Close)

	notFoundSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	t.Cleanup(notFoundSrv.Close)

	t.Run("happy path", func(t *testing.T) {
		t.Parallel()

		dst := filepath.Join(t.TempDir(), "out.txt")

		d := download.New(download.WithContextTimeout(10 * time.Second))

		f, err := d.Download(srv.URL, dst)
		if err != nil {
			t.Fatalf("Download(%q): unexpected error: %v", srv.URL, err)
		}

		if string(f) != dst {
			t.Errorf("Download(%q): returned path = %q, want %q", srv.URL, string(f), dst)
		}

		got, err := os.ReadFile(string(f))
		if err != nil {
			t.Fatalf("ReadFile(%q): %v", string(f), err)
		}

		if string(got) != body {
			t.Errorf("downloaded content = %q, want %q", string(got), body)
		}
	})

	t.Run("404 returns error", func(t *testing.T) {
		t.Parallel()

		dst := filepath.Join(t.TempDir(), "out.txt")

		d := download.New(download.WithContextTimeout(10*time.Second), download.WithMaxRetries(0))

		_, err := d.Download(notFoundSrv.URL, dst)
		if err == nil {
			t.Fatal("Download(404): expected error, got nil")
		}
	})

	t.Run("invalid URL returns error", func(t *testing.T) {
		t.Parallel()

		d := download.New(download.WithContextTimeout(5 * time.Second))

		_, err := d.Download("not-a-url", t.TempDir())
		if err == nil {
			t.Fatal("Download(invalid URL): expected error, got nil")
		}
	})
}

func TestDownload500(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	t.Cleanup(srv.Close)

	dst := filepath.Join(t.TempDir(), "out.txt")

	// Zero retries so the test does not block waiting for retry back-off.
	d := download.New(
		download.WithContextTimeout(10*time.Second),
		download.WithMaxRetries(0),
	)

	_, err := d.Download(srv.URL, dst)
	if err == nil {
		t.Fatal("Download(500): expected error, got nil")
	}
}

func TestDownloadToSubdir(t *testing.T) {
	t.Parallel()

	const body = "subdir content"

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(w, body)
	}))
	t.Cleanup(srv.Close)

	// Destination is a file inside a nested subdirectory of the temp dir.
	dst := filepath.Join(t.TempDir(), "nested", "subdir", "output.txt")

	d := download.New(download.WithContextTimeout(10 * time.Second))

	f, err := d.Download(srv.URL, dst)
	if err != nil {
		t.Fatalf("Download(subdir): unexpected error: %v", err)
	}

	got, err := os.ReadFile(string(f))
	if err != nil {
		t.Fatalf("ReadFile(%q): %v", string(f), err)
	}

	if string(got) != body {
		t.Errorf("Download(subdir): content = %q, want %q", string(got), body)
	}
}
