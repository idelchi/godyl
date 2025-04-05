//go:build !windows

package file

import (
	"fmt"
	"os"
)

// Symlink creates symbolic links for the File to each of the provided symlink Files on Linux or Darwin systems.
// If a symlink already exists, it will skip that symlink and continue without returning an error.
// Returns an error if any symlink creation fails (excluding existing symlinks).
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
