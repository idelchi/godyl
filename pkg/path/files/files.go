package files

import (
	"path/filepath"

	"github.com/idelchi/godyl/pkg/path/file"
)

// Files represents a collection of File objects.
type Files []file.File

// New creates a new Files collection from the provided list of paths.
// The paths are joined with the provided directory to create the full file paths.
// Pass `dir` as an empty string to use the paths as-is.
func New(dir string, paths ...string) (fs Files) {
	for _, path := range paths {
		if path == "" {
			continue
		}

		fs = append(fs, file.File(filepath.Join(dir, path)))
	}

	return
}

// NewFromFiles creates a new Files collection from the provided list of File objects.
func NewFromFiles(files ...file.File) (fs Files) {
	for _, file := range files {
		fs = append(fs, file)
	}

	return
}

// Paths returns a list of string paths representing all Files in the collection.
func (es Files) Paths() (paths []string) {
	for _, e := range es {
		paths = append(paths, e.String())
	}

	return paths
}

// SymlinksFor creates symbolic links for all Files in the collection, linking them to the specified target File.
// It returns an error if the operation fails.
func (es Files) SymlinksFor(file file.File) error {
	return file.Symlink(es...)
}
