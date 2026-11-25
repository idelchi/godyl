package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"slices"
	"strings"
	"testing"

	"github.com/idelchi/godyl/pkg/path/file"
)

// TestMain runs the CLI when invoked as a child process by runCLI.
func TestMain(m *testing.M) {
	if len(os.Args) > 1 && os.Args[1] == "--godyl-test-cli" {
		os.Args = append(os.Args[:1], os.Args[2:]...)

		main()

		return
	}

	m.Run()
}

// runCLI runs godyl with an isolated home and a manifest on stdin.
func runCLI(t *testing.T, dir, manifest string, want int, args ...string) string {
	t.Helper()

	bin, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}

	flags := []string{"--godyl-test-cli", "--cache-dir", "cache", "--no-progress"}
	//nolint:gosec // Runs this test executable, not a user-supplied command.
	cmd := exec.CommandContext(t.Context(), bin, slices.Concat(flags, args)...)

	cmd.Dir = dir
	cmd.Env = []string{
		"HOME=" + dir, "USERPROFILE=" + dir, "TEMP=" + dir, "TMP=" + dir, "TMPDIR=" + dir,
		"PATH=" + os.Getenv("PATH"), "SystemRoot=" + os.Getenv("SystemRoot"), "GOCOVERDIR=" + dir,
	}
	cmd.Stdin = strings.NewReader(manifest)

	var stderr bytes.Buffer

	cmd.Stderr = &stderr

	out, err := cmd.Output()
	if cmd.ProcessState == nil {
		t.Fatalf("starting %v: %v", args, err)
	}

	if got := cmd.ProcessState.ExitCode(); got != want {
		t.Fatalf("%v: want exit %d, got %d\n%s%s", args, want, got, out, stderr.String())
	}

	if want == 2 && (!bytes.Contains(out, []byte("requires download")) || stderr.Len() != 0) {
		t.Fatalf("%v: expected a pending change without an error\n%s%s", args, out, stderr.String())
	}

	return string(out)
}

// TestStatusAndDryRun checks exit codes and cache state before and after a real local install.
func TestStatusAndDryRun(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(w, "fixture payload")
	}))
	t.Cleanup(server.Close)

	const manifest = `- name: fixture
  exe: fixture
  source:
    type: url
  url: %s/fixture{{ .EXTENSION }}
  checksum: none
  version: %s
`

	dir := t.TempDir()
	cacheFile := file.New(dir, "cache", "godyl.json")
	v1 := fmt.Sprintf(manifest, server.URL, "v1.0.0")
	v2 := fmt.Sprintf(manifest, server.URL, "v2.0.0")

	runCLI(t, dir, v1, 2, "status", "-", "--output", "bin")
	runCLI(t, dir, v1, 0, "install", "-", "--output", "bin", "--dry")
	runCLI(t, dir, v1, 0, "status", "-", "--output", "bin", "--tags", "unselected")

	if _, err := os.Stat(cacheFile.Path()); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("preview created a cache or failed to read it: %v", err)
	}

	runCLI(t, dir, v1, 0, "install", "-", "--output", "bin")

	installed, err := cacheFile.Read()
	if err != nil {
		t.Fatal(err)
	}

	runCLI(t, dir, v1, 0, "status", "-", "--output", "bin")

	preview := runCLI(t, dir, v1, 0, "install", "-", "--output", "bin", "--dry")
	if !strings.Contains(preview, `target version "v1.0.0" match`) {
		t.Fatalf("dry run did not skip the installed version:\n%s", preview)
	}

	runCLI(t, dir, v2, 2, "status", "-", "--output", "bin")

	preview = runCLI(t, dir, v2, 0, "install", "-", "--output", "bin", "--dry")
	if !strings.Contains(preview, "requires download") {
		t.Fatalf("dry run did not report the pending update:\n%s", preview)
	}

	runCLI(t, dir, v2, 2, "status", "-", "--output", "bin")

	current, err := cacheFile.Read()
	if err != nil || !bytes.Equal(current, installed) {
		t.Fatalf("preview changed installed-version cache: %s, error: %v", current, err)
	}

	// The first tool is pending but the second cannot resolve its output template.
	mixed := v2 + "- name: broken\n  source: {type: url}\n  output: '{{ .Unknown }}'\n"
	out := runCLI(t, dir, mixed, 1, "status", "-")

	if !strings.Contains(out, "requires download") || !strings.Contains(out, "templating pre api") {
		t.Fatalf("expected both a pending tool and a failed tool:\n%s", out)
	}
}
