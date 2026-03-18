package folder_test

import (
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
