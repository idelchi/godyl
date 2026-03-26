package folder

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
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
//
// An optional permission mode can be provided (default 0o755).
// If multiple values are given, only the first is used.
func (f Folder) Create(perm ...fs.FileMode) error {
	const defaultPerm = 0o755

	mode := fs.FileMode(defaultPerm)

	if len(perm) > 0 {
		mode = perm[0]
	}

	if err := os.MkdirAll(f.Path(), mode); err != nil {
		return fmt.Errorf("creating directory %q: %w", f, err)
	}

	return nil
}

// Chmod modifies the directory's permission bits.
func (f Folder) Chmod(mode fs.FileMode) error {
	if err := os.Chmod(f.Path(), mode); err != nil {
		return fmt.Errorf("changing permissions of directory %q: %w", f, err)
	}

	return nil
}

// Remove recursively deletes the directory and its contents.
func (f Folder) Remove() error {
	if err := os.RemoveAll(f.Path()); err != nil {
		return fmt.Errorf("removing directory %q: %w", f, err)
	}

	return nil
}

// Info retrieves the file information for the directory.
func (f Folder) Info() (fs.FileInfo, error) {
	info, err := os.Stat(f.Path())
	if err != nil {
		return nil, fmt.Errorf("getting folder info for %q: %w", f, err)
	}

	return info, nil
}

// Size returns the size of the folder in bytes.
func (f Folder) Size() (int64, error) {
	var size int64

	err := filepath.WalkDir(f.Path(), func(_ string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() {
			info, err := d.Info()
			if err != nil {
				return err
			}

			size += info.Size()
		}

		return nil
	})

	return size, err
}
