package files

import (
	"path/filepath"

	"github.com/idelchi/godyl/pkg/path/file"
)

// Files represents a collection of filesystem paths.
// Provides batch operations for working with multiple files.
type Files []file.File

// New creates a Files collection from path strings.
// Joins each path with the provided directory to create full paths.
// Empty paths are skipped, and an empty directory uses paths as-is.
func New(dir string, paths ...string) (fs Files) {
	for _, path := range paths {
		if path == "" {
			continue
		}

		fs = append(fs, file.File(filepath.Join(dir, path)))
	}

	return
}

// NewFromFiles creates a Files collection from File objects.
// Takes a list of existing File objects and combines them into
// a single Files collection.
func NewFromFiles(files ...file.File) (fs Files) {
	for _, file := range files {
		fs = append(fs, file)
	}

	return
}

// Paths converts the Files collection to string paths.
// Returns a slice containing the string representation of each file.
func (es Files) Paths() (paths []string) {
	for _, e := range es {
		paths = append(paths, e.String())
	}

	return paths
}

// SymlinksFor creates symlinks from each file to a target.
// Creates a symbolic link at each path in the collection,
// pointing to the specified target file. Returns an error
// if any symlink creation fails.
func (es Files) SymlinksFor(file file.File) error {
	return file.Symlink(es...)
}
