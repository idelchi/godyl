package file

import (
	"fmt"
	"path/filepath"
)

// Files represents a collection of File objects.
type Files []File

// FromStrings creates a Files collection from a list of string paths, relative to the provided directory.
func (Files) FromStrings(dir string, files ...string) Files {
	f := Files{}

	for _, file := range files {
		// Ignore empty strings
		if file == "" {
			continue
		}

		f = append(f, File(filepath.Join(dir, file)))
	}

	return f
}

// Find searches for any of the Files in the given directory.
// It returns the first File found or an error if none are found.
func (es Files) Find(dir string) (File, error) {
	for _, e := range es {
		file, err := e.Find(dir)
		if err == nil {
			return file, nil
		}
	}

	return "", fmt.Errorf("Files %v not found in %q", es.Paths(), dir)
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
func (es Files) SymlinksFor(file File) error {
	return file.Symlink(es...)
}
