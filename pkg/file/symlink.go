//go:build linux || darwin

package file

import (
	"fmt"
	"os"
)

// Symlink creates symlinks for the executable.
func (f File) Symlink(symlinks []string) error {
	for _, symlink := range symlinks {
		if symlink == f.String() {
			continue
		}

		err := os.Symlink(f.String(), symlink)
		if err != nil && !os.IsExist(err) {
			return fmt.Errorf("creating symlink for %q: %w", symlink, err)
		}
	}

	return nil
}
