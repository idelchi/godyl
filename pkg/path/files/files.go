package files

import (
	"fmt"
	"slices"

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

// Add adds a new file to the collection.
// Joins the provided path with the directory to create a full path.
func (es *Files) Add(dir, path string) {
	*es = append(*es, file.New(dir, path))
}

// AddFile adds a file to the collection.
// If the file already exists in the collection, it is not added again.
func (es *Files) AddFile(f file.File) {
	if !slices.Contains(*es, f) {
		*es = append(*es, f)
	}
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
		slice[i] = f.Path()
	}

	return slice
}

// Contains checks if the collection contains a specified file.
func (es Files) Contains(file file.File) bool {
	return slices.Contains(es, file)
}

// Remove removes a file from the collection.
// Returns true if the file was found and removed, false otherwise.
func (es *Files) Remove(file file.File) bool {
	index := slices.Index(*es, file)
	if index == -1 {
		return false
	}

	*es = append((*es)[:index], (*es)[index+1:]...)

	return true
}

// RelativeTo makes all files in the collection relative to the specified base directory.
func (es *Files) RelativeTo(base string) Files {
	// preallocate
	relFiles := make(Files, 0, len(*es))

	for _, f := range *es {
		rel, err := f.RelativeTo(base)
		if err != nil {
			rel = f
		}

		relFiles = append(relFiles, rel)
	}

	return relFiles
}
