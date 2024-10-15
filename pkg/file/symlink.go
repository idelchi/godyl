//go:build linux || darwin

package file

import (
	"fmt"
	"os"
)

// Symlink creates symlinks for the executable.
func (f File) Symlink(symlinks ...File) error {
	for _, symlink := range symlinks {
		if symlink.Name() == f.Name() {
			continue
		}

		err := os.Symlink(f.Name(), symlink.Name())
		if err != nil && !os.IsExist(err) {
			return fmt.Errorf("creating symlink for %q: %w", symlink, err)
		}
	}

	return nil
}
