package file

import (
	"fmt"
	"path/filepath"
)

type Files []File

// FromStrings creates a list of File from a list of strings.
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
// The first File found is returned.
func (es Files) Find(dir string) (File, error) {
	for _, e := range es {
		file, err := e.Find(dir)
		if err == nil {
			return file, nil
		}
	}
	return "", fmt.Errorf("Files %v not found in %q", es.Paths(), dir)
}

// Paths returns a list of paths for all Files.
func (es Files) Paths() (paths []string) {
	for _, e := range es {
		paths = append(paths, e.String())
	}
	return paths
}

// SymlinksFor creates symlinks for all Files.
func (es Files) SymlinksFor(exe File) error {
	return exe.Symlink(es.Paths())
}
