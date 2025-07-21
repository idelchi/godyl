package files

import (
	"fmt"

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

// LinksFor creates links for each file in the collection.
// Creates a link at each path in the collection,
// pointing to the specified target file.
// See file.File.Links for more details.
func (es Files) LinksFor(file file.File) error {
	if err := file.Links(es...); err != nil {
		return fmt.Errorf("creating links for %q: %w", file, err)
	}

	return nil
}

// Exists returns the first file in the collection that exists.
func (es Files) Exists() (file.File, bool) {
	for _, f := range es {
		if f.Exists() {
			return f, true
		}
	}

	return file.File(""), false
}

// AsSlice converts the Files collection to a slice of strings.
func (es Files) AsSlice() []string {
	slice := make([]string, len(es))
	for i, f := range es {
		slice[i] = f.String()
	}
	return slice
}
