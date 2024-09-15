package folder

import (
	"fmt"
	"io/fs"
	"os"
)

// CreateRandomInDir creates a uniquely named directory.
// Creates a directory with a random name inside the specified directory.
// Use empty string for directory to create in system temp directory.
// Pattern is used as a prefix for the random directory name.
func CreateRandomInDir(dir, pattern string) (Folder, error) {
	// Ensure target directory exists before creating subdir
	if dir != "" {
		if err := New(dir).Create(); err != nil {
			return Folder(""), fmt.Errorf("creating temporary directory in %q: %w", dir, err)
		}
	}

	name, err := os.MkdirTemp(dir, pattern)
	if err != nil {
		return Folder(""), fmt.Errorf("creating temporary directory in %q: %w", dir, err)
	}

	return New(name), nil
}

// Create ensures the directory and its parents exist.
// Creates all necessary directories with 0755 permissions.
func (f Folder) Create() error {
	const perm = 0o755

	if err := os.MkdirAll(f.String(), perm); err != nil {
		return fmt.Errorf("creating directory %q: %w", f, err)
	}

	return nil
}

// Remove recursively deletes the directory and its contents.
func (f Folder) Remove() error {
	if err := os.RemoveAll(f.String()); err != nil {
		return fmt.Errorf("removing directory %q: %w", f, err)
	}

	return nil
}

// Info retrieves the file information for the directory.
func (f Folder) Info() (fs.FileInfo, error) {
	info, err := os.Stat(f.String())
	if err != nil {
		return nil, fmt.Errorf("getting folder info for %q: %w", f, err)
	}

	return info, nil
}
