package file_test

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/idelchi/godyl/pkg/path/file"
)

// ---------------------------------------------------------------------------
// Section 1: Pure path tests (no filesystem)
// ---------------------------------------------------------------------------

func TestNew(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		paths []string
		want  string
	}{
		{
			name:  "single component",
			paths: []string{"foo"},
			want:  "foo",
		},
		{
			name:  "two components",
			paths: []string{"a", "b"},
			want:  "a/b",
		},
		{
			name:  "three components",
			paths: []string{"a", "b", "c"},
			want:  "a/b/c",
		},
		{
			name:  "absolute path with multiple components",
			paths: []string{"/usr", "local", "bin"},
			want:  "/usr/local/bin",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := file.New(tc.paths...)
			if got.Path() != tc.want {
				t.Errorf("New(%v).Path() = %q, want %q", tc.paths, got.Path(), tc.want)
			}
		})
	}
}

func TestBase(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "path with directory components and extension",
			input: "a/b/c.txt",
			want:  "c.txt",
		},
		{
			name:  "filename only",
			input: "file.txt",
			want:  "file.txt",
		},
		{
			name:  "absolute path no extension",
			input: "/usr/local/bin/tool",
			want:  "tool",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := file.New(tc.input).Base()
			if got != tc.want {
				t.Errorf("New(%q).Base() = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}

func TestExtension(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "single extension",
			input: "file.txt",
			want:  "txt",
		},
		{
			name:  "double extension returns last segment",
			input: "file.tar.gz",
			want:  "gz",
		},
		{
			name:  "no extension",
			input: "file",
			want:  "",
		},
		{
			name:  "hidden file treated as having extension",
			input: ".hidden",
			want:  "hidden",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := file.New(tc.input).Extension()
			if got != tc.want {
				t.Errorf("New(%q).Extension() = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}

func TestWithoutExtension(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "removes single extension",
			input: "file.txt",
			want:  "file",
		},
		{
			name:  "removes only last extension",
			input: "file.tar.gz",
			want:  "file.tar",
		},
		{
			name:  "no extension unchanged",
			input: "file",
			want:  "file",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := file.New(tc.input).WithoutExtension().Path()
			if got != tc.want {
				t.Errorf("New(%q).WithoutExtension().Path() = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}

func TestWithoutExtensions(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "strips all extensions from double-extension file",
			input: "file.tar.gz",
			want:  "file",
		},
		{
			name:  "strips all extensions from triple-extension file",
			input: "file.a.b.c",
			want:  "file",
		},
		{
			name:  "no extension unchanged",
			input: "file",
			want:  "file",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := file.New(tc.input).WithoutExtensions().Path()
			if got != tc.want {
				t.Errorf("New(%q).WithoutExtensions().Path() = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}

func TestHasExtension(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{
			name:  "file with extension",
			input: "file.txt",
			want:  true,
		},
		{
			name:  "file without extension",
			input: "file",
			want:  false,
		},
		{
			name:  "hidden file has extension",
			input: ".hidden",
			want:  true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := file.New(tc.input).HasExtension()
			if got != tc.want {
				t.Errorf("New(%q).HasExtension() = %v, want %v", tc.input, got, tc.want)
			}
		})
	}
}

func TestWithExtension(t *testing.T) {
	t.Parallel()

	t.Run("adds extension to file without one", func(t *testing.T) {
		t.Parallel()

		got := file.New("path/to/file").WithExtension("txt").Path()
		want := "path/to/file.txt"

		if got != want {
			t.Errorf("WithExtension(\"txt\").Path() = %q, want %q", got, want)
		}
	})

	t.Run("replaces existing extension", func(t *testing.T) {
		t.Parallel()

		got := file.New("path/to/file.old").WithExtension("new").Path()
		want := "path/to/file.new"

		if got != want {
			t.Errorf("WithExtension(\"new\").Path() = %q, want %q", got, want)
		}
	})
}

func TestUnescape(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "decodes percent-encoded space",
			input: "file%20name.txt",
			want:  "file name.txt",
		},
		{
			name:  "no encoding unchanged",
			input: "file.txt",
			want:  "file.txt",
		},
		{
			name:  "decodes percent-encoded plus sign",
			input: "a%2Bb.txt",
			want:  "a+b.txt",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := file.New(tc.input).Unescape().Path()
			if got != tc.want {
				t.Errorf("New(%q).Unescape().Path() = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}

func TestMatches(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		pattern string
		want    bool
	}{
		{
			name:    "double-extension glob match",
			input:   "file.tar.gz",
			pattern: "*.tar.gz",
			want:    true,
		},
		{
			name:    "doublestar glob matches nested go file",
			input:   "a/b/c/file.go",
			pattern: "**/*.go",
			want:    true,
		},
		{
			name:    "extension mismatch",
			input:   "file.txt",
			pattern: "*.go",
			want:    false,
		},
		{
			name:    "exact filename match",
			input:   "file.txt",
			pattern: "file.txt",
			want:    true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, err := file.New(tc.input).Matches(tc.pattern)
			if err != nil {
				t.Fatalf("New(%q).Matches(%q) unexpected error: %v", tc.input, tc.pattern, err)
			}

			if got != tc.want {
				t.Errorf("New(%q).Matches(%q) = %v, want %v", tc.input, tc.pattern, got, tc.want)
			}
		})
	}
}

func TestWithoutFolder(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		input  string
		prefix string
		want   string
	}{
		{
			name:   "strips matching leading folder",
			input:  "a/b/c.txt",
			prefix: "a",
			want:   "b/c.txt",
		},
		{
			name:   "non-matching prefix leaves path unchanged",
			input:  "a/b/c.txt",
			prefix: "x",
			want:   "a/b/c.txt",
		},
		{
			name:   "empty prefix leaves path unchanged",
			input:  "a/b",
			prefix: "",
			want:   "a/b",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := file.New(tc.input).WithoutFolder(tc.prefix).Path()
			if got != tc.want {
				t.Errorf("New(%q).WithoutFolder(%q).Path() = %q, want %q", tc.input, tc.prefix, got, tc.want)
			}
		})
	}
}

func TestIsAbs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{
			name:  "absolute path",
			input: "/absolute/path",
			want:  true,
		},
		{
			name:  "relative path",
			input: "relative/path",
			want:  false,
		},
		{
			name:  "single dot",
			input: ".",
			want:  false,
		},
		{
			name:  "root",
			input: "/",
			want:  true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := file.New(tc.input).IsAbs()
			if got != tc.want {
				t.Errorf("New(%q).IsAbs() = %v, want %v", tc.input, got, tc.want)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Section 2: Filesystem tests (use t.TempDir())
// ---------------------------------------------------------------------------

func TestFileCreateAndExists(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	f := file.New(filepath.Join(dir, "newfile.txt"))

	if f.Exists() {
		t.Fatal("file should not exist before Create()")
	}

	if err := f.Create(); err != nil {
		t.Fatalf("Create() unexpected error: %v", err)
	}

	if !f.Exists() {
		t.Error("Exists() = false after Create(), want true")
	}

	if !f.IsFile() {
		t.Error("IsFile() = false after Create(), want true")
	}

	if f.IsDir() {
		t.Error("IsDir() = true after Create(), want false")
	}
}

func TestFileWriteAndRead(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	f := file.New(filepath.Join(dir, "data.bin"))
	content := []byte("hello world")

	if err := f.Write(content); err != nil {
		t.Fatalf("Write() unexpected error: %v", err)
	}

	got, err := f.Read()
	if err != nil {
		t.Fatalf("Read() unexpected error: %v", err)
	}

	if string(got) != string(content) {
		t.Errorf("Read() = %q, want %q", got, content)
	}
}

func TestFileReadString(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	f := file.New(filepath.Join(dir, "str.txt"))
	content := "hello"

	if err := f.Write([]byte(content)); err != nil {
		t.Fatalf("Write() unexpected error: %v", err)
	}

	got, err := f.ReadString()
	if err != nil {
		t.Fatalf("ReadString() unexpected error: %v", err)
	}

	if got != content {
		t.Errorf("ReadString() = %q, want %q", got, content)
	}
}

func TestFileLines(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		content string
		want    []string
	}{
		{
			// No trailing newline: splitting "line1\nline2\nline3" on "\n"
			// yields exactly ["line1", "line2", "line3"] — 3 elements.
			name:    "no trailing newline yields exact line count",
			content: "line1\nline2\nline3",
			want:    []string{"line1", "line2", "line3"},
		},
		{
			// Trailing newline: splitting "line1\nline2\n" on "\n" yields
			// ["line1", "line2", ""] — 3 elements including a final empty string.
			// This mirrors how strings.Split behaves and is the expected contract.
			name:    "trailing newline produces extra empty element",
			content: "line1\nline2\n",
			want:    []string{"line1", "line2", ""},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			dir := t.TempDir()
			f := file.New(filepath.Join(dir, tc.name+".txt"))

			if err := f.Write([]byte(tc.content)); err != nil {
				t.Fatalf("Write() unexpected error: %v", err)
			}

			got, err := f.Lines()
			if err != nil {
				t.Fatalf("Lines() unexpected error: %v", err)
			}

			if !slices.Equal(got, tc.want) {
				t.Errorf("Lines() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestFileCopy(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	src := file.New(filepath.Join(dir, "src.txt"))
	dst := file.New(filepath.Join(dir, "dst.txt"))
	content := []byte("content")

	if err := src.Write(content); err != nil {
		t.Fatalf("Write() unexpected error: %v", err)
	}

	if err := src.Copy(dst); err != nil {
		t.Fatalf("Copy() unexpected error: %v", err)
	}

	got, err := dst.Read()
	if err != nil {
		t.Fatalf("Read() on copy destination unexpected error: %v", err)
	}

	if string(got) != string(content) {
		t.Errorf("Read() after Copy() = %q, want %q", got, content)
	}
}

func TestFileCopy_SameDestination(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	src := file.New(filepath.Join(dir, "same.txt"))
	content := []byte("data")

	if err := src.Write(content); err != nil {
		t.Fatalf("Write() unexpected error: %v", err)
	}

	// Copying a file to itself must return an error.
	if err := src.Copy(src); err == nil {
		t.Error("Copy(self) = nil, want error when source and destination are identical")
	}

	// The original file content must survive the failed self-copy.
	got, err := src.Read()
	if err != nil {
		t.Fatalf("Read() after failed Copy(self) unexpected error: %v", err)
	}

	if string(got) != string(content) {
		t.Errorf("file content after Copy(self) = %q, want %q", got, content)
	}
}

func TestFileRemove(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	f := file.New(filepath.Join(dir, "todelete.txt"))

	if err := f.Create(); err != nil {
		t.Fatalf("Create() unexpected error: %v", err)
	}

	if !f.Exists() {
		t.Fatal("Exists() = false after Create(), want true")
	}

	if err := f.Remove(); err != nil {
		t.Fatalf("Remove() unexpected error: %v", err)
	}

	if f.Exists() {
		t.Error("Exists() = true after Remove(), want false")
	}
}

func TestFileSize(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	f := file.New(filepath.Join(dir, "sized.txt"))
	content := []byte("12345")

	if err := f.Write(content); err != nil {
		t.Fatalf("Write() unexpected error: %v", err)
	}

	got, err := f.Size()
	if err != nil {
		t.Fatalf("Size() unexpected error: %v", err)
	}

	if got != int64(len(content)) {
		t.Errorf("Size() = %d, want %d", got, len(content))
	}
}

func TestFileHash(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	f := file.New(filepath.Join(dir, "hashed.txt"))
	content := []byte("hello")

	if err := f.Write(content); err != nil {
		t.Fatalf("Write() unexpected error: %v", err)
	}

	got, err := f.Hash()
	if err != nil {
		t.Fatalf("Hash() unexpected error: %v", err)
	}

	// Compute expected hash inline to avoid any hardcoded string mismatch.
	sum := sha256.Sum256(content)
	want := hex.EncodeToString(sum[:])

	if got != want {
		t.Errorf("Hash() = %q, want %q", got, want)
	}
}

func TestFileIsExecutable(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	f := file.New(filepath.Join(dir, "script.sh"))

	if err := f.Write([]byte("#!/bin/sh\n")); err != nil {
		t.Fatalf("Write() unexpected error: %v", err)
	}

	// Ensure the file starts without execute permission by setting mode explicitly.
	if err := os.Chmod(f.Path(), 0o600); err != nil {
		t.Fatalf("Chmod() unexpected error: %v", err)
	}

	isExec, err := f.IsExecutable()
	if err != nil {
		t.Fatalf("IsExecutable() unexpected error before MakeExecutable: %v", err)
	}

	if isExec {
		t.Error("IsExecutable() = true before MakeExecutable(), want false")
	}

	if err := f.MakeExecutable(); err != nil {
		t.Fatalf("MakeExecutable() unexpected error: %v", err)
	}

	isExec, err = f.IsExecutable()
	if err != nil {
		t.Fatalf("IsExecutable() unexpected error after MakeExecutable: %v", err)
	}

	if !isExec {
		t.Error("IsExecutable() = false after MakeExecutable(), want true")
	}
}

func TestFileLinks(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	src := file.New(filepath.Join(dir, "original.txt"))
	lnk := file.New(filepath.Join(dir, "link.txt"))
	content := []byte("data")

	if err := src.Write(content); err != nil {
		t.Fatalf("Write() unexpected error: %v", err)
	}

	if err := src.Links(lnk); err != nil {
		t.Fatalf("Links() unexpected error: %v", err)
	}

	if !lnk.Exists() {
		t.Fatal("link Exists() = false after Links(), want true")
	}

	got, err := lnk.Read()
	if err != nil {
		t.Fatalf("Read() on link unexpected error: %v", err)
	}

	if string(got) != string(content) {
		t.Errorf("Read() via link = %q, want %q", got, content)
	}
}

func TestFileNumberOfLines(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	f := file.New(filepath.Join(dir, "numlines.txt"))

	// "a\nb\nc\n" — strings.Split on "\n" yields ["a","b","c",""] = 4 elements.
	if err := f.Write([]byte("a\nb\nc\n")); err != nil {
		t.Fatalf("Write() unexpected error: %v", err)
	}

	got, err := f.NumberOfLines()
	if err != nil {
		t.Fatalf("NumberOfLines() unexpected error: %v", err)
	}

	const want = 4 // ["a", "b", "c", ""] from strings.Split("a\nb\nc\n", "\n")

	if got != want {
		t.Errorf("NumberOfLines() = %d, want %d", got, want)
	}
}

func TestFileSet(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{
			name:  "non-empty path is set",
			input: "x",
			want:  true,
		},
		{
			name:  "empty path is not set",
			input: "",
			want:  false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := file.New(tc.input).Set()
			if got != tc.want {
				t.Errorf("New(%q).Set() = %v, want %v", tc.input, got, tc.want)
			}
		})
	}
}

func TestFileLargerThan(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	f := file.New(filepath.Join(dir, "data.bin"))
	content := []byte("hello") // 5 bytes

	if err := f.Write(content); err != nil {
		t.Fatalf("Write() unexpected error: %v", err)
	}

	size := int64(len(content))

	tests := []struct {
		name      string
		threshold int64
		want      bool
	}{
		{
			name:      "threshold equals file size returns false",
			threshold: size,
			want:      false,
		},
		{
			name:      "threshold one below file size returns true",
			threshold: size - 1,
			want:      true,
		},
		{
			name:      "threshold above file size returns false",
			threshold: size + 1,
			want:      false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, err := f.LargerThan(tc.threshold)
			if err != nil {
				t.Fatalf("LargerThan(%d) unexpected error: %v", tc.threshold, err)
			}

			if got != tc.want {
				t.Errorf("LargerThan(%d) = %v, want %v", tc.threshold, got, tc.want)
			}
		})
	}
}

func TestFileSmallerThan(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	f := file.New(filepath.Join(dir, "data.bin"))
	content := []byte("hello") // 5 bytes

	if err := f.Write(content); err != nil {
		t.Fatalf("Write() unexpected error: %v", err)
	}

	size := int64(len(content))

	tests := []struct {
		name      string
		threshold int64
		want      bool
	}{
		{
			name:      "threshold equals file size returns false",
			threshold: size,
			want:      false,
		},
		{
			name:      "threshold one above file size returns true",
			threshold: size + 1,
			want:      true,
		},
		{
			name:      "threshold below file size returns false",
			threshold: size - 1,
			want:      false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, err := f.SmallerThan(tc.threshold)
			if err != nil {
				t.Fatalf("SmallerThan(%d) unexpected error: %v", tc.threshold, err)
			}

			if got != tc.want {
				t.Errorf("SmallerThan(%d) = %v, want %v", tc.threshold, got, tc.want)
			}
		})
	}
}

func TestFileUp(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "three-component path returns two-component parent",
			input: "a/b/c",
			want:  "a/b",
		},
		{
			name:  "two-component path returns single-component parent",
			input: "x/y",
			want:  "x",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := file.New(tc.input).Up().Path()
			if got != tc.want {
				t.Errorf("New(%q).Up().Path() = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}

func TestFileCopies(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	src := file.New(filepath.Join(dir, "src.txt"))
	dst1 := file.New(filepath.Join(dir, "dst1.txt"))
	dst2 := file.New(filepath.Join(dir, "dst2.txt"))
	content := []byte("copies content")

	if err := src.Write(content); err != nil {
		t.Fatalf("Write() unexpected error: %v", err)
	}

	if err := src.Copies(dst1, dst2); err != nil {
		t.Fatalf("Copies() unexpected error: %v", err)
	}

	for _, dst := range []file.File{dst1, dst2} {
		got, err := dst.Read()
		if err != nil {
			t.Fatalf("Read() on %q unexpected error: %v", dst.Path(), err)
		}

		if string(got) != string(content) {
			t.Errorf("Read() from %q = %q, want %q", dst.Path(), got, content)
		}
	}
}

func TestFileCopyToNonExistentDir(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	src := file.New(filepath.Join(dir, "source.txt"))

	if err := src.Write([]byte("hello")); err != nil {
		t.Fatalf("Write() unexpected error: %v", err)
	}

	// Destination lives inside a directory that was never created.
	dst := file.New(filepath.Join(dir, "nonexistent-subdir", "dest.txt"))

	if err := src.Copy(dst); err == nil {
		t.Error("Copy() to a non-existent parent directory returned nil, want error")
	}
}

func TestFileHashEmpty(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	f := file.New(filepath.Join(dir, "empty.txt"))

	// Create the file with no content.
	if err := f.Write([]byte{}); err != nil {
		t.Fatalf("Write() unexpected error: %v", err)
	}

	got, err := f.Hash()
	if err != nil {
		t.Fatalf("Hash() on empty file unexpected error: %v", err)
	}

	// sha256 of zero bytes is well-defined; compute it inline.
	sum := sha256.Sum256([]byte{})
	want := hex.EncodeToString(sum[:])

	if got != want {
		t.Errorf("Hash() on empty file = %q, want %q", got, want)
	}
}

func TestFileWhich(t *testing.T) {
	t.Parallel()

	t.Run("finds binary in PATH", func(t *testing.T) {
		t.Parallel()

		// "sh" should be universally available on Linux.
		f := file.New("sh")

		got, err := f.Which()
		if err != nil {
			t.Fatalf("Which() unexpected error: %v", err)
		}

		if !got.IsAbs() {
			t.Errorf("Which() returned non-absolute path %q", got.Path())
		}

		if !got.Exists() {
			t.Errorf("Which() returned path %q that does not exist", got.Path())
		}
	})

	t.Run("not found returns error", func(t *testing.T) {
		t.Parallel()

		f := file.New("definitely-not-a-binary-xyz-123")

		_, err := f.Which()
		if err == nil {
			t.Error("Which() on non-existent binary returned nil, want error")
		}
	})

	t.Run("InPath delegates to Which", func(t *testing.T) {
		t.Parallel()

		if !file.New("sh").InPath() {
			t.Error("InPath() = false for 'sh', want true")
		}

		if file.New("definitely-not-a-binary-xyz-123").InPath() {
			t.Error("InPath() = true for non-existent binary, want false")
		}
	})
}

func TestFileWriteWithPerm(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()

	t.Run("default permission is 0600", func(t *testing.T) {
		t.Parallel()

		f := file.New(filepath.Join(dir, "default.txt"))

		if err := f.Write([]byte("hello")); err != nil {
			t.Fatalf("Write() unexpected error: %v", err)
		}

		info, err := os.Stat(f.Path())
		if err != nil {
			t.Fatalf("Stat() unexpected error: %v", err)
		}

		if got := info.Mode().Perm(); got != 0o600 {
			t.Errorf("default permissions = %o, want %o", got, 0o600)
		}
	})

	t.Run("explicit permission", func(t *testing.T) {
		t.Parallel()

		f := file.New(filepath.Join(dir, "explicit.txt"))
		content := []byte("hello perm")

		if err := f.Write(content, 0o644); err != nil {
			t.Fatalf("Write(data, 0o644) unexpected error: %v", err)
		}

		got, err := f.Read()
		if err != nil {
			t.Fatalf("Read() unexpected error: %v", err)
		}

		if string(got) != string(content) {
			t.Errorf("Read() = %q, want %q", got, content)
		}

		info, err := os.Stat(f.Path())
		if err != nil {
			t.Fatalf("Stat() unexpected error: %v", err)
		}

		if got := info.Mode().Perm(); got != 0o644 {
			t.Errorf("file permissions = %o, want %o", got, 0o644)
		}
	})

	t.Run("multiple perms uses first", func(t *testing.T) {
		t.Parallel()

		f := file.New(filepath.Join(dir, "multi.txt"))

		if err := f.Write([]byte("x"), 0o644, 0o777); err != nil {
			t.Fatalf("Write() unexpected error: %v", err)
		}

		info, err := os.Stat(f.Path())
		if err != nil {
			t.Fatalf("Stat() unexpected error: %v", err)
		}

		if got := info.Mode().Perm(); got != 0o644 {
			t.Errorf("permissions = %o, want %o (should use first perm)", got, 0o644)
		}
	})
}

func TestFileLinesEmpty(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	f := file.New(filepath.Join(dir, "empty.txt"))

	// Write an empty file (zero bytes).
	if err := f.Write([]byte{}); err != nil {
		t.Fatalf("Write() unexpected error: %v", err)
	}

	got, err := f.Lines()
	if err != nil {
		t.Fatalf("Lines() on empty file unexpected error: %v", err)
	}

	// strings.SplitSeq("", "\n") yields a single empty string, matching
	// the same contract as strings.Split("", "\n") = [""].
	want := []string{""}
	if !slices.Equal(got, want) {
		t.Errorf("Lines() on empty file = %v, want %v", got, want)
	}
}
