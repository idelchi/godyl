//nolint:testpackage // Tests exercise unexported artifact key helpers.
package processor

import (
	"archive/zip"
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"sync/atomic"
	"testing"

	"github.com/idelchi/godyl/internal/config/root"
	"github.com/idelchi/godyl/internal/detect"
	"github.com/idelchi/godyl/internal/detect/platform"
	"github.com/idelchi/godyl/internal/tools"
	"github.com/idelchi/godyl/internal/tools/checksum"
	"github.com/idelchi/godyl/internal/tools/exe"
	"github.com/idelchi/godyl/internal/tools/mode"
	"github.com/idelchi/godyl/internal/tools/sources"
	"github.com/idelchi/godyl/internal/tools/sources/install"
	"github.com/idelchi/godyl/internal/tools/strategy"
	"github.com/idelchi/godyl/internal/tools/tags"
	"github.com/idelchi/godyl/internal/tools/tool"
	"github.com/idelchi/godyl/pkg/logger"
)

// TestArtifactKey verifies artifact key behavior.
func TestArtifactKey(t *testing.T) {
	t.Parallel()

	base := install.Data{
		Path: "https://example.com/tool.zip",
		Header: http.Header{
			"x-token": []string{"b", "a"},
		},
		Checksum: checksum.Checksum{
			Type:  checksum.SHA256,
			Value: "abc123",
		},
	}

	same := base

	same.Header = http.Header{
		"X-Token": []string{"a", "b"},
	}

	key, err := newArtifactKey(sources.URL, base)
	if err != nil {
		t.Fatalf("newArtifactKey(): %v", err)
	}

	sameKey, err := newArtifactKey(sources.URL, same)
	if err != nil {
		t.Fatalf("newArtifactKey() with reordered headers: %v", err)
	}

	if key != sameKey {
		t.Fatal("newArtifactKey() should normalize header case and value order")
	}

	differentChecksum := base

	differentChecksum.Checksum.Value = "def456"

	checksumKey, err := newArtifactKey(sources.URL, differentChecksum)
	if err != nil {
		t.Fatalf("newArtifactKey() with different checksum: %v", err)
	}

	if key == checksumKey {
		t.Fatal("newArtifactKey() should differ when checksum differs")
	}

	differentHeader := base

	differentHeader.Header = http.Header{"X-Token": []string{"c"}}

	headerKey, err := newArtifactKey(sources.URL, differentHeader)
	if err != nil {
		t.Fatalf("newArtifactKey() with different header: %v", err)
	}

	if key == headerKey {
		t.Fatal("newArtifactKey() should differ when headers differ")
	}

	differentURL := base

	differentURL.Path = "https://example.com/other.zip"

	urlKey, err := newArtifactKey(sources.URL, differentURL)
	if err != nil {
		t.Fatalf("newArtifactKey() with different URL: %v", err)
	}

	if key == urlKey {
		t.Fatal("newArtifactKey() should differ when URL differs")
	}

	sourceKey, err := newArtifactKey(sources.GITHUB, base)
	if err != nil {
		t.Fatalf("newArtifactKey() with different source: %v", err)
	}

	if key == sourceKey {
		t.Fatal("newArtifactKey() should differ when source differs")
	}

	noVerifySSL := base

	noVerifySSL.NoVerifySSL = true

	noVerifySSLKey, err := newArtifactKey(sources.URL, noVerifySSL)
	if err != nil {
		t.Fatalf("newArtifactKey() with no-verify-ssl: %v", err)
	}

	if key == noVerifySSLKey {
		t.Fatal("newArtifactKey() should differ when no-verify-ssl differs")
	}

	noVerifyChecksum := base

	noVerifyChecksum.NoVerifyChecksum = true

	noVerifyChecksumKey, err := newArtifactKey(sources.URL, noVerifyChecksum)
	if err != nil {
		t.Fatalf("newArtifactKey() with no-verify-checksum: %v", err)
	}

	if key == noVerifyChecksumKey {
		t.Fatal("newArtifactKey() should differ when no-verify-checksum differs")
	}
}

// TestIsReusableArtifact verifies is reusable artifact behavior.
func TestIsReusableArtifact(t *testing.T) {
	t.Parallel()

	const artifactURL = "https://example.com/tool.zip"

	tests := []struct {
		name string
		tool *tool.Tool
		want bool
	}{
		{
			name: "url find is reusable",
			tool: testURLTool(t, "url-find", artifactURL, "tool-a", t.TempDir(), mode.Find),
			want: true,
		},
		{
			name: "url extract is not reusable",
			tool: testURLTool(t, "url-extract", artifactURL, "tool-a", t.TempDir(), mode.Extract),
			want: false,
		},
		{
			name: "go find is not reusable",
			tool: func() *tool.Tool {
				tool := testURLTool(t, "go-find", artifactURL, "tool-a", t.TempDir(), mode.Find)

				tool.Source.Type = sources.GO

				return tool
			}(),
			want: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if got := isReusableArtifact(tc.tool); got != tc.want {
				t.Errorf("isReusableArtifact() = %v, want %v", got, tc.want)
			}
		})
	}
}

// TestProcessorReusesURLArtifactForFindMode verifies processor reuses url artifact for find mode behavior.
func TestProcessorReusesURLArtifactForFindMode(t *testing.T) {
	t.Parallel()

	const (
		toolA = "tool-a"
		toolB = "tool-b"
	)

	archive := testZip(t, map[string]string{
		toolA: "alpha",
		toolB: "beta",
	})

	var getCount atomic.Int64

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/zip")

		if r.Method == http.MethodHead {
			return
		}

		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)

			return
		}

		getCount.Add(1)

		_, _ = w.Write(archive)
	}))
	t.Cleanup(server.Close)

	output := t.TempDir()
	artifactURL := server.URL + "/bundle.zip"

	processor := New(tools.Tools{
		testURLTool(t, "jsonnet", artifactURL, toolA, output, mode.Find),
		testURLTool(t, "jsonnet-lint", artifactURL, toolB, output, mode.Find),
	}, root.Config{
		Cache:      root.Cache{Disabled: true},
		NoProgress: true,
		Parallel:   2,
		Tokens:     root.Tokens{GitHub: "test-token"},
	}, testLogger(t))

	summary, err := processor.Process(tags.IncludeTags{})
	if err != nil {
		t.Fatalf("Process(): %v", err)
	}

	if summary.Successful != 2 || summary.Failed != 0 {
		t.Fatalf(
			"Process() summary = successful %d failed %d errors %v",
			summary.Successful,
			summary.Failed,
			summary.Errors,
		)
	}

	if got := getCount.Load(); got != 1 {
		t.Fatalf("GET request count = %d, want 1", got)
	}

	assertFileContent(t, filepath.Join(output, toolA), "alpha")
	assertFileContent(t, filepath.Join(output, toolB), "beta")
}

// testURLTool constructs a URL-backed tool configured for artifact-reuse tests.
func testURLTool(t *testing.T, name, url, exeName, output string, mode mode.Mode) *tool.Tool {
	t.Helper()

	patterns := exe.Patterns{exeName}
	tool := tool.NewEmptyTool()

	tool.Name = name
	tool.URL = url
	tool.Output = output
	tool.Exe = exe.Exe{
		Name:     exeName,
		Patterns: &patterns,
	}
	tool.Platform = detect.Platform{
		OS: platform.OS{
			Name: runtime.GOOS,
		},
		Architecture: platform.Architecture{
			Name: runtime.GOARCH,
		},
	}
	tool.Source.Type = sources.URL
	tool.Strategy = strategy.Force
	tool.Mode = mode
	tool.Checksum.Type = checksum.None

	return tool
}

// testZip builds an in-memory ZIP archive containing the supplied files.
func testZip(t *testing.T, files map[string]string) []byte {
	t.Helper()

	var buffer bytes.Buffer

	writer := zip.NewWriter(&buffer)

	for name, content := range files {
		header := &zip.FileHeader{
			Name: name,
		}
		header.SetMode(0o755)

		fileWriter, err := writer.CreateHeader(header)
		if err != nil {
			t.Fatalf("CreateHeader(%q): %v", name, err)
		}

		if _, err := io.WriteString(fileWriter, content); err != nil {
			t.Fatalf("WriteString(%q): %v", name, err)
		}
	}

	if err := writer.Close(); err != nil {
		t.Fatalf("Close zip writer: %v", err)
	}

	return buffer.Bytes()
}

// testLogger returns a silent logger suitable for processor tests.
func testLogger(t *testing.T) *logger.Logger {
	t.Helper()

	log, err := logger.NewCustom(logger.SILENT, io.Discard)
	if err != nil {
		t.Fatalf("NewCustom(): %v", err)
	}

	return log
}

// assertFileContent verifies the complete contents of path.
func assertFileContent(t *testing.T, path, want string) {
	t.Helper()

	got, err := os.ReadFile(path) //nolint:gosec // Test helper reads files created by this test.
	if err != nil {
		t.Fatalf("ReadFile(%q): %v", path, err)
	}

	if string(got) != want {
		t.Fatalf("ReadFile(%q) = %q, want %q", path, string(got), want)
	}
}
