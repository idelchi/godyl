//go:build !windows

package file

import (
	"fmt"
	"os"
)

// Symlink creates symbolic links to this file on Unix-like systems.
// Takes multiple target paths and creates a symlink at each location.
// Skips existing symlinks without error, but returns an error if
// symlink creation fails for any other reason. Not available on Windows.
func (f File) Symlink(symlinks ...File) error {
	for _, symlink := range symlinks {
		if symlink.Path() == f.Path() {
			continue
		}

		err := os.Symlink(f.Path(), symlink.Path())
		if err != nil && !os.IsExist(err) {
			return fmt.Errorf("creating symlink for %q: %w", symlink, err)
		}
	}

	return nil
}
