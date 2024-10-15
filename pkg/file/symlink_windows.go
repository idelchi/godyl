package file

import (
	"fmt"
)

// Symlink creates symlinks for the executable.
func (f File) Symlink(symlinks ...File) error {
	for _, symlink := range symlinks {
		if symlink.Name() == f.Name() {
			continue
		}

		if err := f.Copy(symlink); err != nil {
			return fmt.Errorf("copying %q to %q: %w", f, symlink, err)
		}
	}

	return nil
}
