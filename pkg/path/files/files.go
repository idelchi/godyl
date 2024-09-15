package files

import (
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

		fs = append(fs, file.New(dir, path))
	}

	return fs
}

// SymlinksFor creates symlinks from each file to a target.
// Creates a symbolic link at each path in the collection,
// pointing to the specified target file. Returns an error
// if any symlink creation fails.
func (es Files) SymlinksFor(file file.File) error {
	return file.Symlink(es...)
}
