package folder_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/idelchi/godyl/pkg/path/file"
	"github.com/idelchi/godyl/pkg/path/folder"
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
			paths: []string{"a"},
			want:  "a",
		},
		{
			name:  "two components joined",
			paths: []string{"a", "b"},
			want:  "a/b",
		},
		{
			name:  "three components joined",
			paths: []string{"a", "b", "c"},
			want:  "a/b/c",
		},
		{
			name:  "absolute path",
			paths: []string{"/usr", "local"},
			want:  "/usr/local",
		},
		{
			name:  "empty string",
			paths: []string{""},
			want:  "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := folder.New(tc.paths...).Path()
			if got != tc.want {
				t.Errorf("New(%v).Path() = %q, want %q", tc.paths, got, tc.want)
			}
		})
	}
}

func TestJoin(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		base  string
		joins []string
		want  string
	}{
		{
			name:  "join two segments",
			base:  "a",
			joins: []string{"b", "c"},
			want:  "a/b/c",
		},
		{
			name:  "join single segment",
			base:  "x/y",
			joins: []string{"z"},
			want:  "x/y/z",
		},
		{
			name:  "join nothing",
			base:  "a/b",
			joins: []string{},
			want:  "a/b",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := folder.New(tc.base).Join(tc.joins...).Path()
			if got != tc.want {
				t.Errorf("New(%q).Join(%v).Path() = %q, want %q", tc.base, tc.joins, got, tc.want)
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
			name:  "last segment of nested path",
			input: "a/b/c",
			want:  "c",
		},
		{
			name:  "single segment",
			input: "mydir",
			want:  "mydir",
		},
		{
			name:  "two segments",
			input: "parent/child",
			want:  "child",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := folder.New(tc.input).Base()
			if got != tc.want {
				t.Errorf("New(%q).Base() = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}

func TestDir(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "three-segment path returns two-segment parent",
			input: "a/b/c",
			want:  "a/b",
		},
		{
			name:  "two-segment path returns single-segment parent",
			input: "parent/child",
			want:  "parent",
		},
		{
			name:  "single segment returns dot",
			input: "only",
			want:  ".",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := folder.New(tc.input).Dir().Path()
			if got != tc.want {
				t.Errorf("New(%q).Dir().Path() = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}

func TestWithFile(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		dir      string
		filename string
		want     string
	}{
		{
			name:     "simple dir and file",
			dir:      "dir",
			filename: "f.txt",
			want:     "dir/f.txt",
		},
		{
			name:     "nested dir and file",
			dir:      "a/b",
			filename: "c.go",
			want:     "a/b/c.go",
		},
		{
			name:     "file with subdirectory",
			dir:      "root",
			filename: "sub/file.bin",
			want:     "root/sub/file.bin",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := folder.New(tc.dir).WithFile(tc.filename).Path()
			if got != tc.want {
				t.Errorf("New(%q).WithFile(%q).Path() = %q, want %q", tc.dir, tc.filename, got, tc.want)
			}
		})
	}
}

func TestIsSet(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{
			name:  "non-empty path is set",
			input: "a",
			want:  true,
		},
		{
			name:  "empty path is not set",
			input: "",
			want:  false,
		},
		{
			name:  "deeply nested path is set",
			input: "a/b/c",
			want:  true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := folder.New(tc.input).IsSet()
			if got != tc.want {
				t.Errorf("New(%q).IsSet() = %v, want %v", tc.input, got, tc.want)
			}
		})
	}
}

func TestAsFile(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "converts folder path to file path",
			input: "a/b/c",
			want:  "a/b/c",
		},
		{
			name:  "single segment",
			input: "mydir",
			want:  "mydir",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := folder.New(tc.input).AsFile().Path()
			if got != tc.want {
				t.Errorf("New(%q).AsFile().Path() = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}

func TestFolderIsAbs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{
			name:  "absolute path",
			input: "/absolute/dir",
			want:  true,
		},
		{
			name:  "relative path",
			input: "relative/dir",
			want:  false,
		},
		{
			name:  "root",
			input: "/",
			want:  true,
		},
		{
			name:  "dot",
			input: ".",
			want:  false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := folder.New(tc.input).IsAbs()
			if got != tc.want {
				t.Errorf("New(%q).IsAbs() = %v, want %v", tc.input, got, tc.want)
			}
		})
	}
}

func TestFolderAbsolute(t *testing.T) {
	t.Parallel()

	t.Run("relative becomes absolute", func(t *testing.T) {
		t.Parallel()

		f := folder.New("relative/dir")
		got := f.Absolute()

		if !got.IsAbs() {
			t.Errorf("Absolute() returned non-absolute path %q", got.Path())
		}
	})

	t.Run("absolute stays absolute", func(t *testing.T) {
		t.Parallel()

		f := folder.New("/already/absolute")
		got := f.Absolute()

		if got.Path() != "/already/absolute" {
			t.Errorf("Absolute() = %q, want %q", got.Path(), "/already/absolute")
		}
	})
}

func TestCwd(t *testing.T) {
	t.Parallel()

	got, err := folder.Cwd()
	if err != nil {
		t.Fatalf("Cwd() unexpected error: %v", err)
	}

	want, err := os.Getwd()
	if err != nil {
		t.Fatalf("os.Getwd() unexpected error: %v", err)
	}

	if got.Path() != filepath.ToSlash(want) {
		t.Errorf("Cwd().Path() = %q, want %q", got.Path(), filepath.ToSlash(want))
	}

	if !got.IsAbs() {
		t.Errorf("Cwd() returned non-absolute path %q", got.Path())
	}
}

func TestHome(t *testing.T) {
	t.Parallel()

	got, err := folder.Home()
	if err != nil {
		t.Fatalf("Home() unexpected error: %v", err)
	}

	want, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("os.UserHomeDir() unexpected error: %v", err)
	}

	if got.Path() != filepath.ToSlash(want) {
		t.Errorf("Home().Path() = %q, want %q", got.Path(), filepath.ToSlash(want))
	}

	if !got.Exists() {
		t.Errorf("Home() returned path %q that does not exist", got.Path())
	}
}

// ---------------------------------------------------------------------------
// Section 2: Filesystem tests (use t.TempDir)
// ---------------------------------------------------------------------------

func TestFolderCreate(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	f := folder.New(dir, "newdir")

	if f.Exists() {
		t.Fatal("Exists() = true before Create(), want false")
	}

	if err := f.Create(); err != nil {
		t.Fatalf("Create() unexpected error: %v", err)
	}

	if !f.Exists() {
		t.Error("Exists() = false after Create(), want true")
	}
}

func TestFolderCreate_Nested(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	// Create() must create all intermediate directories (MkdirAll behaviour).
	f := folder.New(dir, "a", "b", "c")

	if err := f.Create(); err != nil {
		t.Fatalf("Create() unexpected error for nested path: %v", err)
	}

	if !f.Exists() {
		t.Error("Exists() = false after Create() for nested path, want true")
	}
}

func TestFolderRemove(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	f := folder.New(dir, "toremove")

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

func TestFolderListFiles(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	f := folder.New(dir)

	// Create two files inside the directory.
	names := []string{"alpha.txt", "beta.txt"}
	for _, name := range names {
		if err := file.New(dir, name).Write([]byte("data")); err != nil {
			t.Fatalf("Write(%q) unexpected error: %v", name, err)
		}
	}

	got, err := f.ListFiles()
	if err != nil {
		t.Fatalf("ListFiles() unexpected error: %v", err)
	}

	if len(got) != len(names) {
		t.Fatalf("ListFiles() returned %d files, want %d", len(got), len(names))
	}

	// Verify the actual base names of returned entries match what was created.
	// ListFiles sorts its results, and names is already sorted.
	for i, name := range names {
		if got[i].Base() != name {
			t.Errorf("ListFiles()[%d].Base() = %q, want %q", i, got[i].Base(), name)
		}
	}
}

func TestFolderListFiles_ExcludesSubdirs(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	f := folder.New(dir)

	// One file and one subdirectory — only the file should appear.
	if err := file.New(dir, "file.txt").Write([]byte("x")); err != nil {
		t.Fatalf("Write() unexpected error: %v", err)
	}

	if err := folder.New(dir, "subdir").Create(); err != nil {
		t.Fatalf("Create() unexpected error: %v", err)
	}

	got, err := f.ListFiles()
	if err != nil {
		t.Fatalf("ListFiles() unexpected error: %v", err)
	}

	if len(got) != 1 {
		t.Errorf("ListFiles() returned %d entries, want 1 (subdirs must be excluded)", len(got))
	}
}

func TestFolderListFolders(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	f := folder.New(dir)

	subdirs := []string{"sub1", "sub2"}
	for _, name := range subdirs {
		if err := folder.New(dir, name).Create(); err != nil {
			t.Fatalf("Create(%q) unexpected error: %v", name, err)
		}
	}

	got, err := f.ListFolders()
	if err != nil {
		t.Fatalf("ListFolders() unexpected error: %v", err)
	}

	if len(got) != len(subdirs) {
		t.Fatalf("ListFolders() returned %d folders, want %d", len(got), len(subdirs))
	}

	// Verify the actual base names of returned entries match what was created.
	// ListFolders sorts its results, and subdirs is already sorted.
	for i, name := range subdirs {
		if got[i].Base() != name {
			t.Errorf("ListFolders()[%d].Base() = %q, want %q", i, got[i].Base(), name)
		}
	}
}

func TestFolderListFolders_ExcludesFiles(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	f := folder.New(dir)

	// One subdirectory and one file — only the directory should appear.
	if err := folder.New(dir, "onlysub").Create(); err != nil {
		t.Fatalf("Create() unexpected error: %v", err)
	}

	if err := file.New(dir, "file.txt").Write([]byte("y")); err != nil {
		t.Fatalf("Write() unexpected error: %v", err)
	}

	got, err := f.ListFolders()
	if err != nil {
		t.Fatalf("ListFolders() unexpected error: %v", err)
	}

	if len(got) != 1 {
		t.Errorf("ListFolders() returned %d entries, want 1 (files must be excluded)", len(got))
	}
}

func TestFolderFindFile(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	root := folder.New(dir)

	// Build: dir/sub/target.txt and dir/sub/other.bin
	sub := folder.New(dir, "sub")
	if err := sub.Create(); err != nil {
		t.Fatalf("Create() unexpected error: %v", err)
	}

	if err := file.New(dir, "sub", "target.txt").Write([]byte("hello")); err != nil {
		t.Fatalf("Write(target.txt) unexpected error: %v", err)
	}

	if err := file.New(dir, "sub", "other.bin").Write([]byte("bin")); err != nil {
		t.Fatalf("Write(other.bin) unexpected error: %v", err)
	}

	// The criterion receives the relative path from the root (e.g. "sub/target.txt").
	// Use a doublestar pattern so it matches regardless of nesting depth.
	txtCriteria := func(f file.File) (bool, error) {
		return f.Matches("**/*.txt")
	}

	got, err := root.FindFile(txtCriteria)
	if err != nil {
		t.Fatalf("FindFile() unexpected error: %v", err)
	}

	wantBase := "target.txt"
	if got.Base() != wantBase {
		t.Errorf("FindFile() returned file with Base() = %q, want %q", got.Base(), wantBase)
	}
}

func TestFolderFindFile_NotFound(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	root := folder.New(dir)

	// Place only a .bin file; searching for .go should fail.
	if err := file.New(dir, "data.bin").Write([]byte("x")); err != nil {
		t.Fatalf("Write() unexpected error: %v", err)
	}

	goCriteria := func(f file.File) (bool, error) {
		return f.Matches("*.go")
	}

	_, err := root.FindFile(goCriteria)
	if err == nil {
		t.Error("FindFile() expected error for no match, got nil")
	}
}

// TestFolderFindFile_SingleStar verifies that a single "*" pattern does not
// match across directory separators. A file nested one level deep (sub/f.txt)
// must NOT be matched by "*.txt" when FindFile evaluates the relative path
// "sub/f.txt" against the pattern — doublestar.Match with a bare "*" treats "/"
// as a separator boundary and therefore will not match across it.
func TestFolderFindFile_SingleStar(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	root := folder.New(dir)

	// Place a .txt file one directory level deeper.
	sub := folder.New(dir, "nested")
	if err := sub.Create(); err != nil {
		t.Fatalf("Create() unexpected error: %v", err)
	}

	if err := file.New(dir, "nested", "deep.txt").Write([]byte("x")); err != nil {
		t.Fatalf("Write(deep.txt) unexpected error: %v", err)
	}

	// "*.txt" must not cross the "/" boundary, so "nested/deep.txt" should not match.
	singleStarCriteria := func(f file.File) (bool, error) {
		return f.Matches("*.txt")
	}

	_, err := root.FindFile(singleStarCriteria)
	if err == nil {
		t.Error(
			"FindFile(\"*.txt\") matched a file nested in a subdirectory; single \"*\" must not cross directory separators",
		)
	}
}

func TestFolderSize(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	f := folder.New(dir)

	// Write files with known sizes: 5 bytes + 3 bytes = 8 bytes total.
	files := []struct {
		name    string
		content []byte
	}{
		{"a.txt", []byte("hello")}, // 5 bytes
		{"b.txt", []byte("bye")},   // 3 bytes
	}

	var wantSize int64

	for _, item := range files {
		if err := file.New(filepath.Join(dir, item.name)).Write(item.content); err != nil {
			t.Fatalf("Write(%q) unexpected error: %v", item.name, err)
		}

		wantSize += int64(len(item.content))
	}

	got, err := f.Size()
	if err != nil {
		t.Fatalf("Size() unexpected error: %v", err)
	}

	if got != wantSize {
		t.Errorf("Size() = %d, want %d", got, wantSize)
	}
}

func TestFolderSize_EmptyDir(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	f := folder.New(dir)

	got, err := f.Size()
	if err != nil {
		t.Fatalf("Size() unexpected error on empty dir: %v", err)
	}

	if got != 0 {
		t.Errorf("Size() = %d on empty dir, want 0", got)
	}
}

func TestFolderFromFile(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	f := file.New(filepath.Join(dir, "somefile.txt"))

	if err := f.Write([]byte("x")); err != nil {
		t.Fatalf("Write() unexpected error: %v", err)
	}

	got := folder.FromFile(f)

	// The folder must equal the temp directory (the file's parent).
	if got.Path() != filepath.ToSlash(dir) {
		t.Errorf("FromFile(%q).Path() = %q, want %q", f.Path(), got.Path(), filepath.ToSlash(dir))
	}
}

func TestFolderFindFiles(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	root := folder.New(dir)

	names := []string{"alpha.txt", "beta.txt"}
	for _, name := range names {
		if err := file.New(filepath.Join(dir, name)).Write([]byte("data")); err != nil {
			t.Fatalf("Write(%q) unexpected error: %v", name, err)
		}
	}

	// Also create a non-.txt file to confirm filtering.
	if err := file.New(filepath.Join(dir, "other.bin")).Write([]byte("bin")); err != nil {
		t.Fatalf("Write(other.bin) unexpected error: %v", err)
	}

	txtCriteria := func(f file.File) (bool, error) {
		return f.Matches("*.txt")
	}

	got, err := root.FindFiles(txtCriteria)
	if err != nil {
		t.Fatalf("FindFiles() unexpected error: %v", err)
	}

	if len(got) != len(names) {
		t.Fatalf("FindFiles() returned %d files, want %d", len(got), len(names))
	}

	// FindFiles returns absolute paths sorted; verify both base names are present.
	for i, name := range names {
		if got[i].Base() != name {
			t.Errorf("FindFiles()[%d].Base() = %q, want %q", i, got[i].Base(), name)
		}
	}
}

func TestFolderCreateIdempotent(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	f := folder.New(dir, "idempotent")

	if err := f.Create(); err != nil {
		t.Fatalf("Create() first call unexpected error: %v", err)
	}

	// A second call to Create() on an already-existing directory must not error.
	if err := f.Create(); err != nil {
		t.Errorf("Create() second call unexpected error: %v", err)
	}

	if !f.Exists() {
		t.Error("Exists() = false after double Create(), want true")
	}
}

func TestFolderListFilesNonExistent(t *testing.T) {
	t.Parallel()

	// Point to a path that does not exist on disk.
	f := folder.New("/nonexistent-path-that-should-never-exist")

	_, err := f.ListFiles()
	if err == nil {
		t.Error("ListFiles() on non-existent directory returned nil, want error")
	}
}

func TestFolderCreateWithPerm(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()

	t.Run("default permission is 0755", func(t *testing.T) {
		t.Parallel()

		f := folder.New(dir, "default-perm")

		if err := f.Create(); err != nil {
			t.Fatalf("Create() unexpected error: %v", err)
		}

		info, err := os.Stat(f.Path())
		if err != nil {
			t.Fatalf("Stat() unexpected error: %v", err)
		}

		if got := info.Mode().Perm(); got != 0o755 {
			t.Errorf("default permissions = %o, want %o", got, 0o755)
		}
	})

	t.Run("explicit permission", func(t *testing.T) {
		t.Parallel()

		f := folder.New(dir, "custom-perm")

		if err := f.Create(0o700); err != nil {
			t.Fatalf("Create(0o700) unexpected error: %v", err)
		}

		if !f.Exists() {
			t.Fatal("Exists() = false after Create(0o700), want true")
		}

		info, err := os.Stat(f.Path())
		if err != nil {
			t.Fatalf("Stat() unexpected error: %v", err)
		}

		if got := info.Mode().Perm(); got != 0o700 {
			t.Errorf("directory permissions = %o, want %o", got, 0o700)
		}
	})
}

func TestFolderChmod(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	f := folder.New(dir, "chmod-test")

	if err := f.Create(); err != nil {
		t.Fatalf("Create() unexpected error: %v", err)
	}

	if err := f.Chmod(0o700); err != nil {
		t.Fatalf("Chmod() unexpected error: %v", err)
	}

	info, err := os.Stat(f.Path())
	if err != nil {
		t.Fatalf("Stat() unexpected error: %v", err)
	}

	if got := info.Mode().Perm(); got != 0o700 {
		t.Errorf("permissions after Chmod = %o, want %o", got, 0o700)
	}
}

func TestFolderList(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	f := folder.New(dir)

	// Create a file and a subdirectory.
	if err := file.New(dir, "afile.txt").Write([]byte("x")); err != nil {
		t.Fatalf("Write() unexpected error: %v", err)
	}

	if err := folder.New(dir, "asubdir").Create(); err != nil {
		t.Fatalf("Create() unexpected error: %v", err)
	}

	entries, err := f.List()
	if err != nil {
		t.Fatalf("List() unexpected error: %v", err)
	}

	if len(entries) != 2 {
		t.Fatalf("List() returned %d entries, want 2", len(entries))
	}

	// Verify we got both a file and a directory.
	var hasFile, hasDir bool

	for _, e := range entries {
		if e.IsDir() {
			hasDir = true
		} else {
			hasFile = true
		}
	}

	if !hasFile {
		t.Error("List() missing file entry")
	}

	if !hasDir {
		t.Error("List() missing directory entry")
	}
}

func TestFolderGlob(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	f := folder.New(dir)

	// Create test files.
	for _, name := range []string{"a.go", "b.go", "c.txt"} {
		if err := file.New(dir, name).Write([]byte("x")); err != nil {
			t.Fatalf("Write(%q) unexpected error: %v", name, err)
		}
	}

	// Also create a subdirectory to confirm it doesn't interfere.
	if err := folder.New(dir, "sub").Create(); err != nil {
		t.Fatalf("Create() unexpected error: %v", err)
	}

	t.Run("matches go files", func(t *testing.T) {
		t.Parallel()

		got, err := f.Glob("*.go")
		if err != nil {
			t.Fatalf("Glob(*.go) unexpected error: %v", err)
		}

		if len(got) != 2 {
			t.Fatalf("Glob(*.go) returned %d files, want 2", len(got))
		}
	})

	t.Run("no matches returns empty", func(t *testing.T) {
		t.Parallel()

		got, err := f.Glob("*.rs")
		if err != nil {
			t.Fatalf("Glob(*.rs) unexpected error: %v", err)
		}

		if len(got) != 0 {
			t.Errorf("Glob(*.rs) returned %d files, want 0", len(got))
		}
	})

	t.Run("bad pattern returns error", func(t *testing.T) {
		t.Parallel()

		_, err := f.Glob("[invalid")
		if err == nil {
			t.Error("Glob([invalid) returned nil, want error")
		}
	})
}

func TestFolderWalk(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	root := folder.New(dir)

	// Build: dir/a.txt, dir/sub/b.txt
	if err := file.New(dir, "a.txt").Write([]byte("a")); err != nil {
		t.Fatalf("Write(a.txt) unexpected error: %v", err)
	}

	sub := folder.New(dir, "sub")
	if err := sub.Create(); err != nil {
		t.Fatalf("Create(sub) unexpected error: %v", err)
	}

	if err := file.New(dir, "sub", "b.txt").Write([]byte("b")); err != nil {
		t.Fatalf("Write(b.txt) unexpected error: %v", err)
	}

	t.Run("collects all entries", func(t *testing.T) {
		t.Parallel()

		var paths []string

		err := root.Walk(func(path file.File, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}

			paths = append(paths, path.Base())

			return nil
		})
		if err != nil {
			t.Fatalf("Walk() unexpected error: %v", err)
		}

		// Should have: root dir, a.txt, sub dir, b.txt = 4 entries.
		if len(paths) != 4 {
			t.Errorf("Walk() visited %d entries %v, want 4", len(paths), paths)
		}
	})

	t.Run("SkipDir skips subtree", func(t *testing.T) {
		t.Parallel()

		var fileNames []string

		err := root.Walk(func(path file.File, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if d.IsDir() && path.Base() == "sub" {
				return filepath.SkipDir
			}

			if !d.IsDir() {
				fileNames = append(fileNames, path.Base())
			}

			return nil
		})
		if err != nil {
			t.Fatalf("Walk() unexpected error: %v", err)
		}

		// Only a.txt should be collected; b.txt is inside skipped "sub".
		if len(fileNames) != 1 || fileNames[0] != "a.txt" {
			t.Errorf("Walk() with SkipDir collected %v, want [a.txt]", fileNames)
		}
	})
}

func TestFolderRelativeTo(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()

	// Build dir/a/b — the nested folder whose relative path we will compute.
	child := folder.New(filepath.Join(dir, "a", "b"))
	if err := child.Create(); err != nil {
		t.Fatalf("Create() unexpected error: %v", err)
	}

	parent := folder.New(dir)

	got, err := child.RelativeTo(parent)
	if err != nil {
		t.Fatalf("RelativeTo() unexpected error: %v", err)
	}

	const want = "a/b"
	if got.Path() != want {
		t.Errorf("RelativeTo(%q).Path() = %q, want %q", parent.Path(), got.Path(), want)
	}
}
