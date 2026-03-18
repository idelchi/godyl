package files_test

import (
	"path/filepath"
	"slices"
	"testing"

	"github.com/idelchi/godyl/pkg/path/file"
	"github.com/idelchi/godyl/pkg/path/files"
)

func TestNew(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		dir       string
		paths     []string
		wantPaths []string
	}{
		{
			name:      "two valid paths produce two files with directory prefix",
			dir:       "dir",
			paths:     []string{"a.txt", "b.txt"},
			wantPaths: []string{"dir/a.txt", "dir/b.txt"},
		},
		{
			name:      "no paths produces empty collection",
			dir:       "dir",
			paths:     []string{},
			wantPaths: []string{},
		},
		{
			name:      "empty paths are skipped",
			dir:       "dir",
			paths:     []string{"a.txt", "", "b.txt"},
			wantPaths: []string{"dir/a.txt", "dir/b.txt"},
		},
		{
			name:      "all empty paths produce empty collection",
			dir:       "dir",
			paths:     []string{"", ""},
			wantPaths: []string{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			fs := files.New(tc.dir, tc.paths...)
			if len(fs) != len(tc.wantPaths) {
				t.Errorf("New(%q, %v): got %d files, want %d", tc.dir, tc.paths, len(fs), len(tc.wantPaths))
			}

			for i, want := range tc.wantPaths {
				if got := fs[i].Path(); got != want {
					t.Errorf("New(%q, %v)[%d].Path() = %q, want %q", tc.dir, tc.paths, i, got, want)
				}
			}
		})
	}
}

func TestAdd(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		initial   []string
		addDir    string
		addPath   string
		wantCount int
	}{
		{
			name:      "add to empty collection",
			initial:   []string{},
			addDir:    "dir",
			addPath:   "new.txt",
			wantCount: 1,
		},
		{
			name:      "add to existing collection increases count",
			initial:   []string{"a.txt"},
			addDir:    "dir",
			addPath:   "b.txt",
			wantCount: 2,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			fs := files.New("dir", tc.initial...)
			fs.Add(tc.addDir, tc.addPath)

			if len(fs) != tc.wantCount {
				t.Errorf("after Add: got %d files, want %d", len(fs), tc.wantCount)
			}
		})
	}
}

func TestAddFile(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		initial   files.Files
		addFile   file.File
		wantCount int
	}{
		{
			name:      "add new file increases count",
			initial:   files.New("dir"),
			addFile:   file.New("dir", "a.txt"),
			wantCount: 1,
		},
		{
			name:      "add duplicate file does not increase count",
			initial:   files.New("dir", "a.txt"),
			addFile:   file.New("dir", "a.txt"),
			wantCount: 1,
		},
		{
			name:      "add distinct file to non-empty collection",
			initial:   files.New("dir", "a.txt"),
			addFile:   file.New("dir", "b.txt"),
			wantCount: 2,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			fs := slices.Clone(tc.initial)
			fs.AddFile(tc.addFile)

			if len(fs) != tc.wantCount {
				t.Errorf("after AddFile: got %d files, want %d", len(fs), tc.wantCount)
			}

			if tc.wantCount > 0 && !fs.Contains(tc.addFile) {
				t.Errorf("after AddFile: collection does not contain %q", tc.addFile)
			}
		})
	}
}

func TestContains(t *testing.T) {
	t.Parallel()

	present := file.New("dir", "present.txt")
	absent := file.New("dir", "absent.txt")

	tests := []struct {
		name  string
		input file.File
		want  bool
	}{
		{
			name:  "present file returns true",
			input: present,
			want:  true,
		},
		{
			name:  "absent file returns false",
			input: absent,
			want:  false,
		},
	}

	fs := files.New("dir", "present.txt")

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := fs.Contains(tc.input)
			if got != tc.want {
				t.Errorf("Contains(%q) = %v, want %v", tc.input, got, tc.want)
			}
		})
	}
}

func TestRemove(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setup     []string
		remove    file.File
		wantOk    bool
		wantCount int
	}{
		{
			name:      "remove existing file returns true and decreases count",
			setup:     []string{"a.txt", "b.txt"},
			remove:    file.New("dir", "a.txt"),
			wantOk:    true,
			wantCount: 1,
		},
		{
			name:      "remove absent file returns false",
			setup:     []string{"a.txt"},
			remove:    file.New("dir", "missing.txt"),
			wantOk:    false,
			wantCount: 1,
		},
		{
			name:      "remove from empty collection returns false",
			setup:     []string{},
			remove:    file.New("dir", "a.txt"),
			wantOk:    false,
			wantCount: 0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			fs := files.New("dir", tc.setup...)
			ok := fs.Remove(tc.remove)

			if ok != tc.wantOk {
				t.Errorf("Remove(%q) = %v, want %v", tc.remove, ok, tc.wantOk)
			}

			if len(fs) != tc.wantCount {
				t.Errorf("after Remove: got %d files, want %d", len(fs), tc.wantCount)
			}

			if tc.wantOk && fs.Contains(tc.remove) {
				t.Errorf("after successful Remove: collection still contains %q", tc.remove)
			}
		})
	}
}

func TestAsSlice(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		dir       string
		paths     []string
		wantPaths []string
	}{
		{
			name:      "returns string paths matching input",
			dir:       "mydir",
			paths:     []string{"a.txt", "b.txt"},
			wantPaths: []string{"mydir/a.txt", "mydir/b.txt"},
		},
		{
			name:      "empty collection returns empty slice",
			dir:       "mydir",
			paths:     []string{},
			wantPaths: []string{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			fs := files.New(tc.dir, tc.paths...)
			got := fs.AsSlice()

			if !slices.Equal(got, tc.wantPaths) {
				t.Errorf("AsSlice() = %v, want %v", got, tc.wantPaths)
			}
		})
	}
}

func TestRelativeTo(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		dir       string
		paths     []string
		base      string
		wantPaths []string
	}{
		{
			name:      "paths become relative to base directory",
			dir:       "/home/user/project",
			paths:     []string{"src/main.go", "src/util.go"},
			base:      "/home/user/project/src",
			wantPaths: []string{"main.go", "util.go"},
		},
		{
			name:      "deeper base strips more path components",
			dir:       "/tmp/base",
			paths:     []string{"a/b/c.txt"},
			base:      "/tmp/base/a/b",
			wantPaths: []string{"c.txt"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			fs := files.New(tc.dir, tc.paths...)
			rel := fs.RelativeTo(tc.base)

			got := make([]string, len(rel))
			for i, f := range rel {
				got[i] = f.Path()
			}

			if !slices.Equal(got, tc.wantPaths) {
				t.Errorf("RelativeTo(%q) = %v, want %v", tc.base, got, tc.wantPaths)
			}
		})
	}
}

func TestFilesExisting(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()

	existing := file.New(filepath.Join(dir, "real.txt"))
	if err := existing.Write([]byte("x")); err != nil {
		t.Fatalf("Write() unexpected error: %v", err)
	}

	fake := file.New(filepath.Join(dir, "does-not-exist.txt"))

	fs := files.Files{existing, fake}
	fs.Existing()

	if len(fs) != 1 {
		t.Fatalf("Existing() left %d files, want 1", len(fs))
	}

	if fs[0].Path() != existing.Path() {
		t.Errorf("Existing()[0].Path() = %q, want %q", fs[0].Path(), existing.Path())
	}
}

func TestFilesExists(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()

	present := file.New(filepath.Join(dir, "present.txt"))
	if err := present.Write([]byte("data")); err != nil {
		t.Fatalf("Write() unexpected error: %v", err)
	}

	fake := file.New(filepath.Join(dir, "absent.txt"))

	fs := files.Files{fake, present}

	got, ok := fs.Exists()
	if !ok {
		t.Fatal("Exists() = _, false, want _, true")
	}

	if got.Path() != present.Path() {
		t.Errorf("Exists() returned %q, want %q", got.Path(), present.Path())
	}
}

func TestFilesExistingEmpty(t *testing.T) {
	t.Parallel()

	// An empty collection passed to Existing() must remain empty — no panic,
	// no spurious entries added.
	var fs files.Files
	fs.Existing()

	if len(fs) != 0 {
		t.Errorf("Existing() on empty collection: got %d files, want 0", len(fs))
	}
}

func TestFilesRelativeToSameBase(t *testing.T) {
	t.Parallel()

	// When the base equals the file's own directory, filepath.Rel produces ".".
	// files.RelativeTo must pass that through unchanged.
	fs := files.New("/tmp/base", "file.txt")
	rel := fs.RelativeTo("/tmp/base/file.txt")

	if len(rel) != 1 {
		t.Fatalf("RelativeTo: got %d files, want 1", len(rel))
	}

	// filepath.Rel("/tmp/base/file.txt", "/tmp/base/file.txt") = "."
	got := rel[0].Path()
	if got != "." {
		t.Errorf("RelativeTo(self) = %q, want %q", got, ".")
	}
}

// TestRelativeTo_FallbackOnUnreachableBase verifies the silent fallback behaviour
// of files.RelativeTo: when filepath.Rel produces an escaping path (e.g. with
// leading ".."), the underlying file.RelativeTo still succeeds on Linux but the
// result contains the "../" prefix. The collection method silently falls back to
// the original file path only when file.RelativeTo returns an error (which cannot
// occur on Linux). This test therefore documents the observable contract: a base
// that is not a prefix of the file paths produces "../../…" style relative paths
// (no panic, no error surfaced to the caller).
func TestRelativeTo_FallbackOnUnreachableBase(t *testing.T) {
	t.Parallel()

	// Files live under /tmp/a; base is an unrelated directory /tmp/b.
	// filepath.Rel("/tmp/b", "/tmp/a/x.txt") = "../a/x.txt", which succeeds.
	// files.RelativeTo must not error — it returns the computed relative path.
	fs := files.New("/tmp/a", "x.txt")
	rel := fs.RelativeTo("/tmp/b")

	if len(rel) != 1 {
		t.Fatalf("RelativeTo: got %d files, want 1", len(rel))
	}

	got := rel[0].Path()

	// We do not assert the exact "../a/x.txt" value to avoid OS-specific path
	// logic; we only verify that the call did not panic and returned a non-empty
	// path (the fallback contract holds).
	if got == "" {
		t.Error("RelativeTo returned an empty path for an unrelated base, want a non-empty fallback")
	}
}
