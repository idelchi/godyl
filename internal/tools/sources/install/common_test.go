package install_test

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/idelchi/godyl/internal/tools/checksum"
	"github.com/idelchi/godyl/internal/tools/mode"
	"github.com/idelchi/godyl/internal/tools/sources/install"
	"github.com/idelchi/godyl/pkg/path/file"
	"github.com/idelchi/godyl/pkg/path/folder"
)

// testData serves a checksummed local download with HEAD and range support.
func testData(t *testing.T, out, name string, body []byte) install.Data {
	t.Helper()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeContent(w, r, name, time.Time{}, bytes.NewReader(body))
	}))
	t.Cleanup(server.Close)

	return install.Data{
		Mode:     mode.Extract.String(),
		Output:   out,
		Path:     server.URL + "/" + name,
		Checksum: checksum.Checksum{Type: checksum.SHA256, Value: fmt.Sprintf("%x", sha256.Sum256(body))},
	}
}

// testTar builds a tar from headers; regular files all contain the supplied text.
func testTar(t *testing.T, content string, headers ...tar.Header) []byte {
	t.Helper()

	var buf bytes.Buffer

	tw := tar.NewWriter(&buf)

	for _, h := range headers {
		if h.Typeflag != tar.TypeDir {
			h.Size = int64(len(content))
		}

		if err := tw.WriteHeader(&h); err != nil {
			t.Fatal(err)
		}

		if h.Size != 0 {
			if _, err := io.WriteString(tw, content); err != nil {
				t.Fatal(err)
			}
		}
	}

	if err := tw.Close(); err != nil {
		t.Fatal(err)
	}

	return buf.Bytes()
}

// expectFile checks the installed bytes and, on Unix, the permission bits.
func expectFile(t *testing.T, path, want string, perm os.FileMode) {
	t.Helper()

	f := file.New(path)

	got, err := f.ReadString()
	if err != nil || got != want {
		t.Fatalf("read %s: got %q, want %q (err %v)", path, got, want, err)
	}

	info, err := f.Info()
	if err != nil {
		t.Fatal(err)
	}

	if runtime.GOOS != "windows" && info.Mode().Perm() != perm {
		t.Errorf("%s: mode %o, want %o", path, info.Mode().Perm(), perm)
	}
}

// TestDownloadExtract verifies archive merging, metadata and cleanup of read-only staging directories.
func TestDownloadExtract(t *testing.T) {
	t.Parallel()

	out := t.TempDir()
	target := filepath.Join(out, "release")
	stamp := time.Unix(1700000000, 0)
	body := testTar(t, "new", tar.Header{Name: "release/", Typeflag: tar.TypeDir, Mode: 0o555},
		tar.Header{Name: "release/tool", Mode: 0o750, ModTime: stamp},
		tar.Header{Name: "empty/", Typeflag: tar.TypeDir, Mode: 0o755})

	if err := os.WriteFile(filepath.Join(out, "other"), []byte("keep"), 0o600); err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() { _ = folder.New(target).Chmod(0o700) })

	_, err := install.Download(testData(t, out, "release.tar", body))
	if err != nil {
		t.Fatal(err)
	}

	expectFile(t, filepath.Join(target, "tool"), "new", 0o750)
	expectFile(t, filepath.Join(out, "other"), "keep", 0o600)

	empty, err := os.Stat(filepath.Join(out, "empty"))
	if err != nil || !empty.IsDir() {
		t.Fatalf("empty directory not installed: %v", err)
	}

	info, err := os.Stat(filepath.Join(target, "tool"))
	if err != nil || !info.ModTime().Equal(stamp) {
		t.Fatalf("archive timestamp not preserved: %v, %v", info, err)
	}
}

// TestDownloadInvalidArchive preserves installed files even if the first tar member was successfully unpacked.
func TestDownloadInvalidArchive(t *testing.T) {
	t.Parallel()

	out := t.TempDir()
	path := filepath.Join(out, "tool")

	if err := file.New(path).Write([]byte("old"), 0o700); err != nil {
		t.Fatal(err)
	}

	body := testTar(t, "new payload", tar.Header{Name: "tool", Mode: 0o755}, tar.Header{Name: "other", Mode: 0o644})
	// Keep the first member intact but truncate the second one inside its payload.
	body = body[:bytes.LastIndex(body, []byte("new payload"))+1]

	_, err := install.Download(testData(t, out, "broken.tar", body))
	if err == nil {
		t.Fatal("truncated tar unexpectedly succeeded")
	}

	expectFile(t, path, "old", 0o700)

	if entries, err := os.ReadDir(out); err != nil || len(entries) != 1 {
		t.Fatalf("unexpected output after failed download: %v (%v)", entries, err)
	}
}

// TestDownloadZipMerge replaces a zip member but leaves the existing parent directory and neighboring file alone.
func TestDownloadZipMerge(t *testing.T) {
	t.Parallel()

	out := folder.New(t.TempDir())
	bin := out.Join("bin")

	if err := bin.Create(0o700); err != nil {
		t.Fatal(err)
	}

	if err := bin.WithFile("keep").Write([]byte("keep")); err != nil {
		t.Fatal(err)
	}

	if err := bin.WithFile("tool").Write([]byte("old")); err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer

	writer := zip.NewWriter(&buf)
	header := &zip.FileHeader{Name: "bin/tool"}
	header.SetMode(0o750)

	member, err := writer.CreateHeader(header)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := io.WriteString(member, "zip"); err != nil {
		t.Fatal(err)
	}

	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}

	if _, err := install.Download(testData(t, out.Path(), "release.zip", buf.Bytes())); err != nil {
		t.Fatal(err)
	}

	expectFile(t, bin.WithFile("tool").Path(), "zip", 0o750)
	expectFile(t, bin.WithFile("keep").Path(), "keep", 0o600)

	info, err := bin.Info()
	if err != nil {
		t.Fatal(err)
	}

	if runtime.GOOS != "windows" && info.Mode().Perm() != 0o700 {
		t.Fatalf("existing directory mode = %o, want 0700", info.Mode().Perm())
	}
}

// TestAcquireReadOnlyRoot verifies cleanup through the same artifact handle retained by the shared cache.
func TestAcquireReadOnlyRoot(t *testing.T) {
	t.Parallel()

	body := testTar(t, "cached", tar.Header{Name: "./", Typeflag: tar.TypeDir, Mode: 0o555},
		tar.Header{Name: "tool", Mode: 0o755})

	a, err := install.Acquire(testData(t, t.TempDir(), "cached.tar", body))
	if err != nil {
		t.Fatal(err)
	}

	if err := a.Remove(); err != nil {
		t.Fatal(err)
	}

	if a.Destination.Exists() {
		t.Fatal("removed artifact still exists")
	}
}

// TestFindSharedArtifact verifies independent installations leave the cached artifact and its permissions untouched.
func TestFindSharedArtifact(t *testing.T) {
	t.Parallel()

	source := filepath.Join(t.TempDir(), "source")
	if err := os.WriteFile(source, []byte("shared"), 0o600); err != nil {
		t.Fatal(err)
	}

	for _, name := range []string{"one", "two"} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			out := t.TempDir()
			if _, err := install.Find(file.New(source), install.Data{Output: out, Exe: name}); err != nil {
				t.Fatal(err)
			}

			expectFile(t, filepath.Join(out, name), "shared", 0o711)
			expectFile(t, source, "shared", 0o600)
		})
	}
}

// TestFindFailedReplacement verifies a non-empty directory target survives without a leftover staging directory.
func TestFindFailedReplacement(t *testing.T) {
	t.Parallel()

	out := t.TempDir()
	if err := os.Mkdir(filepath.Join(out, "tool"), 0o700); err != nil {
		t.Fatal(err)
	}

	keep := file.New(out, "tool", "keep")
	if err := keep.Write([]byte("original")); err != nil {
		t.Fatal(err)
	}

	source := filepath.Join(t.TempDir(), "source")
	if err := os.WriteFile(source, []byte("new"), 0o600); err != nil {
		t.Fatal(err)
	}

	if _, err := install.Find(file.New(source), install.Data{Output: out, Exe: "tool"}); err == nil {
		t.Fatal("unexpectedly replaced a directory with an executable")
	}

	entries, err := os.ReadDir(out)
	if err != nil || len(entries) != 1 || !entries[0].IsDir() || entries[0].Name() != "tool" {
		t.Fatalf("unexpected output after failed replacement: %v (%v)", entries, err)
	}

	expectFile(t, keep.Path(), "original", 0o600)
}

// TestDownloadOutputSymlinks follows existing output file and directory links just like go-getter.
func TestDownloadOutputSymlinks(t *testing.T) {
	t.Parallel()

	if runtime.GOOS == "windows" {
		t.Skip("requires permission to create symlinks")
	}

	out := t.TempDir()
	external := t.TempDir()
	referent := filepath.Join(external, "tool")

	if err := os.WriteFile(referent, []byte("original"), 0o600); err != nil {
		t.Fatal(err)
	}

	if err := os.Symlink(referent, filepath.Join(out, "tool")); err != nil {
		t.Fatal(err)
	}

	if _, err := install.Download(testData(t, out, "tool", []byte("replacement"))); err != nil {
		t.Fatal(err)
	}

	expectFile(t, referent, "replacement", 0o600)

	if got, err := os.Readlink(filepath.Join(out, "tool")); err != nil || got != referent {
		t.Fatalf("output symlink was replaced: got %q, err %v", got, err)
	}

	if err := os.Symlink(external, filepath.Join(out, "linked")); err != nil {
		t.Fatal(err)
	}

	body := testTar(t, "upgrade", tar.Header{Name: "linked/tool", Mode: 0o600})

	_, err := install.Download(testData(t, out, "link.tar", body))
	if err != nil {
		t.Fatal(err)
	}

	expectFile(t, referent, "upgrade", 0o600)
}
